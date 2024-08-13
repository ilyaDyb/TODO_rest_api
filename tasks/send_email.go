package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/hibiken/asynq"
	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const TypeEmailDelivery = "email:deliver"

type EmailDeliveryPayload struct {
	RecieverEmail string
	Message []byte
}

func NewEmailDeliveryTask(reciever string, msg []byte) (*asynq.Task, error) {
    payload, err := json.Marshal(EmailDeliveryPayload{RecieverEmail: reciever, Message: msg})
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
			"service": "asynq",
		}).Errorf("Failed to start EmailDeliveryTask with error: %v", err)
        return nil, err
    }
    return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
    var p EmailDeliveryPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        logger.Log.WithFields(logrus.Fields{
			"service": "asynq",
		}).Errorf("Failed to unmarchal data with error: %v", err)
        return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
    }

    password := os.Getenv("SMTP_PASSWORD")
    sender := os.Getenv("SMTP_EMAIL")
    receiver := []string{
        p.RecieverEmail,
    }

    smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    message := p.Message
    auth := smtp.PlainAuth("", sender, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, receiver, message)
    if err != nil {
        log.Println(err)
        return fmt.Errorf("error: %s", err)
    }
    logger.Log.WithFields(logrus.Fields{
        "service": "asynq",
        "email": receiver,
    }).Info("Email was sent successfully")
    return nil
}