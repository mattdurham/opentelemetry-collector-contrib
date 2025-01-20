// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheusremotewriteexporter

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/prometheus/prompb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func doNothingExportSink(_ context.Context, reqL []*prompb.WriteRequest) error {
	_ = reqL
	return nil
}

func TestWALCreation_nilConfig(t *testing.T) {
	config := (*WALConfig)(nil)
	pwal := newWAL(config, doNothingExportSink)
	require.Nil(t, pwal)
}

func TestWALCreation_nonNilConfig(t *testing.T) {
	config := &WALConfig{Directory: t.TempDir()}
	pwal := newWAL(config, doNothingExportSink)
	require.NotNil(t, pwal)
	assert.NoError(t, pwal.stop())
}

func orderByLabelValueForEach(reqL []*prompb.WriteRequest) {
	for _, req := range reqL {
		orderByLabelValue(req)
	}
}

func orderByLabelValue(wreq *prompb.WriteRequest) {
	// Sort the timeSeries by their labels.
	type byLabelMessage struct {
		label  *prompb.Label
		sample *prompb.Sample
	}

	for _, timeSeries := range wreq.Timeseries {
		bMsgs := make([]*byLabelMessage, 0, len(wreq.Timeseries)*10)
		for i := range timeSeries.Labels {
			bMsgs = append(bMsgs, &byLabelMessage{
				label:  &timeSeries.Labels[i],
				sample: &timeSeries.Samples[i],
			})
		}
		sort.Slice(bMsgs, func(i, j int) bool {
			return bMsgs[i].label.Value < bMsgs[j].label.Value
		})

		for i := range bMsgs {
			timeSeries.Labels[i] = *bMsgs[i].label
			timeSeries.Samples[i] = *bMsgs[i].sample
		}
	}

	// Now finally sort stably by timeseries value for
	// which just .String() is good enough for comparison.
	sort.Slice(wreq.Timeseries, func(i, j int) bool {
		ti, tj := wreq.Timeseries[i], wreq.Timeseries[j]
		return ti.String() < tj.String()
	})
}

func TestWALStopManyTimes(t *testing.T) {
	tempDir := t.TempDir()
	config := &WALConfig{
		Directory:         tempDir,
		TruncateFrequency: 60 * time.Microsecond,
		BufferSize:        1,
	}
	pwal := newWAL(config, doNothingExportSink)
	require.NotNil(t, pwal)

	// Ensure that invoking .stop() multiple times doesn't cause a panic, but actually
	// First close should NOT return an error.
	require.NoError(t, pwal.stop())
	for i := 0; i < 4; i++ {
		// Every invocation to .stop() should return an errAlreadyClosed.
		require.ErrorIs(t, pwal.stop(), errAlreadyClosed)
	}
}

func TestWAL_persist(t *testing.T) {
	// Unit tests that requests written to the WAL persist.
	config := &WALConfig{Directory: t.TempDir()}

	pwal := newWAL(config, doNothingExportSink)
	require.NotNil(t, pwal)

	// 1. Write out all the entries.
	reqL := []*prompb.WriteRequest{
		{
			Timeseries: []prompb.TimeSeries{
				{
					Labels:  []prompb.Label{{Name: "ts1l1", Value: "ts1k1"}},
					Samples: []prompb.Sample{{Value: 1, Timestamp: 100}},
				},
			},
		},
		{
			Timeseries: []prompb.TimeSeries{
				{
					Labels:  []prompb.Label{{Name: "ts2l1", Value: "ts2k1"}},
					Samples: []prompb.Sample{{Value: 2, Timestamp: 200}},
				},
				{
					Labels:  []prompb.Label{{Name: "ts1l1", Value: "ts1k1"}},
					Samples: []prompb.Sample{{Value: 1, Timestamp: 100}},
				},
			},
		},
	}

	ctx := context.Background()
	require.NoError(t, pwal.retrieveWALIndices())
	t.Cleanup(func() {
		assert.NoError(t, pwal.stop())
	})

	require.NoError(t, pwal.persistToWAL(reqL))

	// 2. Read all the entries from the WAL itself, guided by the indices available,
	// and ensure that they are exactly in order as we'd expect them.
	wal := pwal.wal
	start, err := wal.FirstIndex()
	require.NoError(t, err)
	end, err := wal.LastIndex()
	require.NoError(t, err)

	var reqLFromWAL []*prompb.WriteRequest
	for i := start; i <= end; i++ {
		req, err := pwal.readPrompbFromWAL(ctx, i)
		require.NoError(t, err)
		reqLFromWAL = append(reqLFromWAL, req)
	}

	orderByLabelValueForEach(reqL)
	orderByLabelValueForEach(reqLFromWAL)
	require.Equal(t, reqLFromWAL[0], reqL[0])
	require.Equal(t, reqLFromWAL[1], reqL[1])
}

func TestWal(t *testing.T) {

}

func TestWALDuplicateDataPrevention(t *testing.T) {
	// Create a temporary directory for WAL
	config := &WALConfig{
		Directory:         t.TempDir(),
		BufferSize:        1,
		TruncateFrequency: 1 * time.Second,
	}

	// Track exported requests to detect duplicates
	ids := make(map[string]prompb.TimeSeries)
	exportSink := func(_ context.Context, reqL []*prompb.WriteRequest) error {
		for _, req := range reqL {
			for _, ts := range req.Timeseries {
				if _, found := ids[ts.Labels[0].Name]; found {
					t.Errorf("duplicate data prevents duplicate label %q", ts.Labels[0].Name)
				}
				ids[ts.Labels[0].Name] = ts
			}
		}
		return nil
	}

	// Create WAL with custom export sink
	pwal := newWAL(config, exportSink)
	require.NotNil(t, pwal)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	nop := zap.NewNop()
	ctx = contextWithLogger(ctx, nop)
	err := pwal.run(ctx)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		// Persist duplicate requests to WAL
		require.NoError(t, pwal.persistToWAL(makeReq(i)))

	}

	require.Eventuallyf(t, func() bool {
		return len(ids) == 10*100 // i * j
	}, 5*time.Second, 100*time.Millisecond, "exported count expected 1_000, received %d", len(ids))

	// The context cancel has to be called so the mutex isnt continually locked on stop.
	cancel()
	// Stop WAL
	require.NoError(t, pwal.stop())
}

func makeReq(i int) []*prompb.WriteRequest {
	wr := make([]*prompb.WriteRequest, 0)
	for j := 0; j < 1; j++ {
		ts := &prompb.TimeSeries{
			Labels:  []prompb.Label{{Name: fmt.Sprintf("test_metric_%d_%d", j, i), Value: strconv.Itoa(i)}},
			Samples: []prompb.Sample{{Value: 42, Timestamp: time.Now().UnixNano()}},
		}
		wr = append(wr, &prompb.WriteRequest{Timeseries: []prompb.TimeSeries{*ts}})
	}
	return wr
}
