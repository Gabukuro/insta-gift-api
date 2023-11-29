package prediction

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Gabukuro/insta-gift-api/internal/pkg/uuid"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/rs/zerolog"
)

type (
	Service struct {
		repository              Repository
		logger                  *zerolog.Logger
		SNSClient               snsiface.SNSAPI
		predictionEventTopicArn string
	}

	ServiceParams struct {
		Ctx                     context.Context
		Repository              Repository
		Logger                  *zerolog.Logger
		SNSClient               snsiface.SNSAPI
		PredictionEventTopicArn string
	}

	Message struct {
		ID         uuid.UUID       `json:"id,omitempty"`
		Entity     string          `json:"entity,omitempty"`
		ProducedAt time.Time       `json:"produced_at,omitempty"`
		Payload    MessagePayload  `json:"payload,omitempty"`
		Metadata   MessageMetadata `json:"meta,omitempty"`
	}

	MessagePayload any

	MessageMetadata struct {
		PredictionId string `json:"prediction_id,omitempty"`
		Username     string `json:"username,omitempty"`
	}
)

func NewService(opt ServiceParams) *Service {
	return &Service{
		repository:              opt.Repository,
		logger:                  opt.Logger,
		SNSClient:               opt.SNSClient,
		predictionEventTopicArn: opt.PredictionEventTopicArn,
	}
}

func (s *Service) GetPredictionByUsername(ctx context.Context, username string, prediction *Prediction) error {
	return s.repository.GetPredictionByUsername(ctx, username, prediction)
}

func (s *Service) CreatePrediction(ctx context.Context, username string) (*Prediction, error) {
	prediction := &Prediction{
		Username: username,
		Status:   PredictionStatusPending,
	}
	prediction.NewUUID()

	err := s.repository.CreatePrediction(ctx, prediction)
	if err != nil {
		return nil, err
	}

	_, err = s.SendPredictionMessage(ctx, prediction)
	if err != nil {
		return nil, err
	}

	return prediction, nil
}

func (s *Service) CheckIfPredictionExistsAndReturnItsStatus(ctx context.Context, username string) (*PredictionStatus, bool) {
	prediction := &Prediction{
		Username: username,
	}

	err := s.repository.GetPredictionByUsername(ctx, username, prediction)
	if err != nil {
		return nil, false
	}

	return &prediction.Status, true
}

func (s *Service) SendPredictionMessage(ctx context.Context, prediction *Prediction) (*sns.PublishOutput, error) {
	message := &Message{
		ID:         uuid.New(),
		Entity:     "prediction_created",
		ProducedAt: time.Now().UTC(),
		Payload:    prediction,
		Metadata: MessageMetadata{
			PredictionId: prediction.ID.String(),
			Username:     prediction.Username,
		},
	}

	return s.Publish(ctx, message)
}

func (s *Service) Publish(ctx context.Context, message *Message) (*sns.PublishOutput, error) {
	bytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	publishInput := &sns.PublishInput{
		TopicArn: &s.predictionEventTopicArn,
		Message:  aws.String(string(bytes)),
	}

	return s.SNSClient.Publish(publishInput)
}
