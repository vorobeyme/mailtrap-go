package mailtrap

import (
	"fmt"
	"net/http"
)

type AttachmentsServiceContract interface {
	List(accountID, inboxID, messageID int) ([]*Attachment, *Response, error)
	Get(accountID, inboxID, messageID, attachmentID int) (*Attachment, *Response, error)
}

type AttachmentsService struct {
	client *Client
}

var _ AttachmentsServiceContract = &AttachmentsService{}

// Attachment represents a Mailtrap attachment schema.
type Attachment struct {
	ID                  int    `json:"id"`
	MessageID           int    `json:"message_id"`
	Filename            string `json:"filename"`
	AttachmentType      string `json:"attachment_type"`
	ContentType         string `json:"content_type"`
	ContentID           string `json:"content_id"`
	TransferEncoding    string `json:"transfer_encoding"`
	AttachmentSize      int    `json:"attachment_size"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	AttachmentHumanSize string `json:"attachment_human_size"`
	DownloadPath        string `json:"download_path"`
}

// List returns message attachments by inboxID and messageID.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/bcb1ef001e32d-get-attachments
func (s *AttachmentsService) List(
	accountID, inboxID, messageID int,
) ([]*Attachment, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/attachments", accountID, inboxID, messageID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var attach []*Attachment
	resp, err := s.client.Do(req, &attach)
	if err != nil {
		return nil, resp, err
	}

	return attach, resp, err
}

// Get returns message single attachment by ID.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/e2e15ad4475a4-get-single-attachment
func (s *AttachmentsService) Get(
	accountID, inboxID, messageID, attachmentID int,
) (*Attachment, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/inboxes/%d/messages/%d/attachments/%d", accountID, inboxID, messageID, attachmentID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var attach *Attachment
	res, err := s.client.Do(req, &attach)
	if err != nil {
		return nil, res, err
	}

	return attach, res, err
}
