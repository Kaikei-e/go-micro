package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

type Main struct {
	router *gin.Engine
}

func (m *Main) initServer() error {
	var err error

	err = common.LoadConfig()
	if err != nil {
		return err
	}

	err = databases.Database.Init()
	if err != nil {
		return err
	}

	if common.Config.EnableGinFieldLog {
		f, _ := os.Create("logs/gin.log")
		if common.Config.EnableGinConsoleLog {
			gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
		} else {
			gin.DefaultWriter = io.MultiWriter(f)
		}

	} else {
		if common.Config.EnableGinConsoleLog {
			gin.DefaultWriter = io.MultiWriter()
		}
	}

	m.router = gin.Default()

	return nil

}

// @title MovieManagement Service API Document
// @version 1.0
// @description List APIs of MovieManagement Service
// @termsOfService http://swagger.io/terms/

// @host 107.113.53.47:8809
// @BasePath /api/v1

func main() {
	m := Main{}

	if m.initServer() != nil {
		return
	}

	defer databases.Database.Close()

	c := controllers.Movie{}

	v1 := m.router.Group("/api/v1")
	{
		v1.POST("/login", c.Login)
		v1.Get("/movies/list", c.ListMovies)

		// APIs need to use token string
		v1.Use(jwt.Auth(common.Config.JwtSecretPassword))
		v1.POST("/movies", c.AddMovie)
	}

	m.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	m.router.Run(common.Config.Port)
}
