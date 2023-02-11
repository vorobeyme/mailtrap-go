package mailtrap

import (
	"fmt"
	"net/http"
)

const (
	getAccountUsersEndpoint    = "/accounts/%d/account_accesses"
	deleteAccountUsersEndpoint = "/accounts/%d/account_accesses/%d"
)

type AccountUsersServiceContract interface {
	ListAccountUsers(accountID int, params *ListAccountUsersParams) ([]*AccountUser, *Response, error)
	RemoveAccountUser(accountID, accountAccessID int) (*Response, error)
}

type AccountUsersService struct {
	client *Client
}

var _ AccountUsersServiceContract = &AccountUsersService{}

// Accounts represents a Mailtrap account.
type AccountUser struct {
	ID int `json:"id"`
	// specifier_type can return user, invite
	SpecifierType string                 `json:"specifier_type"`
	Resources     []accountUserResources `json:"resources"`
	Specifier     accountUserSpecifier   `json:"specifier"`
	Permissions   Permissions            `json:"permissions"`
}

// accountUserResources represents a Mailtrap account users resources.
type accountUserResources struct {
	// resource_type can return inbox, project, billing, account, mailsend_domain
	ResourceType string `json:"resource_type"`
	ResourceID   int    `json:"resource_id"`
	// access_level can return 1000 (account owner), 100 (admin), 10 (viewer).
	AccessLevel int `json:"access_level"`
}

// accountUserSpecifier represents a Mailtrap account users specifier.
type accountUserSpecifier struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ListAccountUsersParams represents the available ListAccountUsers() query parameters.
type ListAccountUsersParams struct {
	ProjectIDs *[]int    `url:"project_ids,omitempty" json:"project_ids,omitempty"`
	InboxIDs   *[]string `url:"inbox_ids,omitempty" json:"inbox_ids,omitempty"`
}

// ListAccountUsers get list of all account users.
//
// You need to have account admin or owner permissions for this endpoint to work.
// If you specify project_ids, inbox_ids, the endpoint returns users filtered by these resources.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/be0de9a48df49-list-all-users-in-account
func (s *AccountUsersService) ListAccountUsers(
	accountID int,
	params *ListAccountUsersParams,
) ([]*AccountUser, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf(getAccountUsersEndpoint, accountID), params)
	if err != nil {
		return nil, nil, err
	}

	var accUser []*AccountUser
	res, err := s.client.Do(req, &accUser)
	if err != nil {
		return nil, res, err
	}

	return accUser, res, err
}

// RemoveAccountUser remove user by their ID from the account.
// You need to be an account admin/owner for this endpoint to work.
//
// https://api-docs.mailtrap.io/docs/mailtrap-api-docs/569947e980f71-remove-user-from-the-account
func (s *AccountUsersService) RemoveAccountUser(accountID, accountAccessID int) (*Response, error) {
	uri := fmt.Sprintf(deleteAccountUsersEndpoint, accountID, accountAccessID)
	req, err := s.client.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
