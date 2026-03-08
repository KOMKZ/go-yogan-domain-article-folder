module github.com/KOMKZ/go-yogan-domain-article-folder

go 1.24.1

require (
	github.com/KOMKZ/go-yogan-domain-article v0.0.0
	github.com/KOMKZ/go-yogan-domain-folder v0.0.0
	github.com/samber/do/v2 v2.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/samber/go-type-to-string v1.8.0 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gorm.io/gorm v1.31.1 // indirect
)

replace github.com/KOMKZ/go-yogan-domain-article => ../go-yogan-domain-article

replace github.com/KOMKZ/go-yogan-domain-folder => ../go-yogan-domain-folder

require github.com/KOMKZ/go-yogan-framework v0.0.0

replace github.com/KOMKZ/go-yogan-framework => ../../go-yogan-framework
