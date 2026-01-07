module github.com/pingplex/pingplex

go 1.25.0

require (
	github.com/go-core-fx/config v0.1.0
	github.com/go-core-fx/fiberfx v0.1.0
	github.com/go-core-fx/logger v0.0.1
	github.com/go-core-fx/redisfx v0.0.0-20251029094515-c9e3d82dfaa2
	github.com/gocql/gocql v1.7.0
	github.com/prometheus/client_golang v1.23.2
	github.com/scylladb/gocqlx/v3 v3.0.4
	go.uber.org/fx v1.24.0
	go.uber.org/zap v1.27.1
)

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/ansrivas/fiberprometheus/v2 v2.14.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-core-fx/fxutil v0.0.0-20251027105421-acea37162eb9 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gofiber/contrib/fiberzap/v2 v2.1.6 // indirect
	github.com/gofiber/fiber/v2 v2.52.10 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/parsers/dotenv v1.1.0 // indirect
	github.com/knadh/koanf/parsers/yaml v1.1.0 // indirect
	github.com/knadh/koanf/providers/env/v2 v2.0.0 // indirect
	github.com/knadh/koanf/providers/file v1.2.0 // indirect
	github.com/knadh/koanf/v2 v2.3.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.4 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/redis/go-redis/v9 v9.17.2 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.68.0 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.39.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

// Use the latest version of scylladb/gocql; check for updates at https://github.com/scylladb/gocql/releases
replace github.com/gocql/gocql => github.com/scylladb/gocql v1.17.0
