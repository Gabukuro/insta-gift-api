package prediction

import (
	"github.com/Gabukuro/insta-gift-api/internal/domain/product"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type (
	httpHandler struct {
		service        *Service
		productService *product.Service
		logger         *zerolog.Logger
	}

	Response struct {
		Prediction *Prediction `json:"prediction,omitempty"`
		Message    string      `json:"message,omitempty"`
	}
)

func NewHTTPHandler(app *fiber.App, service *Service, productService *product.Service, logger *zerolog.Logger) {
	handler := &httpHandler{
		service:        service,
		productService: productService,
		logger:         logger,
	}

	predictionGroup := app.Group("/prediction/:username")

	predictionGroup.Get("/", handler.GetPredictionByUsername)
	predictionGroup.Post("/", handler.CreatePrediction)
}

func (h *httpHandler) GetPredictionByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	prediction := &Prediction{}
	err := h.service.GetPredictionByUsername(ctx.Context(), username, prediction)
	if err != nil {
		return err
	}

	prediction.Products, err = h.productService.GetProductsByPredictionID(ctx.Context(), prediction.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(&Response{Prediction: prediction})
}

func (h *httpHandler) CreatePrediction(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	if _, exists := h.service.CheckIfPredictionExistsAndReturnItsStatus(ctx.Context(), username); exists {
		return h.GetPredictionByUsername(ctx)
	}

	prediction, err := h.service.CreatePrediction(ctx.Context(), username)
	if err != nil {
		return err
	}

	return ctx.JSON(&Response{Prediction: prediction})
}
