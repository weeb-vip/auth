package notification

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/weeb-vip/auth/config"
)

type snsService struct {
	svc *sns.SNS
	cfg config.SNSConfig
}

func NewSNSService(cfg config.SNSConfig) NotificationService {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sns.New(sess)

	return &snsService{
		svc,
		cfg,
	}
}

func (service snsService) Publish(topic string, message any) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	output, err := service.svc.Publish(&sns.PublishInput{
		Message:  aws.String(string(messageJSON)),
		TopicArn: aws.String(topic),
	})

	log.Println(output)

	if err != nil {
		return err
	}

	return nil
}
