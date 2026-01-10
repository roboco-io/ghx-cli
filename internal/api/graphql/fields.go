package graphql

import (
	"time"

	gql "github.com/shurcooL/graphql"
)

// Field creation mutations and queries

// CreateFieldMutation represents the createProjectV2Field mutation
type CreateFieldMutation struct {
	CreateProjectV2Field struct {
		ProjectV2Field ProjectV2Field `graphql:"projectV2Field"`
	} `graphql:"createProjectV2Field(input: $input)"`
}

// UpdateFieldMutation represents the updateProjectV2Field mutation
type UpdateFieldMutation struct {
	UpdateProjectV2Field struct {
		ProjectV2Field ProjectV2Field `graphql:"projectV2Field"`
	} `graphql:"updateProjectV2Field(input: $input)"`
}

// DeleteFieldMutation represents the deleteProjectV2Field mutation
type DeleteFieldMutation struct {
	DeleteProjectV2Field struct {
		ProjectV2Field ProjectV2Field `graphql:"projectV2Field"`
	} `graphql:"deleteProjectV2Field(input: $input)"`
}

// CreateSingleSelectFieldOptionMutation represents the createProjectV2SingleSelectFieldOption mutation
type CreateSingleSelectFieldOptionMutation struct {
	CreateProjectV2SingleSelectFieldOption struct {
		ProjectV2SingleSelectFieldOption ProjectV2SingleSelectFieldOption `graphql:"projectV2SingleSelectFieldOption"`
	} `graphql:"createProjectV2SingleSelectFieldOption(input: $input)"`
}

// UpdateSingleSelectFieldOptionMutation represents the updateProjectV2SingleSelectFieldOption mutation
type UpdateSingleSelectFieldOptionMutation struct {
	UpdateProjectV2SingleSelectFieldOption struct {
		ProjectV2SingleSelectFieldOption ProjectV2SingleSelectFieldOption `graphql:"projectV2SingleSelectFieldOption"`
	} `graphql:"updateProjectV2SingleSelectFieldOption(input: $input)"`
}

// DeleteSingleSelectFieldOptionMutation represents the deleteProjectV2SingleSelectFieldOption mutation
type DeleteSingleSelectFieldOptionMutation struct {
	DeleteProjectV2SingleSelectFieldOption struct {
		ProjectV2SingleSelectFieldOption ProjectV2SingleSelectFieldOption `graphql:"projectV2SingleSelectFieldOption"`
	} `graphql:"deleteProjectV2SingleSelectFieldOption(input: $input)"`
}

// Field input types

// CreateFieldInput represents input for creating a field
type CreateFieldInput struct {
	ProjectID           gql.ID                 `json:"projectId"`
	Name                gql.String             `json:"name"`
	DataType            ProjectV2FieldDataType `json:"dataType"`
	SingleSelectOptions []SingleSelectOption   `json:"singleSelectOptions,omitempty"`
}

// SingleSelectOption represents a single select option for field creation
type SingleSelectOption struct {
	Name        gql.String  `json:"name"`
	Color       gql.String  `json:"color"`
	Description *gql.String `json:"description,omitempty"`
}

// UpdateFieldInput represents input for updating a field
type UpdateFieldInput struct {
	Name    *gql.String `json:"name,omitempty"`
	FieldID gql.ID      `json:"fieldId"`
}

// DeleteFieldInput represents input for deleting a field
type DeleteFieldInput struct {
	FieldID gql.ID `json:"fieldId"`
}

// CreateSingleSelectFieldOptionInput represents input for creating a single select option
type CreateSingleSelectFieldOptionInput struct {
	FieldID     gql.ID      `json:"fieldId"`
	Name        gql.String  `json:"name"`
	Color       gql.String  `json:"color"`
	Description *gql.String `json:"description,omitempty"`
}

// UpdateSingleSelectFieldOptionInput represents input for updating a single select option
type UpdateSingleSelectFieldOptionInput struct {
	Name        *gql.String `json:"name,omitempty"`
	Color       *gql.String `json:"color,omitempty"`
	Description *gql.String `json:"description,omitempty"`
	OptionID    gql.ID      `json:"singleSelectOptionId"`
}

// DeleteSingleSelectFieldOptionInput represents input for deleting a single select option
type DeleteSingleSelectFieldOptionInput struct {
	OptionID gql.ID `json:"singleSelectOptionId"`
}

// Variable builders

// BuildCreateFieldVariables builds variables for creating a field
func BuildCreateFieldVariables(input *CreateFieldInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildUpdateFieldVariables builds variables for updating a field
func BuildUpdateFieldVariables(input *UpdateFieldInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildDeleteFieldVariables builds variables for deleting a field
func BuildDeleteFieldVariables(input *DeleteFieldInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildCreateSingleSelectFieldOptionVariables builds variables for creating an option
func BuildCreateSingleSelectFieldOptionVariables(input *CreateSingleSelectFieldOptionInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildUpdateSingleSelectFieldOptionVariables builds variables for updating an option
func BuildUpdateSingleSelectFieldOptionVariables(input *UpdateSingleSelectFieldOptionInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// BuildDeleteSingleSelectFieldOptionVariables builds variables for deleting an option
func BuildDeleteSingleSelectFieldOptionVariables(input *DeleteSingleSelectFieldOptionInput) map[string]interface{} {
	return map[string]interface{}{
		"input": *input,
	}
}

// Extended field info for display
type FieldInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        string
	Name      string
	DataType  ProjectV2FieldDataType
	Options   []FieldOptionInfo
}

type FieldOptionInfo struct {
	Description *string
	ID          string
	Name        string
	Color       string
}

// Color constants for single select options
const (
	SingleSelectColorGray   = "GRAY"
	SingleSelectColorRed    = "RED"
	SingleSelectColorOrange = "ORANGE"
	SingleSelectColorYellow = "YELLOW"
	SingleSelectColorGreen  = "GREEN"
	SingleSelectColorBlue   = "BLUE"
	SingleSelectColorPurple = "PURPLE"
	SingleSelectColorPink   = "PINK"
)

// ValidSingleSelectColors returns all valid colors for single select options
func ValidSingleSelectColors() []string {
	return []string{
		SingleSelectColorGray,
		SingleSelectColorRed,
		SingleSelectColorOrange,
		SingleSelectColorYellow,
		SingleSelectColorGreen,
		SingleSelectColorBlue,
		SingleSelectColorPurple,
		SingleSelectColorPink,
	}
}
