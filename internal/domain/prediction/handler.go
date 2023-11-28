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

	Response struct {
		Prediction *Prediction `json:"prediction,omitempty"`
		Message    string      `json:"message,omitempty"`
	}
)

func NewHTTPHandler(app *fiber.App, service *Service, logger *zerolog.Logger) {
	handler := &httpHandler{
		service: service,
		logger:  logger,
	}

	predictionGroup := app.Group("/prediction/:username")

	predictionGroup.Get("/", handler.GetPredictionByUsername)
	predictionGroup.Post("/", handler.CreatePrediction)
}

func (h *httpHandler) GetPredictionByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	prediction := &Prediction{}
	if err := h.service.GetPredictionByUsername(ctx.Context(), username, prediction); err != nil {
		return err
	}

	return ctx.JSON(&Response{Prediction: prediction})
}

func (h *httpHandler) CreatePrediction(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	if status, exists := h.service.CheckIfPredictionExistsAndReturnItsStatus(ctx.Context(), username); exists {
		switch *status {
		case PredictionStatusPending:
			return ctx.JSON(&Response{Message: "Prediction is pending"})
		case PredictionStatusCompleted:
			return h.GetPredictionByUsername(ctx)
		}
	}

	err := h.service.CreatePrediction(ctx.Context(), username)
	if err != nil {
		return err
	}

	return ctx.JSON(&Response{Message: "Prediction created"})
}
