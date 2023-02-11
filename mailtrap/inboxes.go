package mailtrap

import (
	"fmt"
	"net/http"
)

const (
	getInboxesEndpoint       = "/accounts/%d/inboxes"
	getInboxEndpoint         = "/accounts/%d/inboxes/%d"
	deleteInboxEndpoint      = "/accounts/%d/inboxes/%d"
	updateInboxEndpoint      = "/accounts/%d/inboxes/%d"
	cleanInboxEndpoint       = "/accounts/%d/inboxes/%d/clean"
	markAsReadEndpoint       = "/accounts/%d/inboxes/%d/all_read"
	resetCredentialsEndpoint = "/accounts/%d/inboxes/%d/reset_credentials"
	resetEmailEndpoint       = "/accounts/%d/inboxes/%d/reset_email_username"
	enableEmailEndpoint      = "/accounts/%d/inboxes/%d/toggle_email_username"
	createInboxEndpoint      = "/accounts/%d/projects/%d/inboxes"
)

type InboxesServiceContract interface {
	ListInboxes(accountID int) ([]*Inbox, *Response, error)
	GetInbox(accountID, inboxID int) (*Inbox, *Response, error)
	DeleteInbox(accountID, inboxID int) (*Response, error)
	UpdateInbox(accountID, inboxID int, updRequest *UpdateInboxRequest) (*Inbox, *Response, error)
	CleanInbox(accountID, inboxID int) (*Inbox, *Response, error)
	MarkAsRead(accountID, inboxID int) (*Inbox, *Response, error)
	ResetCredentials(accountID, inboxID int) (*Inbox, *Response, error)
	EnableEmail(accountID, inboxID int) (*Inbox, *Response, error)
	ResetEmail(accountID, inboxID int) (*Inbox, *Response, error)
	CreateInbox(accountID, inboxID int, name string) (*Inbox, *Response, error)
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
	LastMessageSentAt       string      `json:"last_message_sent_at"` // string or null
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

// CreateInbox create an inbox in a project.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/86631e73937e2-create-an-inbox
func (s *InboxesService) CreateInbox(accountID, inboxID int, name string) (*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(createInboxEndpoint, accountID, inboxID)
	payload := &createInboxRequest{
		Inbox: struct {
			Name string `json:"name"`
		}{Name: name},
	}

	return s.makeRequest(endpoint, http.MethodGet, payload)
}

type UpdateInboxRequest struct {
	Name          string `json:"name"`
	EmailUsername string `json:"email_username"`
}

// UpdateInbox update inbox name, inbox email username.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/768067eceee9d-update-an-inbox
func (s *InboxesService) UpdateInbox(accountID, inboxID int, updRequest *UpdateInboxRequest) (*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(updateInboxEndpoint, accountID, inboxID)
	payload := struct {
		Inbox *UpdateInboxRequest `json:"inbox"`
	}{updRequest}

	return s.makeRequest(endpoint, http.MethodPatch, payload)
}

// ListInboxes returns the list of inboxes.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/49dd3b9d6806f-get-a-list-of-inboxes
func (s *InboxesService) ListInboxes(accountID int) ([]*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(getInboxesEndpoint, accountID)
	req, err := s.client.NewRequest(http.MethodGet, endpoint, nil)
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

// GetInbox returns attributes of the inbox.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/432a39abe34b3-get-inbox-attributes
func (s *InboxesService) GetInbox(accountID, inboxID int) (*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(getInboxEndpoint, accountID, inboxID)
	return s.makeRequest(endpoint, http.MethodGet, nil)
}

// DeleteInbox delete an inbox with all its emails.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/e624770632299-delete-project
func (s *InboxesService) DeleteInbox(accountID, inboxID int) (*Response, error) {
	endpoint := fmt.Sprintf(deleteInboxEndpoint, accountID, inboxID)
	req, err := s.client.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// CleanInbox delete all messages (emails) from inbox.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/8a1e782a64fd0-clean-inbox
func (s *InboxesService) CleanInbox(accountID, inboxID int) (*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(cleanInboxEndpoint, accountID, inboxID)
	return s.makeRequest(endpoint, http.MethodPatch, nil)
}

// MarkAsRead mark all messages in the inbox as read.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/8a38b0494dff1-mark-as-read
func (s *InboxesService) MarkAsRead(accountID, inboxID int) (*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(markAsReadEndpoint, accountID, inboxID)
	return s.makeRequest(endpoint, http.MethodPatch, nil)
}

// ResetCredentials resets SMTP credentials of the inbox.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/403fd0f1315e6-reset-credentials
func (s *InboxesService) ResetCredentials(accountID, inboxID int) (*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(resetCredentialsEndpoint, accountID, inboxID)
	return s.makeRequest(endpoint, http.MethodPatch, nil)
}

// EnableEmail enables email address.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/a4b31a4c40ae4-enable-email-address
func (s *InboxesService) EnableEmail(accountID, inboxID int) (*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(enableEmailEndpoint, accountID, inboxID)
	return s.makeRequest(endpoint, http.MethodPatch, nil)
}

// ResetEmail reset email address
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/5ebb1ca46e3d0-reset-email-address
func (s *InboxesService) ResetEmail(accountID, inboxID int) (*Inbox, *Response, error) {
	endpoint := fmt.Sprintf(resetEmailEndpoint, accountID, inboxID)
	return s.makeRequest(endpoint, http.MethodPatch, nil)
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
