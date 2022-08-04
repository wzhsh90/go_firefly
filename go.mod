module firefly

go 1.15

replace github.com/firefly => ./

require (
	github.com/antonmedv/expr v1.9.0
	github.com/asaskevich/EventBus v0.0.0-20180315140547-d46933a94f05
	github.com/fatih/color v1.13.0
	github.com/foolin/goview v0.3.0
	github.com/fsnotify/fsnotify v1.5.4
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.8.1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-playground/validator/v10 v10.10.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/huandu/go-sqlbuilder v1.15.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20220307211146-efcb8507fb70 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.0.0-20200103221440-774c71fcf114 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.3.5
	gorm.io/gorm v1.23.8

)
