package mailtrap

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const sendEmailEndpoint = "/send"

// SendEmailServiceContract is an interface for interfacing with the email
// sending endpoints of the Mailtrap API
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/67f1d70aeb62c-send-email
type SendEmailServiceContract interface {
	Send(request *SendEmailRequest) (*SendEmailResponse, *Response, error)
}

// SendEmailService handles communication with the email sending API.
type SendEmailService struct {
	client *Client
}

var _ SendEmailServiceContract = &SendEmailService{}

// SendEmailRequest represents the request to send email.
type SendEmailRequest struct {
	From EmailAddress   `json:"from"` // required
	To   []EmailAddress `json:"to"`   // required
	Cc   []EmailAddress `json:"cc"`
	Bcc  []EmailAddress `json:"bcc"`

	// An array of objects where you can specify any attachments you want to include.
	Attachments []EmailAttachment `json:"attachments"`

	// An object containing key/value pairs of header names and the value to substitute for them.
	// The key/value pairs must be strings.
	// You must ensure these are properly encoded if they contain unicode characters.
	// These headers cannot be one of the reserved headers.
	Headers map[string]string `json:"headers"`

	// Values that are specific to the entire send that will be carried along with the email and its activity data.
	// Total size of custom variables in JSON form must not exceed 1000 bytes.
	CustomVars map[string]string `json:"custom_variables"`

	// The global or 'message level' subject of your email.
	// This may be overridden by subject lines set in personalizations.
	// required
	Subject string `json:"subject"`

	// Text version of the body of the email. Can be used along with html to create a fallback for non-html clients.
	// Required in the absence of html.
	Text string `json:"text"`

	// HTML version of the body of the email. Can be used along with text to create a fallback for non-html clients.
	// Required in the absence of text.
	HTML     string `json:"html"`
	Category string `json:"category"`
}

// EmailAddress represents an email address.
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// EmailAttachment represents an email attacments.
type EmailAttachment struct {
	// The Base64 encoded content of the attachment.
	// required
	Content string `json:"content"`
	// The MIME type of the content you are attaching (e.g., “text/plain” or “text/html”).
	AttachType string `json:"type"`

	// The attachment's filename.
	// required
	Filename string `json:"filename"`

	// The attachment's content-disposition, specifying how you would like the attachment to be displayed.
	// For example, “inline” results in the attached file are displayed automatically within the message
	// while “attachment” results in the attached file require some action to be taken before it is displayed,
	// such as opening or downloading the file.
	//
	// Allowed values: inline, attachment
	// Default: attachment
	Disposition string `json:"disposition"`

	// The attachment's content ID.
	// This is used when the disposition is set to “inline” and the attachment is an image,
	// allowing the file to be displayed within the body of your email.
	ContentID string `json:"content_id"`
}

// SendEmailResponse contains response from email sending API.
type SendEmailResponse struct {
	Success    bool     `json:"success"`
	MessageIDs []string `json:"message_ids"`
}

// Send email
func (s *SendEmailService) Send(request *SendEmailRequest) (*SendEmailResponse, *Response, error) {
	if request == nil {
		return nil, nil, errors.New("request `SendEmailRequest` to send mail is mandatory")
	}

	if err := request.validate(); err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, sendEmailEndpoint, request)
	if err != nil {
		return nil, nil, err
	}

	response := new(SendEmailResponse)
	res, err := s.client.Do(req, response)
	if err != nil {
		return nil, res, err
	}

	return response, res, err
}

// Send email request validation
func (r *SendEmailRequest) validate() error {
	if r.From.Email == "" {
		return errors.New("'from' address is required")
	}

	if len(r.To) == 0 {
		return errors.New("'to' address is required")
	}
	for _, v := range r.To {
		if v.Email == "" {
			return errors.New("'email' is required in 'to' address")
		}
	}

	if len(r.Attachments) > 0 {
		var errMsg []string
		for _, v := range r.Attachments {
			if v.Content == "" {
				errMsg = append(errMsg, "'content' is required in attachment")
			}
			if v.Filename == "" {
				errMsg = append(errMsg, "'filename' is required in attachment")
			}
		}
		if len(errMsg) > 0 {
			return errors.New(strings.Join(errMsg, "; "))
		}
	}

	if r.Subject == "" {
		return errors.New("'subject' is required")
	}

	if r.Text == "" && r.HTML == "" {
		return errors.New("one of 'text' or 'html' is required")
	}

	const categoryMaxLength int = 255
	if len(r.Category) > categoryMaxLength {
		return fmt.Errorf("'category' is greater than %d chars", categoryMaxLength)
	}

	return nil
}
