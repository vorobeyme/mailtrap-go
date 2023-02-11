package mailtrap

import (
	"fmt"
	"net/http"
)

const (
	forwardMessageEndpoint      = "/accounts/%d/inboxes/%d/messages/%d/forward"
	updateMessageEndpoint       = "/accounts/%d/inboxes/%d/messages/%d/"
	deleteMessageEndpoint       = "/accounts/%d/inboxes/%d/messages/%d/"
	showMessageEndpoint         = "/accounts/%d/inboxes/%d/messages/%d"
	listMessagesEndpoint        = "/accounts/%d/inboxes/%d/messages"
	getMessageSpamScoreEndpoint = "/accounts/%d/inboxes/%d/messages/%d/spam_report"
	getTextMessageEndpoint      = "/accounts/%d/inboxes/%d/messages/%d/body.txt"
	getRawMessageEndpoint       = "/accounts/%d/inboxes/%d/messages/%d/body.raw"
	getSourceMessageEndpoint    = "/accounts/%d/inboxes/%d/messages/%d/body.htmlsource"
	getHTMLMessageEndpoint      = "/accounts/%d/inboxes/%d/messages/%d/body.html"
	getEmlMessageEndpoint       = "/accounts/%d/inboxes/%d/messages/%d/body.eml"
)

type MessagesServiceContract interface {
	ShowMessage(accountID, inboxID, messageID int) (*Message, *Response, error)
	UpdateMessage(accountID, inboxID, messageID int, updRequest *UpdateMessageRequest) (*Message, *Response, error)
	DeleteMessage(accountID, inboxID, messageID int) (*Message, *Response, error)
	ListMessages(accountID, inboxID int) (*[]Message, *Response, error)
	ForwardMessage(accountID, inboxID, messageID int, email string) (*Response, error)
	GetMessageSpamSource(accountID, inboxID, messageID int) (*MessageSpamReport, *Response, error)
	GetTextMessage(accountID, inboxID, messageID int) (*Message, *Response, error)
	GetRawMessage(accountID, inboxID, messageID int) (*Message, *Response, error)
	GetMessageSource(accountID, inboxID, messageID int) (*Message, *Response, error)
	GetHTMLMessage(accountID, inboxID, messageID int) (*Message, *Response, error)
	GetMessageAsEml(accountID, inboxID, messageID int) (*Message, *Response, error)
}

type MessagesService struct {
	client *Client
}

var _ MessagesServiceContract = &MessagesService{}

