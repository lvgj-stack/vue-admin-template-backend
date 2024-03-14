package gobase

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/Mr-LvGJ/jota/access_log"

	"github.com/Mr-LvGJ/gobase/pkg/gobase/controller/v1/table"

	"github.com/Mr-LvGJ/gobase/pkg/common/constant"
	"github.com/Mr-LvGJ/gobase/pkg/common/core"
	"github.com/Mr-LvGJ/gobase/pkg/common/errno"
	"github.com/Mr-LvGJ/gobase/pkg/common/middleware"
	"github.com/Mr-LvGJ/gobase/pkg/common/setting"
	"github.com/Mr-LvGJ/gobase/pkg/common/token"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/controller/v1/user"
	"github.com/Mr-LvGJ/gobase/pkg/gobase/store"
)

func LoadRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine, mw ...gin.HandlerFunc) {
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.RequestID())
	g.Use(access_log.GinAccessLogInterceptor(access_log.WithLogConfig(setting.C().AccessLog)))
}

func installController(g *gin.Engine) {
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.InternalServerError, nil)
	})
	g.GET("/healthz", func(c *gin.Context) {
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})
	dbClient := store.Client()
	userController := user.NewUserController(dbClient)
	tableController := table.NewTableController(dbClient)
	pprof.Register(g)
	v1 := g.Group("/v1")
	{
		userV1 := v1.Group("/users")
		{
			userV1.POST("/login", userController.Login)
			userV1.POST("", userController.Create)
			userV1.Use(authMiddleware())
			userV1.POST("/logout", userController.Logout)
			userV1.GET("/info", userController.Info)
			userV1.DELETE(":name", userController.Delete)
			userV1.PUT(":name", userController.Update)
			userV1.GET("", userController.List)
			userV1.GET(":name", userController.Get)
		}
		tableV1 := v1.Group("/tables")
		{
			tableV1.Use(authMiddleware())
			tableV1.GET("/list", tableController.List)
			tableV1.POST("", tableController.Create)
		}
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, errno.ErrToken, nil)
			c.Abort()
			return
		}
		c.Set(constant.XUsernameKey, username)
		c.Next()
	}
}
