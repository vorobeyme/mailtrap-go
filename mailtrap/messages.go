package mailtrap

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"
)

type MessagesServiceContract interface {
	List(accountID, inboxID int) ([]*Message, *Response, error)
	Get(accountID, inboxID, messageID int) (*Message, *Response, error)
	Update(accountID, inboxID, messageID int, updateReq *UpdateMessageRequest) (*Message, *Response, error)
	Delete(accountID, inboxID, messageID int) (*Response, error)
	Forward(accountID, inboxID, messageID int, email string) (*Response, error)
	SpamReport(accountID, inboxID, messageID int) (*SpamReport, *Response, error)
	AsRaw(accountID, inboxID, messageID int) (string, *Response, error)
	AsText(accountID, inboxID, messageID int) (string, *Response, error)
	AsHTML(accountID, inboxID, messageID int) (string, *Response, error)
	AsHTMLSource(accountID, inboxID, messageID int) (string, *Response, error)
	AsEML(accountID, inboxID, messageID int) (string, *Response, error)
}

type MessagesService struct {
	client *Client
}

var _ MessagesServiceContract = &MessagesService{}

// Message represents a Mailtrap message.
type Message struct {
	ID                   int              `json:"id"`
	InboxID              int              `json:"inbox_id"`
	Subject              string           `json:"subject"`
	SentAt               time.Time        `json:"sent_at"`
	FromEmail            string           `json:"from_email"`
	FromName             string           `json:"from_name"`
	ToEmail              string           `json:"to_email"`
	ToName               string           `json:"to_name"`
	EmailSize            int              `json:"email_size"`
	IsRead               bool             `json:"is_read"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
	HTMLBodySize         int              `json:"html_body_size"`
	TextBodySize         int              `json:"text_body_size"`
	HumanSize            string           `json:"human_size"`
	HTMLPath             string           `json:"html_path"`
	TxtPath              string           `json:"txt_path"`
	RawPath              string           `json:"raw_path"`
	DownloadPath         string           `json:"download_path"`
	HTMLSourcePath       string           `json:"html_source_path"`
	BlacklistsReportInfo bool             `json:"blacklists_report_info"`
	SMTPInfo             *MessageSMTPInfo `json:"smtp_information"`
}

// MessageSMTPInfo represents a Mailtrap message SMTP information.
type MessageSMTPInfo struct {
	Ok   bool `json:"ok"`
	Data struct {
		MailFromAddr string `json:"mail_from_addr"`
		ClientIP     string `json:"client_ip"`
	} `json:"data"`
}

// SpamReport represents Mailtrap message spam analysis report.
type SpamReport struct {
	Report struct {
		ResponseCode    int           `json:"ResponseCode"`
		ResponseMessage string        `json:"ResponseMessage"`
		ResponseVersion string        `json:"ResponseVersion"`
		Score           float64       `json:"Score"`
		Spam            bool          `json:"Spam"`
		Threshold       float64       `json:"Threshold"`
		Details         []interface{} `json:"Details"` // string or null
	} `json:"report"`
}

// List returns all messages in inboxs.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/a80869adf4489-get-messages
func (s *MessagesService) List(accountID, inboxID int) ([]*Message, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages", accountID, inboxID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var msg []*Message
	res, err := s.client.Do(req, &msg)
	if err != nil {
		return nil, res, err
	}

	return msg, res, nil
}

// Get returns email message with its attributes by ID.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/c1708cf554d6e-show-email-message
func (s *MessagesService) Get(accountID, inboxID, messageID int) (*Message, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d", accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var msg *Message
	res, err := s.client.Do(req, &msg)
	if err != nil {
		return nil, res, err
	}

	return msg, res, err
}

type UpdateMessageRequest struct {
	IsRead bool `json:"is_read"`
}

// Update updates message attributes.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) Update(
	accountID, inboxID, messageID int,
	updateReq *UpdateMessageRequest,
) (*Message, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d", accountID, inboxID, messageID)
	payload := struct {
		Message *UpdateMessageRequest `json:"message"`
	}{updateReq}

	req, err := s.client.NewRequest(http.MethodPatch, u, payload)
	if err != nil {
		return nil, nil, err
	}

	var msg *Message
	res, err := s.client.Do(req, &msg)
	if err != nil {
		return nil, res, err
	}

	return msg, res, nil
}

// Delete removes message from inbox.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) Delete(accountID, inboxID, messageID int) (*Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d", accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

type forwardRequest struct {
	Email string `json:"email"`
}

// Forward forwards message to an email address.
// The email address must be confirmed by the recipient in advance.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) Forward(
	accountID, inboxID, messageID int,
	email string,
) (*Response, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("Forward 'email' is invalid.")
	}

	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/forward", accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodPost, u, &forwardRequest{Email: email})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// SpamReport returns a brief spam report by message ID.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/000f54556fc6e-get-message-spam-score
func (s *MessagesService) SpamReport(accountID, inboxID, messageID int) (*SpamReport, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/spam_report", accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var report *SpamReport
	res, err := s.client.Do(req, &report)
	if err != nil {
		return nil, res, err
	}

	return report, res, nil
}

// AsRaw returns raw email body.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) AsRaw(accountID, inboxID, messageID int) (string, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/body.raw", accountID, inboxID, messageID)
	return s.makeRequest(u, http.MethodGet, "text/plain")
}

// AsText returns text email body, if it exists.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) AsText(accountID, inboxID, messageID int) (string, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/body.txt", accountID, inboxID, messageID)
	return s.makeRequest(u, http.MethodGet, "text/plain")
}

// AsHTML returns formatted HTML email body. Not applicable for plain text emails.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) AsHTML(accountID, inboxID, messageID int) (string, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/body.html", accountID, inboxID, messageID)
	return s.makeRequest(u, http.MethodGet, "text/html")
}

// AsHTMLSource returns HTML source of email.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) AsHTMLSource(accountID, inboxID, messageID int) (string, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/body.htmlsource", accountID, inboxID, messageID)
	return s.makeRequest(u, http.MethodGet, "text/html")
}

// AsEML returns email message in .eml format.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) AsEML(accountID, inboxID, messageID int) (string, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/body.eml", accountID, inboxID, messageID)
	return s.makeRequest(u, http.MethodGet, "message/rfc822")
}

func (s *MessagesService) makeRequest(endpoint, httpMethod string, acceptHeader string) (string, *Response, error) {
	req, err := s.client.NewRequest(httpMethod, endpoint, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Accept", acceptHeader)

	var respStr string
	res, err := s.client.Do(req, &respStr)
	if err != nil {
		return "", res, err
	}

	return respStr, res, nil
}
