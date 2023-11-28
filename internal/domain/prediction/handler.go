package prediction

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type (
	httpHandler struct {
		service *Service
		logger  *zerolog.Logger
	}
)

func NewHTTPHandler(app *fiber.App, service *Service, logger *zerolog.Logger) {
	handler := &httpHandler{
		service: service,
		logger:  logger,
	}

	predictionGroup := app.Group("/prediction/:username")

	predictionGroup.Get("/", handler.GetPredictionByUsername)
	predictionGroup.Post("/", handler.EnqueuePrediction)
}

func (h *httpHandler) GetPredictionByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	prediction := &Prediction{}
	if err := h.service.GetPredictionByUsername(ctx.Context(), username, prediction); err != nil {
		return err
	}

	return ctx.JSON(prediction)
}

func (h *httpHandler) EnqueuePrediction(ctx *fiber.Ctx) error {
	return nil
}
