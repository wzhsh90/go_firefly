module firefly

go 1.15

replace github.com/firefly => ./

require (
	github.com/asaskevich/EventBus v0.0.0-20180315140547-d46933a94f05
	github.com/foolin/goview v0.3.0
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.8.1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-sql-driver/mysql v1.6.0
	github.com/upper/db/v4 v4.5.4
	go.uber.org/zap v1.13.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0

)