// Message represents a Mailtrap message.
type Message struct {
	ID                   int    `json:"id"`
	InboxID              int    `json:"inbox_id"`
	Subject              string `json:"subject"`
	SentAt               string `json:"sent_at"`
	FromEmail            string `json:"from_email"`
	FromName             string `json:"from_name"`
	ToEmail              string `json:"to_email"`
	ToName               string `json:"to_name"`
	EmailSize            int    `json:"email_size"`
	IsRead               bool   `json:"is_read"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	HtmlBodySize         int    `json:"html_body_size"`
	TextBodySize         int    `json:"text_body_size"`
	HumanSize            string `json:"human_size"`
	HTMLPath             string `json:"html_path"`
	TxtPath              string `json:"txt_path"`
	RawPath              string `json:"raw_path"`
	DownloadPath         string `json:"download_path"`
	HTMLSourcePath       string `json:"html_source_path"`
	BlacklistsReportInfo bool   `json:"blacklists_report_info"`
	SMTPInformation      struct {
		Ok   bool `json:"ok"`
		Data struct {
			MailFromAddr string `json:"mail_from_addr"`
			ClientIP     string `json:"client_ip"`
		} `json:"data"`
	} `json:"smtp_information"`
}

// MessageSpamReport represents Mailtrap message spam analysis report.
type MessageSpamReport struct {
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

// ShowMessage returns email message with its attributes by ID.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/c1708cf554d6e-show-email-message
func (s *MessagesService) ShowMessage(accountID, inboxID, messageID int) (*Message, *Response, error) {
	endpoint := fmt.Sprintf(showMessageEndpoint, accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodGet, endpoint, nil)
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
	IsRead string `json:"is_read"`
}

// UpdateMessage updates message attributes
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) UpdateMessage(
	accountID, inboxID, messageID int,
	updRequest *UpdateMessageRequest,
) (*Message, *Response, error) {
	endpoint := fmt.Sprintf(updateMessageEndpoint, accountID, inboxID, messageID)
	payload := struct {
		Message *UpdateMessageRequest `json:"message"`
	}{updRequest}

	req, err := s.client.NewRequest(http.MethodPatch, endpoint, payload)
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

// DeleteMessage delete message from inbox.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) DeleteMessage(accountID, inboxID, messageID int) (*Message, *Response, error) {
	endpoint := fmt.Sprintf(deleteMessageEndpoint, accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodDelete, endpoint, nil)
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

// ListMessages returns all messages in inboxs.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) ListMessages(accountID, inboxID int) (*[]Message, *Response, error) {
	endpoint := fmt.Sprintf(listMessagesEndpoint, accountID, inboxID)
	req, err := s.client.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var msg *[]Message
	res, err := s.client.Do(req, &msg)
	if err != nil {
		return nil, res, err
	}

	return msg, res, nil
}

// ForwardMessage forward message to an email address.
// The email address must be confirmed by the recipient in advance.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) ForwardMessage(
	accountID, inboxID, messageID int,
	email string,
) (*Response, error) {
	endpoint := fmt.Sprintf(forwardMessageEndpoint, accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetMessageSpamSource returns a brief spam report by message ID.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/000f54556fc6e-get-message-spam-score
func (s *MessagesService) GetMessageSpamSource(
	accountID, inboxID, messageID int,
) (*MessageSpamReport, *Response, error) {
	endpoint := fmt.Sprintf(getMessageSpamScoreEndpoint, accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var report *MessageSpamReport
	res, err := s.client.Do(req, &report)
	if err != nil {
		return nil, res, err
	}

	return report, res, nil
}

// GetTextMessage returns text email body, if it exists.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) GetTextMessage(accountID, inboxID, messageID int) (*Message, *Response, error) {
	endpoint := fmt.Sprintf(getTextMessageEndpoint, accountID, inboxID, messageID)
	return s.makeRequest(endpoint, http.MethodGet)
}

// GetRawMessage returns raw email body.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) GetRawMessage(accountID, inboxID, messageID int) (*Message, *Response, error) {
	endpoint := fmt.Sprintf(getRawMessageEndpoint, accountID, inboxID, messageID)
	return s.makeRequest(endpoint, http.MethodGet)
}

// GetMessageSource returns HTML source of email.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) GetMessageSource(accountID, inboxID, messageID int) (*Message, *Response, error) {
	endpoint := fmt.Sprintf(getSourceMessageEndpoint, accountID, inboxID, messageID)
	return s.makeRequest(endpoint, http.MethodGet)
}

// GetHTMLMessage returns formatted HTML email body. Not applicable for plain text emails.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) GetHTMLMessage(accountID, inboxID, messageID int) (*Message, *Response, error) {
	endpoint := fmt.Sprintf(getHTMLMessageEndpoint, accountID, inboxID, messageID)
	return s.makeRequest(endpoint, http.MethodGet)
}

// GetMessageAsEml returns email message in .eml format.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/53cf46462fba5-update-message
func (s *MessagesService) GetMessageAsEml(accountID, inboxID, messageID int) (*Message, *Response, error) {
	endpoint := fmt.Sprintf(getEmlMessageEndpoint, accountID, inboxID, messageID)
	return s.makeRequest(endpoint, http.MethodGet)
}

func (s *MessagesService) makeRequest(endpoint, httpMethod string) (*Message, *Response, error) {
	req, err := s.client.NewRequest(httpMethod, endpoint, nil)
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
