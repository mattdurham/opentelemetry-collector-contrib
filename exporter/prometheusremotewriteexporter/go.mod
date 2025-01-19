module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter

go 1.23

toolchain go1.23.4

require (
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/fsnotify/fsnotify v1.8.0
	github.com/go-kit/log v0.2.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang/snappy v0.0.4
	github.com/grafana/walqueue v0.0.0-20250113171943-e5fe545d1408
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.117.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.117.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.117.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.117.0
	github.com/prometheus/client_golang v1.20.5
	github.com/prometheus/prometheus v0.55.1
	github.com/stretchr/testify v1.10.0
	github.com/tidwall/wal v1.1.8
	go.opentelemetry.io/collector/component v0.117.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/component/componenttest v0.117.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/config/confighttp v0.117.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/config/configopaque v1.23.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/config/configretry v1.23.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/config/configtelemetry v0.117.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/config/configtls v1.23.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/confmap v1.23.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/consumer/consumererror v0.117.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/exporter v0.117.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/exporter/exportertest v0.117.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/featuregate v1.23.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/collector/pdata v1.23.1-0.20250117002813-e970f8bb1258
	go.opentelemetry.io/otel v1.32.0
	go.opentelemetry.io/otel/metric v1.32.0
	go.opentelemetry.io/otel/sdk/metric v1.32.0
	go.opentelemetry.io/otel/trace v1.32.0
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deneonet/benc v1.1.2 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gammazero/deque v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grafana/regexp v0.0.0-20240518133315-a468a5bfb3bc // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.61.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/tidwall/gjson v1.10.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tinylib/msgp v1.2.4 // indirect
	github.com/vladopajic/go-actor v0.9.1-0.20241115212052-39d92aec6093 // indirect
	go.opentelemetry.io/collector/client v1.23.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/config/configauth v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.23.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/consumer v1.23.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/consumer/consumertest v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/consumer/xconsumer v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/exporter/xexporter v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/extension v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/extension/auth v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/extension/xextension v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/pipeline v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/receiver v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/receiver/receivertest v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/receiver/xreceiver v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/collector/semconv v0.117.1-0.20250117002813-e970f8bb1258 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.56.0 // indirect
	go.opentelemetry.io/otel/sdk v1.32.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241104194629-dd2ea8efbc28 // indirect
	google.golang.org/grpc v1.69.4 // indirect
	google.golang.org/protobuf v1.36.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../internal/coreinternal

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry => ../../pkg/resourcetotelemetry

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus => ../../pkg/translator/prometheus

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite => ../../pkg/translator/prometheusremotewrite

retract (
	v0.76.2
	v0.76.1
	v0.65.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil => ../../pkg/pdatautil

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest => ../../pkg/pdatatest

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden => ../../pkg/golden
