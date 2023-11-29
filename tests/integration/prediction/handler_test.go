package prediction_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Gabukuro/insta-gift-api/internal/domain/prediction"
	"github.com/Gabukuro/insta-gift-api/internal/domain/product"
	awsprovider "github.com/Gabukuro/insta-gift-api/internal/pkg/awsprovider/mock"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/database"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/log"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/router"
	"github.com/Gabukuro/insta-gift-api/internal/pkg/testhelper"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestPredictionHandler(t *testing.T) {
	t.Parallel()

	tests := []testhelper.Test{
		{
			Description: "should get prediction by username",
			Input: &testhelper.Input{
				Method: http.MethodGet,
				Path:   "/prediction/test",
			},
			Expected: &testhelper.Expected{
				Status: http.StatusOK,
				Body:   `{"prediction":{"id":"!uuid","created_at":"!date","updated_at":"!date","username":"test","feedback_rating":5,"status":"completed"}}`,
			},
		},
		{
			Description: "should create prediction",
			Input: &testhelper.Input{
				Method: http.MethodPost,
				Path:   "/prediction/test2",
			},
			Expected: &testhelper.Expected{
				Status: http.StatusOK,
				Body:   `{"prediction":{"id":"!uuid","created_at":"!date","updated_at":"!date","username":"test2","feedback_rating":0,"status":"pending"}}`,
			},
		},
	}

	RunTests(t, tests)
}

func RunTests(t *testing.T, tests []testhelper.Test) {
	ctx := context.Background()
	logger := log.New(zerolog.ErrorLevel)
	testDb := testhelper.NewTestDatabase(logger)
	testDbUrl, err := testDb.Create()
	assert.Equal(t, nil, err)

	defer func() {
		err := testDb.Drop()
		assert.Equal(t, nil, err)
	}()

	db := database.New(testDbUrl, 1, logger).Connect()
	defer func() {
		err := db.Close()
		assert.Equal(t, nil, err)
	}()

	bunDB := bun.NewDB(db.DB, pgdialect.New())

	routerInstance := router.NewRouter(&router.Options{
		AppName: "test",
		Logger:  logger,
	})

	mcokCtrl := gomock.NewController(t)
	snsClient := awsprovider.NewMockSNSAPI(mcokCtrl)
	snsClient.EXPECT().Publish(gomock.Any()).Return(nil, nil)

	productRepo := product.NewRepository(bunDB, logger)
	productService := product.NewService(product.ServiceParams{
		Ctx:        ctx,
		Repository: productRepo,
		Logger:     logger,
	})

	predictionRepo := prediction.NewRepository(bunDB, logger)
	predictionService := prediction.NewService(prediction.ServiceParams{
		Ctx:                     ctx,
		Repository:              predictionRepo,
		Logger:                  logger,
		SNSClient:               snsClient,
		PredictionEventTopicArn: "arn:aws:sns:us-east-2:000000000000:profile-analysis-events",
	})

	prediction.NewHTTPHandler(routerInstance.App(), predictionService, productService, logger)

	testhelper.TestAllRequests(tests, t, routerInstance)
}
