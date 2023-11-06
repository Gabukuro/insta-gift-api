package config

import (
	"encoding/json"
	"strings"

	"github.com/Gabukuro/insta-gift-api/internal/pkg/awsprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
)

type (
	Config struct {
		GoEnv                string `env:"GO_ENV" envDefault:"development"`
		AwsRegion            string `env:"AWS_REGION" envDefault:"us-east-1"`
		SecretsManagerRegion string `env:"SECRETS_MANAGER_REGION" envDefault:""`
		AwsProvider          awsprovider.AWSProvider
		awsSecretsManager    secretsmanageriface.SecretsManagerAPI
		logger               *zerolog.Logger

		cachedSecrets map[string]string
	}

	secrets map[string]any
)

func New(logger *zerolog.Logger) (cfg *Config) {
	cfg = new(Config)
	cfg.logger = logger
	cfg.cachedSecrets = make(map[string]string)

	cfg.parseEnv()
	cfg.initAWS()

	return
}

func (cfg *Config) parseEnv() {
	opts := &env.Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			cfg.logger.Debug().Msgf("Set %s to %v (default? %v)\n", tag, value, isDefault)
		},
	}
	if err := env.Parse(cfg, *opts); err != nil {
		cfg.logger.Fatal().Err(err).Msg("failed to load config")
	}
}

func (cfg *Config) initAWS() {
	cfg.AwsProvider = awsprovider.New(cfg.AwsRegion, cfg.GoEnv)

	awsRegion := cfg.AwsRegion
	if cfg.SecretsManagerRegion != "" {
		awsRegion = cfg.SecretsManagerRegion
	}

	cfg.awsSecretsManager = awsprovider.New(awsRegion, cfg.GoEnv).SecretManager()
}

func (cfg *Config) GetSecretValue(secret string, secretData any) {
	if cfg.verifyCachedSecret(secret, secretData) {
		return
	}

	secretId := aws.String(cfg.GoEnv + "." + secret)
	secretValueInput := &secretsmanager.GetSecretValueInput{
		SecretId: secretId,
	}

	receivedSecret, err := cfg.awsSecretsManager.GetSecretValue(secretValueInput)
	if err != nil {
		cfg.logger.Error().Err(err).Msgf("failed to load secret: %s", secret)
		panic(err)
	}

	cfg.cachedSecrets[secret] = *receivedSecret.SecretString

	err = json.Unmarshal([]byte(*receivedSecret.SecretString), secretData)
	if err != nil {
		cfg.logger.Error().Err(err).Msgf("Failed to unmarshal secret: %s", secret)
		panic(err)
	}
}

func (cfg *Config) verifyCachedSecret(secret string, secretData any) bool {
	secretString, ok := cfg.cachedSecrets[secret]
	if ok {
		err := json.Unmarshal([]byte(secretString), secretData)
		if err != nil {
			cfg.logger.Error().Err(err).Msgf("Failed to unmarshal secret: %s", secret)
			panic(err)
		}
	}

	return ok
}

func (cfg *Config) GetSecretString(property string) (secret string) {
	secretMap := make(map[string]any)
	secretOk := false

	for index, property := range strings.Split(property, ".") {
		if index == 0 {
			cfg.GetSecretValue(property, &secretMap)
			continue
		}

		if actualProperty, ok := secretMap[property]; ok {
			if actualPropertyMap, ok := actualProperty.(map[string]any); ok {
				secretMap = actualPropertyMap
				continue
			}

			if actualPropertyString, ok := actualProperty.(string); ok {
				secret = actualPropertyString
				secretOk = true
			}
		}
	}

	if secret == "" && !secretOk {
		panic("secret not found:" + property)
	}

	return
}

func (cfg *Config) GetDataBaseSecret(database string) (databaseURL string) {
	secrets := make(secrets)
	cfg.GetSecretValue("database", &secrets)

	if database, ok := secrets[database]; ok {
		if databaseString, ok := database.(string); ok {
			return databaseString
		}
	}

	return
}
