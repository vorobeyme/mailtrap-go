package mailtrap

import (
	"encoding/json"
	"io"
)

// Events is the wrapper around the Webhook event.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/b9cdfe3d25137-receive-events
type Events struct {
	Events []Event `json:"events"`
}

// Event represents an email event that user is subscribed to.
type Event struct {
	Event           string            `json:"event"`
	Email           string            `json:"email"`
	Category        string            `json:"category"`
	MessageID       string            `json:"message_id"`
	CustomVariables map[string]string `json:"custom_variables"`
	EventID         string            `json:"event_id"`
	Timestamp       int               `json:"timestamp"`
	Response        string            `json:"response"`
	ResponseCode    int               `json:"response_code"`
	Reason          string            `json:"reason"`
	IP              string            `json:"ip"`
	UserAgent       string            `json:"user_agent"`
	URL             string            `json:"url"`
}

func DecodeWebhook(r io.Reader) (*Events, error) {
	e := new(Events)
	if err := json.NewDecoder(r).Decode(&e); err != nil {
		return nil, err
	}
	return e, nil
}
