package awsprovider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type (
	provider struct {
		Region  string
		Env     string
		session *session.Session
	}

	AWSProvider interface {
		Session(force bool) (*session.Session, error)
		SQS() sqsiface.SQSAPI
		SNS() snsiface.SNSAPI
		SecretManager() secretsmanageriface.SecretsManagerAPI
	}
)

func New(region string, env string) AWSProvider {
	return &provider{Region: region, Env: env}
}

func (awsProvider *provider) Session(force bool) (*session.Session, error) {
	if awsProvider.session != nil && !force {
		return awsProvider.session, nil
	}

	var err error

	config := &aws.Config{
		Region: aws.String(awsProvider.Region),
	}

	// force localstack credentials
	if awsProvider.Env == "" || awsProvider.Env == "development" {
		config.Endpoint = aws.String("http://s3.localhost.localstack.cloud:4566")
	}

	awsProvider.session, err = session.NewSession(config)

	if err != nil {
		return nil, err
	}

	return awsProvider.session, nil
}

func (awsProvider *provider) SQS() sqsiface.SQSAPI {
	ses := session.Must(awsProvider.Session(false))
	return sqs.New(ses)
}

func (awsProvider *provider) SNS() snsiface.SNSAPI {
	ses := session.Must(awsProvider.Session(false))
	return sns.New(ses)
}

func (awsProvider *provider) SecretManager() secretsmanageriface.SecretsManagerAPI {
	ses := session.Must(awsProvider.Session(false))
	return secretsmanager.New(ses)
}
