package mailtrap

import (
	"fmt"
	"net/http"
)

type AccountUsersServiceContract interface {
	List(accountID int, params *ListAccountUsersParams) ([]*AccountUser, *Response, error)
	Delete(accountID, accountAccessID int) (*Response, error)
}

type AccountUsersService struct {
	client *client
}

var _ AccountUsersServiceContract = &AccountUsersService{}

// AccountUsers represents a Mailtrap account users.
type AccountUser struct {
	ID int `json:"id"`
	// specifier_type can return user, invite
	SpecifierType string                 `json:"specifier_type"`
	Resources     []AccountUserResources `json:"resources"`
	Specifier     AccountUserSpecifier   `json:"specifier"`
	Permissions   Permissions            `json:"permissions"`
}

// AccountUserResources represents a Mailtrap account users resources.
type AccountUserResources struct {
	// resource_type can return inbox, project, billing, account, mailsend_domain
	ResourceType string `json:"resource_type"`
	ResourceID   int    `json:"resource_id"`
	// access_level can return 1000 (account owner), 100 (admin), 10 (viewer).
	AccessLevel int `json:"access_level"`
}

// AccountUserSpecifier represents a Mailtrap account users specifier.
type AccountUserSpecifier struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ListAccountUsersParams represents the available List() query parameters.
type ListAccountUsersParams struct {
	ProjectIDs *[]int    `url:"project_ids,omitempty" json:"project_ids,omitempty"`
	InboxIDs   *[]string `url:"inbox_ids,omitempty" json:"inbox_ids,omitempty"`
}

// List returns list of all account users.
//
// You need to have account admin or owner permissions for this endpoint to work.
// If you specify project_ids, inbox_ids, the endpoint returns users filtered by these resources.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/be0de9a48df49-list-all-users-in-account
func (s *AccountUsersService) List(
	accountID int,
	params *ListAccountUsersParams,
) ([]*AccountUser, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/account_accesses", accountID)
	req, err := s.client.NewRequest(http.MethodGet, u, params)
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

// Delete removes user by their ID from the account.
// You need to be an account admin/owner for this endpoint to work.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/569947e980f71-remove-user-from-the-account
func (s *AccountUsersService) Delete(accountID, accountAccessID int) (*Response, error) {
	u := fmt.Sprintf("/accounts/%d/account_accesses/%d", accountID, accountAccessID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
