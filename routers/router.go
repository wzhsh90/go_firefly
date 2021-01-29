package routers

import (
	"encoding/gob"
	"firefly/config"
	"firefly/consts"
	"firefly/event"
	"firefly/middleware"
	models "firefly/model"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
	"strings"
	"time"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.String()
		loginedFlag := false
		excludeArr := []string{"login", "static", "file", "doc"}
		for _, v := range excludeArr {
			if strings.Contains(url, v) {
				loginedFlag = true
			}
		}
		if !loginedFlag {
			v := middleware.GetCurrentUser(c)
			if v != nil && v.Id != "" {
				loginedFlag = true
			}
		}
		if loginedFlag || url == "/" || url == "" {
			c.Next()
			return
		} else {
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" && !strings.Contains(url, "login") {
				rest := models.RestResult{}
				rest.Code = 100
				c.AbortWithStatusJSON(200, rest)
			} else {
				c.Redirect(302, "/login")
				return
			}
		}

	}
}
func InitRouter(r *gin.Engine, server config.Server) {
	gob.Register(&models.LoginUser{})
	r.MaxMultipartMemory = 100 << 20 //100M
	r.Use(middleware.Session())
	//r.Use(middleware.Cors())
	//r.Use(middleware.ZapHttpLogger())
	r.Use(AuthMiddleWare())
	r.Static("/static", "./static")
	//r.HTMLRender = ginview.Default()
	ctime := strconv.FormatInt(time.Now().Unix(), 10)
	r.HTMLRender = ginview.New(goview.Config{
		Root:      "views",
		Extension: ".html",
		//Master:    "layouts/master",
		Partials: []string{},
		Funcs: template.FuncMap{
			"htmlBr": func(a string) template.HTML {
				return template.HTML(strings.ReplaceAll(a, "\n", "<br/>"))
			},
			"htmlOrg": func(a string) template.HTML {
				return template.HTML(a)
			},
			"urlFn": func(s string) template.URL {
				return template.URL(s)
			},
			"addFn": func(a, b int) int {
				return a + b
			},
			"neFn": func(a, b int64) bool {
				return a != b
			},
			"rankFn": func(idx int, page, pageSize int64) int64 {
				return int64(idx) + (page-1)*pageSize + 1
			},
			"add64Fn": func(a, b int64) int64 {
				return a + b
			},
			"subFn": func(a, b int) int {
				return a - b
			},
			"sub64Fn": func(a, b int64) int64 {
				return a - b
			},
			"modFn": func(a, b int) bool {
				return a%b == 0
			},
			"mod64Fn": func(a, b int64) bool {
				return a%b == 0
			},
			"LookupPath": func(a string) string {
				return a + "?v=" + ctime
			},
		},
		DisableCache: false,
		Delims:       goview.Delims{Left: "${", Right: "}"},
	})
	event.Publish(consts.ROUTER_INIT_EVENT, r)
	_ = r.Run(server.Port) // listen and serve on 0.0.0.0:8080
}
