package connection

import (
	"fmt"
	"sanjay/config"

	"github.com/twilio/twilio-go"
)

func Velidation(r *config.Twillow) bool {
	if r.Sid == "" || r.WatappNum == "" || r.Token == "" {
		return false
	}
	return true
}
func NewTwillowConn() (*twilio.RestClient, error) {
	r := config.LoadEnv().GetTillowInfo()
	if exits := Velidation(r); !exits {
		return nil, fmt.Errorf("twillow env missing")
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: r.Sid,
		Password: r.Token,
	})
	return client, nil
}
