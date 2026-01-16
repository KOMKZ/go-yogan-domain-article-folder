module github.com/KOMKZ/go-yogan-domain-article-folder

go 1.24.1

require (
	github.com/KOMKZ/go-yogan-domain-article v0.0.0
	github.com/KOMKZ/go-yogan-domain-folder v0.0.0
)

require (
	github.com/KOMKZ/go-yogan-framework v0.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/text v0.32.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gorm.io/gorm v1.31.1 // indirect
)

replace (
	github.com/KOMKZ/go-yogan-domain-article => ../go-yogan-domain-article
	github.com/KOMKZ/go-yogan-domain-folder => ../go-yogan-domain-folder
	github.com/KOMKZ/go-yogan-framework => ../../go-yogan-framework
)
