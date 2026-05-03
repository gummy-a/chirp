package email

import (
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func NewAwsConfig(ctx context.Context) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AUTH_SERVICE_AWS_SES_REGION")),
	)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewEmailSender(cfg aws.Config, ctx context.Context) *EmailSender {
	return &EmailSender{
		cfg: cfg,
		ctx: ctx,
	}
}

type EmailSender struct {
	cfg aws.Config
	ctx context.Context
}

func (e *EmailSender) Send(to value_object.Email, numberCode value_object.NumberCode, tmpAccountID *value_object.TemporaryAccountID) error {
	env := os.Getenv("AUTH_SERVICE_APP_ENV")
	if env == "development" {
		fmt.Printf("email: %s, numbercode: %d\n", to.String(), numberCode)
		return nil
	}

	client := sesv2.NewFromConfig(e.cfg)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(os.Getenv("AUTH_SERVICE_AWS_SES_FROM_ADDRESS")),
		Destination: &types.Destination{
			ToAddresses: []string{
				to.String(),
			},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String("verification code"),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(fmt.Sprintf("visit %s/signup/?token=%s to verify your account.\n Your verification code is %d.\n", os.Getenv("AUTH_SERVICE_FRONTEND_URL"), tmpAccountID.String(), numberCode)),
					},
				},
			},
		},
	}

	_, err := client.SendEmail(e.ctx, input)
	if err != nil {
		return err
	}
	return nil
}
