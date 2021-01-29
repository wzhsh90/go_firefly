package middleware

import (
	"firefly/config"
	"firefly/consts"
	models "firefly/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Session 初始化session
func Session() gin.HandlerFunc {
	bootConfig := config.GetAppConfig()
	store := cookie.NewStore([]byte(bootConfig.Server.SessionKey))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"})
	return sessions.Sessions(bootConfig.Server.SessionName, store)
}

func SaveCurrentUser(ctx *gin.Context, user *models.LoginUser) {
	session := sessions.Default(ctx)
	session.Set(consts.SEESION_USER_NAME, user)
	err := session.Save()
	if err != nil {
		config.ZapAppLogger.Sugar().Infow("SaveCurrentUser",
			"xx", err.Error(),
		)
	}
}
func RemoveUserSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete(consts.SEESION_USER_NAME)
	err := session.Save()
	if err != nil {
		config.ZapAppLogger.Sugar().Infow("RemoveUserSession",
			"xx", err.Error(),
		)
	}
}
func GetCurrentUser(ctx *gin.Context) *models.LoginUser {
	session := sessions.Default(ctx)

	org := session.Get(consts.SEESION_USER_NAME)
	if org != nil {
		return org.(*models.LoginUser)
	} else {
		return nil
	}
}
