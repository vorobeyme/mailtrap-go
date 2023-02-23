package mailtrap

import (
	"strings"
	"testing"
)

func TestWebhook_DecodeWebhook(t *testing.T) {
	_, _, teardown := setupTestingClient()
	defer teardown()

	const webhookEvents = `{
		"events": [
			{
				"event": "delivery",
				"email": "john@example.com",
				"category": "Password reset",
				"message_id": "12345678-abcd-efgh-yyyy-1111111111",
				"event_id": "98765432-abcd-edfg-xxxx-2222222222",
				"custom_variables": {
					"user_id": "45982",
					"batch_id": "PSJ-12"
				},
				"timestamp": 123456789011
			}
		]
	}`

	jsonData := strings.NewReader(webhookEvents)
	res, err := DecodeWebhook(jsonData)
	if err != nil {
		t.Errorf("DecodeWebhook returned error: %v", err)
	}

	expected := Events{Events: []Event{
		{
			Event:           "delivery",
			Email:           "john@example.com",
			Category:        "Password reset",
			MessageID:       "12345678-abcd-efgh-yyyy-1111111111",
			EventID:         "98765432-abcd-edfg-xxxx-2222222222",
			CustomVariables: map[string]string{"user_id": "45982", "batch_id": "PSJ-12"},
			Timestamp:       123456789011,
		},
	}}

	testJSONMarshal(t, &expected, webhookEvents)

	if cv := res.Events[0].CustomVariables; len(cv) != 2 {
		t.Errorf("DecodeWebhook expected 2 variables, got: %v", len(cv))
	}

	data := strings.NewReader(`{{"bad": "json"}}`)
	_, err = DecodeWebhook(data)
	if err == nil {
		t.Error("DecodeWebhook err = nil, want error")
	}
}
