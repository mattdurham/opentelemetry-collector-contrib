package prometheusremotewriteexporter

import (
	"context"
	"github.com/go-kit/log"
	"github.com/golang/snappy"
	"github.com/grafana/walqueue/implementations/prometheus"
	"github.com/grafana/walqueue/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter/internal/metadatatest"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkImplementations(b *testing.B) {
	seriesCount := atomic.Uint32{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		data, err = snappy.Decode(nil, data)
		if err != nil {
			panic(err)
		}
		var req prompb.WriteRequest
		err = req.Unmarshal(data)
		if err != nil {
			panic(err)
		}
		require.NotEmpty(b, req.Timeseries)
		for _, ts := range req.Timeseries {
			seriesCount.Add(uint32(len(ts.Samples)))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	intSumBatch := testdata.GenerateMetricsManyMetricsSameResource(1_000)

	type tt struct {
		send func(m pmetric.Metrics)
		name string
	}
	tests := []tt{
		{
			send: buildPR(b, srv),
			name: "existing_prometheusremotewriteexporter",
		},
		{
			send: buildWalQueue(b, srv),
			name: "new_walqueue",
		},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			runs := 0
			for i := 0; i < b.N; i++ {
				test.send(intSumBatch)
				require.Eventually(b, func() bool {
					return seriesCount.Load() == 2002
				}, 5*time.Second, 100*time.Millisecond)
				runs++
				seriesCount.Store(0)
			}
		})
	}

}
func createPR(b *testing.B, srv *httptest.Server) *prwExporter {

	// Adjusted retry settings for faster testing
	retrySettings := configretry.BackOffConfig{
		Enabled:         true,
		InitialInterval: 100 * time.Millisecond, // Shorter initial interval
		MaxInterval:     1 * time.Second,        // Shorter max interval
		MaxElapsedTime:  2 * time.Second,        // Shorter max elapsed time
	}
	clientConfig := confighttp.NewDefaultClientConfig()
	clientConfig.Endpoint = srv.URL
	clientConfig.ReadBufferSize = 512 * 1024
	clientConfig.WriteBufferSize = 512 * 1024
	cfg := &Config{
		Namespace:         "",
		ClientConfig:      clientConfig,
		MaxBatchSizeBytes: 3000000,
		RemoteWriteQueue:  RemoteWriteQueue{NumConsumers: 1},
		TargetInfo: &TargetInfo{
			Enabled: true,
		},
		CreatedMetric: &CreatedMetric{
			Enabled: true,
		},
		BackOffConfig: retrySettings,
	}

	cfg.WAL = &WALConfig{
		Directory: b.TempDir(),
		// This is the writerequests to write, which in our tests actually contains 2002 samples.
		// So always write
		BufferSize: 1,
		// For it to write more often.
		TruncateFrequency: 1 * time.Second,
	}

	buildInfo := component.BuildInfo{
		Description: "OpenTelemetry Collector",
		Version:     "1.0",
	}
	tel := metadatatest.SetupTelemetry()
	set := tel.NewSettings()
	// detailed level enables otelhttp client instrumentation which we dont want to test here
	set.MetricsLevel = configtelemetry.LevelBasic
	set.BuildInfo = buildInfo

	prwe, nErr := newPRWExporter(cfg, set)

	require.NoError(b, nErr)

	return prwe
}
func buildPR(b *testing.B, srv *httptest.Server) func(m pmetric.Metrics) {
	ctx, cancel := context.WithCancel(context.Background())

	prwe := createPR(b, srv)
	b.Cleanup(func() {
		cancel()
		prwe.Shutdown(ctx)
	})
	require.NoError(b, prwe.Start(ctx, componenttest.NewNopHost()))
	return func(m pmetric.Metrics) {
		err := prwe.PushMetrics(ctx, m)
		require.NoError(b, err)
	}
}

func buildWalQueue(b *testing.B, srv *httptest.Server) func(m pmetric.Metrics) {
	q, err := prometheus.NewQueue("test", types.ConnectionConfig{
		URL:           srv.URL,
		Timeout:       1 * time.Second,
		BatchCount:    100,
		FlushInterval: 500 * time.Millisecond,
		Connections:   2,
	}, b.TempDir(), 100, 100*time.Millisecond, 1*time.Hour, prom.NewRegistry(), "test", log.NewNopLogger())
	require.NoError(b, err)
	q.Start()
	b.Cleanup(func() {
		q.Stop()
	})
	// Only using prwe here to ensure we use the same exporter settings
	prwe := createPR(b, srv)
	return func(m pmetric.Metrics) {
		app := q.Appender(context.Background())
		tsMap, err := prometheusremotewrite.FromMetrics(m, prwe.exporterSettings)
		for _, ts := range tsMap {
			lbls := make(labels.Labels, len(ts.Labels))
			for i, lbl := range ts.Labels {
				lbls[i] = labels.Label{
					Name:  lbl.Name,
					Value: lbl.Value,
				}
			}
			for _, sample := range ts.Samples {
				// Timestamp is nano but we want milli so using the below so we dont hit ttl issues.
				_, aErr := app.Append(0, lbls, time.Now().UnixMilli(), sample.Value)
				require.NoError(b, aErr)
			}
		}
		err = app.Commit()
		require.NoError(b, err)
	}

}
