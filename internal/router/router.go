package router

import (
	"baas-auth/internal/configs"
	"baas-auth/internal/controllers"
	"baas-auth/internal/services"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func InitAPI(appConfig *configs.Config, service *services.Service) humacli.CLI {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		humaConfig := huma.DefaultConfig("Auth API", "0.2.0")
		humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"bearer": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		}
		humaConfig.Servers = []*huma.Server{
			{URL: fmt.Sprintf("http://localhost:%d", options.Port)},
		}

		app := fiber.New()
		app.Use(logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		}))

		api := humafiber.New(app, humaConfig)
		v1Api := huma.NewGroup(api, "/api/v1")

		userController := controllers.InitUserController(appConfig, service)

		huma.AutoRegister(v1Api, userController)

		hooks.OnStart(func() {
			app.Listen(fmt.Sprintf(":%d", options.Port))
		})
		hooks.OnStop(func() {
			service.DB.Close()
		})
	})

	return cli
}
