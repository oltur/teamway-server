package main

import (
	"github.com/oltur/teamway-server/controller"
	"github.com/oltur/teamway-server/docs"
	_ "github.com/oltur/teamway-server/docs"
)

var version string

// @title           Teamway test task: Server
// @description     This is a Teamway test task: Server, based on celler example.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url NA
// http://www.swagger.io/support
// @contact.email  olturua@gmail.com

// @version 0
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	s := version
	if s == "" {
		s = "(DEVELOPMENT BUILD)"
	}

	docs.SwaggerInfo.Version = s
	r, _ := controller.SetupRouter()
	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
}
