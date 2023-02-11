package mailtrap

import (
	"fmt"
	"net/http"
)

const (
	createProjectEndpoint = "/accounts/%d/projects"
	getProjectsEndpoint   = "/accounts/%d/projects"
	getProjectEndpoint    = "/accounts/%d/projects/%d"
	updateProjectEndpoint = "/accounts/%d/projects/%d"
	deleteProjectEndpoint = "/accounts/%d/projects/%d"
)

type ProjectsServiceContract interface {
	CreateProject(accountID int, name string) (*Project, *Response, error)
	ListProjects(accountID int) ([]*Project, *Response, error)
	GetProject(accountID, projectID int) (*Project, *Response, error)
	UpdateProject(accountID, projectID int, name string) (*Project, *Response, error)
	DeleteProject(accountID, projectID int) (*Response, error)
}

type ProjectsService struct {
	client *Client
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
type ProjectRequest struct {
	Project struct {
		Name string `json:"name"`
	} `json:"project"`
}

// ListProjects list projects and their inboxes.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/c088109b11d07-get-a-list-of-projects
func (s *ProjectsService) ListProjects(accountID int) ([]*Project, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, fmt.Sprintf(getProjectsEndpoint, accountID), nil)
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

// GetProject get the project and its inboxes.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/3c60381e63410-get-project-by-id
func (s *ProjectsService) GetProject(accountID, projectID int) (*Project, *Response, error) {
	endpoint := fmt.Sprintf(getProjectEndpoint, accountID, projectID)
	req, err := s.client.NewRequest(http.MethodGet, endpoint, nil)
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

// DeleteProject delete project and its inboxes.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/e624770632299-delete-project
func (s *ProjectsService) DeleteProject(accountID, projectID int) (*Response, error) {
	endpoint := fmt.Sprintf(deleteProjectEndpoint, accountID, projectID)
	req, err := s.client.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UpdateProject update project name.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/73bdfaac8c86c-update-project
func (s *ProjectsService) UpdateProject(accountID, projectID int, name string) (*Project, *Response, error) {
	endpoint := fmt.Sprintf(updateProjectEndpoint, accountID, projectID)
	updateReq := &ProjectRequest{
		Project: struct {
			Name string `json:"name"`
		}{Name: name},
	}

	req, err := s.client.NewRequest(http.MethodPatch, endpoint, updateReq)
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

// CreateProject creates project.
//
// See: https://api-docs.mailtrap.io/docs/mailtrap-api-docs/ee252e413d78a-create-project
func (s *ProjectsService) CreateProject(accountID int, name string) (*Project, *Response, error) {
	endpoint := fmt.Sprintf(createProjectEndpoint, accountID)
	createReq := &ProjectRequest{
		Project: struct {
			Name string `json:"name"`
		}{Name: name},
	}

	req, err := s.client.NewRequest(http.MethodPost, endpoint, createReq)
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
