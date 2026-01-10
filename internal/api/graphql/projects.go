package graphql

import (
	"time"

	gql "github.com/shurcooL/graphql"
)

// ProjectV2 represents a GitHub Project v2
type ProjectV2 struct {
	CreatedAt   time.Time `graphql:"createdAt"`
	UpdatedAt   time.Time `graphql:"updatedAt"`
	Description *string   `graphql:"description"`
	Owner       struct {
		ID    string `graphql:"id"`
		Login string `graphql:"login"`
		Type  string `graphql:"__typename"`
	} `graphql:"owner"`
	ID     string `graphql:"id"`
	Title  string `graphql:"title"`
	URL    string `graphql:"url"`
	Fields struct {
		Nodes []ProjectV2Field `graphql:"nodes"`
	} `graphql:"fields(first: 20)"`
	Items struct {
		Nodes []ProjectV2Item `graphql:"nodes"`
	} `graphql:"items(first: 100)"`
	Number int  `graphql:"number"`
	Closed bool `graphql:"closed"`
}

// ProjectV2Field represents a custom field in a project
type ProjectV2Field struct {
	ID       string                 `graphql:"id"`
	Name     string                 `graphql:"name"`
	DataType ProjectV2FieldDataType `graphql:"dataType"`

	Options struct {
		Nodes []ProjectV2SingleSelectFieldOption `graphql:"nodes"`
	} `graphql:"... on ProjectV2SingleSelectField { options(first: 20) }"`
}

// ProjectV2FieldDataType represents the data type of a field
type ProjectV2FieldDataType string

const (
	ProjectV2FieldDataTypeText         ProjectV2FieldDataType = "TEXT"
	ProjectV2FieldDataTypeNumber       ProjectV2FieldDataType = "NUMBER"
	ProjectV2FieldDataTypeDate         ProjectV2FieldDataType = "DATE"
	ProjectV2FieldDataTypeSingleSelect ProjectV2FieldDataType = "SINGLE_SELECT"
	ProjectV2FieldDataTypeIteration    ProjectV2FieldDataType = "ITERATION"
)

// ProjectV2SingleSelectFieldOption represents an option for a single select field
type ProjectV2SingleSelectFieldOption struct {
	Description *string `graphql:"description"`
	ID          string  `graphql:"id"`
	Name        string  `graphql:"name"`
	Color       string  `graphql:"color"`
}

// ProjectV2Item represents an item in a project
type ProjectV2Item struct {
	CreatedAt   time.Time `graphql:"createdAt"`
	UpdatedAt   time.Time `graphql:"updatedAt"`
	ID          string    `graphql:"id"`
	FieldValues struct {
		Nodes []ProjectV2ItemFieldValue `graphql:"nodes"`
	} `graphql:"fieldValues(first: 20)"`
	Content struct {
		DraftBody   *string `graphql:"... on DraftIssue { body }"`
		PRTitle     string  `graphql:"... on PullRequest { title }"`
		IssueURL    string  `graphql:"... on Issue { url }"`
		IssueState  string  `graphql:"... on Issue { state }"`
		TypeName    string  `graphql:"__typename"`
		PRURL       string  `graphql:"... on PullRequest { url }"`
		PRState     string  `graphql:"... on PullRequest { state }"`
		DraftTitle  string  `graphql:"... on DraftIssue { title }"`
		IssueTitle  string  `graphql:"... on Issue { title }"`
		IssueNumber int     `graphql:"... on Issue { number }"`
		PRNumber    int     `graphql:"... on PullRequest { number }"`
		IssueClosed bool    `graphql:"... on Issue { closed }"`
		PRClosed    bool    `graphql:"... on PullRequest { closed }"`
	} `graphql:"content"`
}

// ProjectV2ItemFieldValue represents a field value for an item
type ProjectV2ItemFieldValue struct {
	TextValue         *string    `graphql:"... on ProjectV2ItemFieldTextValue { text }"`
	NumberValue       *float64   `graphql:"... on ProjectV2ItemFieldNumberValue { number }"`
	DateValue         *time.Time `graphql:"... on ProjectV2ItemFieldDateValue { date }"`
	SingleSelectValue *struct {
		ID   string `graphql:"id"`
		Name string `graphql:"name"`
	} `graphql:"... on ProjectV2ItemFieldSingleSelectValue { singleSelectOption }"`
	IterationValue *struct {
		ID    string `graphql:"id"`
		Title string `graphql:"title"`
	} `graphql:"... on ProjectV2ItemFieldIterationValue { iteration }"`
	Field struct {
		ID   string `graphql:"id"`
		Name string `graphql:"name"`
	} `graphql:"field"`
}

// Queries

// ListUserProjectsQuery lists projects for a user
type ListUserProjectsQuery struct {
	User struct {
		ProjectsV2 struct {
			PageInfo PageInfo    `graphql:"pageInfo"`
			Nodes    []ProjectV2 `graphql:"nodes"`
		} `graphql:"projectsV2(first: $first, after: $after)"`
	} `graphql:"user(login: $login)"`
}

// ListOrgProjectsQuery lists projects for an organization
type ListOrgProjectsQuery struct {
	Organization struct {
		ProjectsV2 struct {
			PageInfo PageInfo    `graphql:"pageInfo"`
			Nodes    []ProjectV2 `graphql:"nodes"`
		} `graphql:"projectsV2(first: $first, after: $after)"`
	} `graphql:"organization(login: $login)"`
}

