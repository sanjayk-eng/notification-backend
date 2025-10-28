package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const fromNumber = "whatsapp:+14155238886"

type Twillow struct {
	Client *twilio.RestClient
}

func NewTwillowClient(conn *twilio.RestClient) *Twillow {
	return &Twillow{
		Client: conn,
	}
}

func (t *Twillow) SendMessage(fullPhone, otp string) (string, error) {
	message := fmt.Sprintf("Your verification code is %s. It will expire in 1 minute.", otp)

	params := &openapi.CreateMessageParams{}
	params.SetTo("whatsapp:" + fullPhone)
	params.SetFrom(fromNumber)
	params.SetBody(message)

	resp, err := t.Client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Twilio API error: %v", err)
		return "", fmt.Errorf("failed to send OTP: %v", err)
	}

	// --- Check SID ---
	if resp.Sid == nil {
		return "", errors.New("message SID missing in response")
	}
	messageSID := *resp.Sid
	log.Printf("Message SID: %s", messageSID)

	// --- Check Initial Status ---
	initialStatus := ""
	if resp.Status != nil {
		initialStatus = *resp.Status
		log.Printf("Initial Status: %s", initialStatus)
	}

	// --- Wait for Twilio to process (2–3 sec delay) ---
	time.Sleep(3 * time.Second)

	// --- Fetch Latest Message Status ---
	statusResp, err := t.Client.Api.FetchMessage(messageSID, nil)
	if err != nil {
		log.Printf("Failed to fetch message status: %v", err)
		return messageSID, fmt.Errorf("failed to fetch message status: %v", err)
	}

	finalStatus := ""
	if statusResp.Status != nil {
		finalStatus = *statusResp.Status
	} else {
		finalStatus = "unknown"
	}

	log.Printf("Final message status for %s: %s", messageSID, finalStatus)

	// --- Return SID + Final Status Together ---
	return fmt.Sprintf("%s|%s", messageSID, finalStatus), nil
}
