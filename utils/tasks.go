package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"golang.org/x/net/context"
)

const TypeEmailDelivery = "email:deliver"

type EmailDeliveryPayload struct {
	RecieverEmail string
	SenderEmail string
}

func NewEmailDeliveryTask(reciever string, sender string) (*asynq.Task, error) {
    payload, err := json.Marshal(EmailDeliveryPayload{RecieverEmail: reciever, SenderEmail: sender})
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
    log.Printf("Sending Email to User: user_email=%s, sender_email=%s", p.RecieverEmail, p.SenderEmail)
    // Email delivery code ...
    return nil
}