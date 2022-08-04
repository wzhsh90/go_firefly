package models

type BootConfig struct {
	Server Server
	Db     DataSource
}
type Server struct {
	Mode        string
	Port        string
	Profile     string
	SessionName string
	SessionKey  string
	TplPath     string
}

func (c *Server) Prod() bool {
	return c.Profile == "prod"
}

type DataSource struct {
	Url       string
	MaxLife   int
	MaxIdle   int
	MaxOpen   int
	LogEnable bool
}
