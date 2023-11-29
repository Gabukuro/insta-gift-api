package router

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
)

type (
	Router struct {
		app    *fiber.App
		logger *zerolog.Logger
	}

	Options struct {
		AppName string
		Logger  *zerolog.Logger
	}
)

func NewRouter(options *Options) *Router {
	routerInstance := &Router{
		app: fiber.New(
			fiber.Config{
				AppName:     options.AppName,
				JSONEncoder: json.Marshal,
				JSONDecoder: json.Unmarshal,
			},
		),
		logger: options.Logger,
	}

	routerInstance.app.Use(cors.New())

	routerInstance.app.Get("/status", func(ctx *fiber.Ctx) (err error) {
		_ = ctx.SendString("OK")
		return nil
	})

	return routerInstance
}

func (r *Router) Start(addr string) {
	if err := r.app.Listen(addr); err != nil {
		panic(err)
	}
}

func (r *Router) Stop() {
	_ = r.App().Shutdown()
}

func (r *Router) App() *fiber.App {
	return r.app
}
