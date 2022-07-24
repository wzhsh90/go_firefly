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
	Dialect   string
	User      string
	Password  string
	Database  string
	Host      string
	Socket    string
	MaxLife   int
	MaxIdle   int
	MaxOpen   int
	LogEnable bool
}
