package models

type BootConfig struct {
	Server Server
	Db     DataSource
}
type Server struct {
	Mode        string
	Port        string
	SessionName string
	SessionKey  string
	TplPath     string
}
type DataSource struct {
	Url       string
	MaxLife   int
	MaxIdle   int
	MaxOpen   int
	LogEnable bool
}
