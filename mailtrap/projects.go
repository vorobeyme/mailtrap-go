package mailtrap

import (
	"fmt"
	"net/http"
)

// ProjectsServiceContract defines the methods available to projects.
type ProjectsServiceContract interface {
	List(accountID int) ([]*Project, *Response, error)
	Get(accountID, projectID int) (*Project, *Response, error)
	Create(accountID int, name string) (*Project, *Response, error)
	Update(accountID, projectID int, name string) (*Project, *Response, error)
	Delete(accountID, projectID int) (*Response, error)
}

type ProjectsService struct {
	client *client
}

var _ ProjectsServiceContract = &ProjectsService{}

// Project represents a Mailtrap project.
type Project struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ShareLinks struct {
		Admin  string `json:"admin"`
		Viewer string `json:"viewer"`
	} `json:"share_links"`
	Inboxes     []Inbox     `json:"inboxes"`
	Permissions Permissions `json:"permissions"`
}

// ProjectRequest represents the request to create / update project.
type projectRequest struct {
	Project struct {
		Name string `json:"name"`
	} `json:"project"`
}

// List returns the list of projects and their inboxes.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/c088109b11d07-get-a-list-of-projects
func (s *ProjectsService) List(accountID int) ([]*Project, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/projects", accountID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var project []*Project
	res, err := s.client.Do(req, &project)
	if err != nil {
		return nil, res, err
	}

	return project, res, err
}

// Get returns the project and its inboxes.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/3c60381e63410-get-project-by-id
func (s *ProjectsService) Get(accountID, projectID int) (*Project, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/projects/%d", accountID, projectID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var project *Project
	res, err := s.client.Do(req, &project)
	if err != nil {
		return nil, res, err
	}

	return project, res, err
}

// Delete removes project and its inboxes.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/e624770632299-delete-project
func (s *ProjectsService) Delete(accountID, projectID int) (*Response, error) {
	u := fmt.Sprintf("/accounts/%d/projects/%d", accountID, projectID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Update updates project name.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/73bdfaac8c86c-update-project
func (s *ProjectsService) Update(accountID, projectID int, name string) (*Project, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/projects/%d", accountID, projectID)
	payload := &projectRequest{
		Project: struct {
			Name string `json:"name"`
		}{Name: name},
	}

	req, err := s.client.NewRequest(http.MethodPatch, u, payload)
	if err != nil {
		return nil, nil, err
	}

	var project *Project
	res, err := s.client.Do(req, &project)
	if err != nil {
		return nil, res, err
	}

	return project, res, err
}

// Create creates a Mailtrap project.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/ee252e413d78a-create-project
func (s *ProjectsService) Create(accountID int, name string) (*Project, *Response, error) {
	u := fmt.Sprintf("/accounts/%d/projects", accountID)
	payload := &projectRequest{
		Project: struct {
			Name string `json:"name"`
		}{Name: name},
	}

	req, err := s.client.NewRequest(http.MethodPost, u, payload)
	if err != nil {
		return nil, nil, err
	}

	var project *Project
	res, err := s.client.Do(req, &project)
	if err != nil {
		return nil, res, err
	}

	return project, res, err
}
