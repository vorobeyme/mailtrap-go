package mailtrap

import "net/http"

type AccountsServiceContract interface {
	List() ([]*Account, *Response, error)
}

type AccountsService struct {
	client *client
}

var _ AccountsServiceContract = &AccountsService{}

// Account represents a Mailtrap account schema.
type Account struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	AccessLevels []int  `json:"access_levels"`
}

// List returns a list of Mailtrap accounts.
//
// See https://api-docs.mailtrap.io/docs/mailtrap-api-docs/4cfa4c61eae3c-get-all-accounts
func (s *AccountsService) List() ([]*Account, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/accounts", nil)
	if err != nil {
		return nil, nil, err
	}

	var accounts []*Account
	res, err := s.client.Do(req, &accounts)
	if err != nil {
		return nil, res, err
	}

	return accounts, res, nil
}
