package mailtrap

import (
	"fmt"
	"net/http"
)

type PermissionsServiceContract interface {
	ListResources(accountID int) ([]*Resource, *Response, error)
	Manage(accountID, accountAccessID int, permissionReq *[]PermissionRequest) (*Response, error)
}

type PermissionsService struct {
	client *Client
}

var _ PermissionsServiceContract = &PermissionsService{}

// Resource represents the resources nested according to their hierarchy.
type Resource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	// AccessLevel represents the access level of the token used to make the request.
	AccessLevel int        `json:"access_level"`
	Resource    []*Resource `json:"resources"`
}

// Permissions represents a Mailtrap permissions schema.
type Permissions struct {
	CanRead    bool `json:"can_read"`
	CanUpdate  bool `json:"can_update"`
	CanDestroy bool `json:"can_destroy"`
	CanLeave   bool `json:"can_leave"`
}

type PermissionRequest struct {
	// ResourceID is an ID of the resource.
	ResourceID int `json:"resource_id,omitempty"`

	// ResourceType can be account, billing, project, inbox or mailsend_domain.
	ResourceType string `json:"resource_type,omitempty"`

	// AccessLevel can be admin or viewer or their numbers 100 and 10 respectively
	AccessLevel string `json:"access_level,omitempty"`

	// Destroy - if true, instead of creating/updating the permission, it destroys it
	Destroy bool `json:"_destroy,omitempty"`
}

// permissionRequest represents the request to manage permissions.
type permissionRequest struct {
	Permissions *[]PermissionRequest `json:"permissions"`
}

// GetResources returns all resources in account to which the token has admin access.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/595e78d9c870b-get-resources
func (s *PermissionsService) ListResources(accountID int) ([]*Resource, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/permissions/resources", accountID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var resource []*Resource
	res, err := s.client.Do(req, &resource)
	if err != nil {
		return nil, res, err
	}

	return resource, res, err
}

// Manage manages user or token permissions.
//
// If send a combination of resource_type and resource_id that already exists, the permission is updated.
// If the combination doesn’t exist, the permission is created.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/78d0b711767ea-manage-user-or-token-permissions
func (s *PermissionsService) Manage(
	accountID, accountAccessID int,
	permissionReq *[]PermissionRequest,
) (*Response, error) {
	u := fmt.Sprintf("/accounts/%d/account_accesses/%d/permissions/bulk", accountID, accountAccessID)
	req, err := s.client.NewRequest(http.MethodPut, u, &permissionRequest{Permissions: permissionReq})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