// GetProjectQuery gets a specific project by number
type GetProjectQuery struct {
	Organization struct {
		ProjectV2 ProjectV2 `graphql:"projectV2(number: $number)"`
	} `graphql:"organization(login: $orgLogin)"`
}

// GetUserProjectQuery gets a specific user project by number
type GetUserProjectQuery struct {
	User struct {
		ProjectV2 ProjectV2 `graphql:"projectV2(number: $number)"`
	} `graphql:"user(login: $userLogin)"`
}

// PageInfo represents pagination information
type PageInfo struct {
	StartCursor     string `graphql:"startCursor"`
	EndCursor       string `graphql:"endCursor"`
	HasNextPage     bool   `graphql:"hasNextPage"`
	HasPreviousPage bool   `graphql:"hasPreviousPage"`
}

// Mutations

// CreateProjectMutation creates a new project
type CreateProjectMutation struct {
	CreateProjectV2 struct {
		ProjectV2 ProjectV2 `graphql:"projectV2"`
	} `graphql:"createProjectV2(input: $input)"`
}

// UpdateProjectMutation updates a project
type UpdateProjectMutation struct {
	UpdateProjectV2 struct {
		ProjectV2 ProjectV2 `graphql:"projectV2"`
	} `graphql:"updateProjectV2(input: $input)"`
}

// DeleteProjectMutation deletes a project
type DeleteProjectMutation struct {
	DeleteProjectV2 struct {
		ProjectV2 ProjectV2 `graphql:"projectV2"`
	} `graphql:"deleteProjectV2(input: $input)"`
}

// AddItemToProjectMutation adds an item to a project
type AddItemToProjectMutation struct {
	AddProjectV2ItemByID struct {
		Item ProjectV2Item `graphql:"item"`
	} `graphql:"addProjectV2ItemById(input: $input)"`
}

// UpdateItemFieldMutation updates a field value for an item
type UpdateItemFieldMutation struct {
	UpdateProjectV2ItemFieldValue struct {
		ProjectV2Item ProjectV2Item `graphql:"projectV2Item"`
	} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
}

// RemoveItemFromProjectMutation removes an item from a project
type RemoveItemFromProjectMutation struct {
	DeleteProjectV2Item struct {
		DeletedItemID string `graphql:"deletedItemId"`
	} `graphql:"deleteProjectV2Item(input: $input)"`
}

// Input Types

// CreateProjectInput represents input for creating a project
type CreateProjectInput struct {
	OwnerID     gql.ID      `json:"ownerId"`
	Title       gql.String  `json:"title"`
	Description *gql.String `json:"description,omitempty"`
	Readme      *gql.String `json:"readme,omitempty"`
	Visibility  *gql.String `json:"visibility,omitempty"`
	Repository  *gql.ID     `json:"repositoryId,omitempty"`
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	Title     *gql.String  `json:"title,omitempty"`
	Closed    *gql.Boolean `json:"closed,omitempty"`
	ProjectID gql.ID       `json:"projectId"`
}

// DeleteProjectInput represents input for deleting a project
type DeleteProjectInput struct {
	ProjectID gql.ID `json:"projectId"`
}

// AddItemInput represents input for adding an item to a project
type AddItemInput struct {
	ProjectID gql.ID `json:"projectId"`
	ContentID gql.ID `json:"contentId"`
}

// UpdateItemFieldInput represents input for updating an item field
type UpdateItemFieldInput struct {
	Value     interface{} `json:"value"`
	ProjectID gql.ID      `json:"projectId"`
	ItemID    gql.ID      `json:"itemId"`
	FieldID   gql.ID      `json:"fieldId"`
}

// RemoveItemInput represents input for removing an item from a project
type RemoveItemInput struct {
	ProjectID gql.ID `json:"projectId"`
	ItemID    gql.ID `json:"itemId"`
}

// Variable Builders

// BuildCreateProjectVariables builds variables for project creation
func BuildCreateProjectVariables(input *CreateProjectInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildUpdateProjectVariables builds variables for project update
func BuildUpdateProjectVariables(input *UpdateProjectInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildDeleteProjectVariables builds variables for project deletion
func BuildDeleteProjectVariables(input *DeleteProjectInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildAddItemVariables builds variables for adding an item
func BuildAddItemVariables(input *AddItemInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildUpdateItemFieldVariables builds variables for updating an item field
func BuildUpdateItemFieldVariables(input *UpdateItemFieldInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildRemoveItemVariables builds variables for removing an item
func BuildRemoveItemVariables(input *RemoveItemInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildListProjectsVariables builds variables for listing projects
func BuildListProjectsVariables(login string, first int, after *string) map[string]interface{} {
	vars := map[string]interface{}{
		"login": gql.String(login),
		"first": gql.Int(first), //nolint:gosec // first is always within int32 range
	}
	if after != nil {
		vars["after"] = gql.String(*after)
	} else {
		vars["after"] = (*gql.String)(nil)
	}
	return vars
}

// BuildGetProjectVariables builds variables for getting a project
func BuildGetProjectVariables(login string, number int, isOrg bool) map[string]interface{} {
	if isOrg {
		return map[string]interface{}{
			"orgLogin": gql.String(login),
			"number":   gql.Int(number), //nolint:gosec // number is always within int32 range
		}
	}
	return map[string]interface{}{
		"userLogin": gql.String(login),
		"number":    gql.Int(number), //nolint:gosec // number is always within int32 range
	}
}
