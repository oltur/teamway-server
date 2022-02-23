package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetupRouter() (*gin.Engine, *Controller) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.MaxMultipartMemory = 20 << 20 // 20 MiB

	// TODO: Change to specific CORS rules?
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-Total-Count", "Authorization"}
	config.ExposeHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-Total-Count", "Authorization"}
	r.Use(cors.New(config))

	c := NewController()

	v1 := r.Group("/api/v1")
	{
		testTaken := v1.Group("/test-taken")
		{
			testTaken.Use(c.Auth())
			testTaken.GET("/next", c.GetNextQuestion)
			testTaken.POST("", c.AnswerQuestion)
			testTaken.GET("", c.GetTestResult)
		}
		user := v1.Group("/user")
		{
			user.POST("", c.AddUser)
			user.POST("/login", c.Login)
			user.POST("/logout", c.Auth(), c.Logout)
			user.GET(":id", c.Auth(), c.ShowUser)
			user.DELETE(":id", c.Auth(), c.DeleteUser)
			user.PATCH(":id", c.Auth(), c.UpdateUser)
		}
		test := v1.Group("/test")
		{
			test.GET(":id", c.ShowTest)
			test.GET("", c.ListTests)
			test.POST("", c.Auth(), c.AddTest)
			test.DELETE(":id", c.Auth(), c.DeleteTest)
		}
		testByTitle := v1.Group("/test-by-title")
		{
			testByTitle.GET(":title", c.ShowTestByTitle)
		}
		tools := v1.Group("/utils")
		{
			tools.GET("/ping", c.Ping)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r, c
}
