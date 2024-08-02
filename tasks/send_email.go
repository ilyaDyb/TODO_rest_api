package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/hibiken/asynq"
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
        return nil, err
    }
    return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
    var p EmailDeliveryPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
    }
    log.Printf("Sending Email to User: user_email=%s, message=%s", p.RecieverEmail, p.Message)
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
    return nil
}