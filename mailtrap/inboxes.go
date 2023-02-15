package mailtrap

import (
	"fmt"
	"net/http"
)

type InboxesServiceContract interface {
	Create(accountID, inboxID int, name string) (*Inbox, *Response, error)
	Update(accountID, inboxID int, updRequest *UpdateInboxRequest) (*Inbox, *Response, error)
	List(accountID int) ([]*Inbox, *Response, error)
	Get(accountID, inboxID int) (*Inbox, *Response, error)
	Delete(accountID, inboxID int) (*Response, error)
	Clean(accountID, inboxID int) (*Inbox, *Response, error)
	MarkAsRead(accountID, inboxID int) (*Inbox, *Response, error)
	ResetCredentials(accountID, inboxID int) (*Inbox, *Response, error)
	EnableEmail(accountID, inboxID int) (*Inbox, *Response, error)
	ResetEmail(accountID, inboxID int) (*Inbox, *Response, error)
}

type InboxesService struct {
	client *Client
}

var _ InboxesServiceContract = &InboxesService{}

// Inbox represents a Mailtrap inbox.
type Inbox struct {
	ID                      int         `json:"id"`
	Name                    string      `json:"name"`
	Username                string      `json:"username"`
	Password                string      `json:"password"`
	MaxSize                 int         `json:"max_size"`
	Status                  string      `json:"status"`
	EmailUsername           string      `json:"email_username"`
	EmailUsernameEnabled    bool        `json:"email_username_enabled"`
	SentMessagesCount       int         `json:"sent_messages_count"`
	ForwardedMessagesCount  int         `json:"forwarded_messages_count"`
	Used                    bool        `json:"used"`
	ForwardFromEmailAddress string      `json:"forward_from_email_address"`
	ProjectID               int         `json:"project_id"`
	Domain                  string      `json:"domain"`
	POP3Domain              string      `json:"pop3_domain"`
	EmailDomain             string      `json:"email_domain"`
	EmailsCount             int         `json:"emails_count"`
	EmailsUnreadCount       int         `json:"emails_unread_count"`
	LastMessageSentAt       string      `json:"last_message_sent_at"`
	SMTPPorts               []int       `json:"smtp_ports"`
	POP3Ports               []int       `json:"pop3_ports"`
	MaxMessageSize          int         `json:"max_message_size"`
	Permissions             Permissions `json:"permissions"`
}

type createInboxRequest struct {
	Inbox struct {
		Name string `json:"name"`
	} `json:"inbox"`
}

func (s *InboxesService) Messages() *MessagesService {
	return s.client.Messages
}

// Create creates an inbox in a project.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/86631e73937e2-create-an-inbox
func (s *InboxesService) Create(accountID, inboxID int, name string) (*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/projects/%d/inboxes", accountID, inboxID)
	payload := &createInboxRequest{
		Inbox: struct {
			Name string `json:"name"`
		}{Name: name},
	}

	return s.makeRequest(u, http.MethodPost, payload)
}

type UpdateInboxRequest struct {
	Name          string `json:"name,omitempty"`
	EmailUsername string `json:"email_username,omitempty"`
}

// Update updates inbox name, inbox email username.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/768067eceee9d-update-an-inbox
func (s *InboxesService) Update(accountID, inboxID int, opts *UpdateInboxRequest) (*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d", accountID, inboxID)
	payload := struct {
		Inbox *UpdateInboxRequest `json:"inbox"`
	}{opts}

	return s.makeRequest(u, http.MethodPatch, payload)
}

// List returns the list of inboxes.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/49dd3b9d6806f-get-a-list-of-inboxes
func (s *InboxesService) List(accountID int) ([]*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes", accountID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var inbox []*Inbox
	res, err := s.client.Do(req, &inbox)
	if err != nil {
		return nil, res, err
	}

	return inbox, res, err
}

// Get returns attributes of the inbox.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/432a39abe34b3-get-inbox-attributes
func (s *InboxesService) Get(accountID, inboxID int) (*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d", accountID, inboxID)
	return s.makeRequest(u, http.MethodGet, nil)
}

// Delete removes an inbox with all its emails.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/e624770632299-delete-project
func (s *InboxesService) Delete(accountID, inboxID int) (*Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d", accountID, inboxID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Clean delete all messages (emails) from inbox.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/8a1e782a64fd0-clean-inbox
func (s *InboxesService) Clean(accountID, inboxID int) (*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/clean", accountID, inboxID)
	return s.makeRequest(u, http.MethodPatch, nil)
}

// MarkAsRead mark all messages in the inbox as read.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/8a38b0494dff1-mark-as-read
func (s *InboxesService) MarkAsRead(accountID, inboxID int) (*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/all_read", accountID, inboxID)
	return s.makeRequest(u, http.MethodPatch, nil)
}

// ResetCredentials resets SMTP credentials of the inbox.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/403fd0f1315e6-reset-credentials
func (s *InboxesService) ResetCredentials(accountID, inboxID int) (*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/reset_credentials", accountID, inboxID)
	return s.makeRequest(u, http.MethodPatch, nil)
}

// EnableEmail enables email address.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/a4b31a4c40ae4-enable-email-address
func (s *InboxesService) EnableEmail(accountID, inboxID int) (*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/toggle_email_username", accountID, inboxID)
	return s.makeRequest(u, http.MethodPatch, nil)
}

// ResetEmail reset email address
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/5ebb1ca46e3d0-reset-email-address
func (s *InboxesService) ResetEmail(accountID, inboxID int) (*Inbox, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/reset_email_username", accountID, inboxID)
	return s.makeRequest(u, http.MethodPatch, nil)
}

func (s *InboxesService) makeRequest(endpoint, httpMethod string, payload interface{}) (*Inbox, *Response, error) {
	req, err := s.client.NewRequest(httpMethod, endpoint, payload)
	if err != nil {
		return nil, nil, err
	}

	var inbox *Inbox
	res, err := s.client.Do(req, &inbox)
	if err != nil {
		return nil, res, err
	}

	return inbox, res, nil
}
