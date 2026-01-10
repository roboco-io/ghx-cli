package graphql

import (
	"testing"

	gql "github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

func TestFieldMutations(t *testing.T) {
	t.Run("CreateField mutation structure", func(t *testing.T) {
		mutation := &CreateFieldMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("UpdateField mutation structure", func(t *testing.T) {
		mutation := &UpdateFieldMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("DeleteField mutation structure", func(t *testing.T) {
		mutation := &DeleteFieldMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("CreateSingleSelectFieldOption mutation structure", func(t *testing.T) {
		mutation := &CreateSingleSelectFieldOptionMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("UpdateSingleSelectFieldOption mutation structure", func(t *testing.T) {
		mutation := &UpdateSingleSelectFieldOptionMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("DeleteSingleSelectFieldOption mutation structure", func(t *testing.T) {
		mutation := &DeleteSingleSelectFieldOptionMutation{}
		assert.NotNil(t, mutation)
	})
}

func TestFieldVariableBuilders(t *testing.T) {
	t.Run("BuildCreateFieldVariables creates proper variables", func(t *testing.T) {
		input := &CreateFieldInput{
			ProjectID: gql.ID("project-id"),
			Name:      gql.String("Priority"),
			DataType:  ProjectV2FieldDataTypeText,
		}

		variables := BuildCreateFieldVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildCreateFieldVariables with single select options", func(t *testing.T) {
		input := &CreateFieldInput{
			ProjectID: gql.ID("project-id"),
			Name:      gql.String("Priority"),
			DataType:  ProjectV2FieldDataTypeSingleSelect,
			SingleSelectOptions: []SingleSelectOption{
				{Name: gql.String("High"), Color: gql.String("GRAY")},
				{Name: gql.String("Medium"), Color: gql.String("GRAY")},
				{Name: gql.String("Low"), Color: gql.String("GRAY")},
			},
		}

		variables := BuildCreateFieldVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildUpdateFieldVariables creates proper variables", func(t *testing.T) {
		newName := gql.String("Updated Priority")
		input := &UpdateFieldInput{
			FieldID: gql.ID("field-id"),
			Name:    &newName,
		}

		variables := BuildUpdateFieldVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildDeleteFieldVariables creates proper variables", func(t *testing.T) {
		input := &DeleteFieldInput{
			FieldID: gql.ID("field-id"),
		}

		variables := BuildDeleteFieldVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildCreateSingleSelectFieldOptionVariables creates proper variables", func(t *testing.T) {
		desc := gql.String("Critical priority items")
		input := &CreateSingleSelectFieldOptionInput{
			FieldID:     gql.ID("field-id"),
			Name:        gql.String("Critical"),
			Color:       gql.String(SingleSelectColorRed),
			Description: &desc,
		}

		variables := BuildCreateSingleSelectFieldOptionVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildCreateSingleSelectFieldOptionVariables without description", func(t *testing.T) {
		input := &CreateSingleSelectFieldOptionInput{
			FieldID: gql.ID("field-id"),
			Name:    gql.String("Critical"),
			Color:   gql.String(SingleSelectColorRed),
		}

		variables := BuildCreateSingleSelectFieldOptionVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildUpdateSingleSelectFieldOptionVariables creates proper variables", func(t *testing.T) {
		newName := gql.String("Very High")
		newColor := gql.String(SingleSelectColorOrange)
		newDescription := gql.String("Very high priority")

		input := &UpdateSingleSelectFieldOptionInput{
			OptionID:    gql.ID("option-id"),
			Name:        &newName,
			Color:       &newColor,
			Description: &newDescription,
		}

		variables := BuildUpdateSingleSelectFieldOptionVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildDeleteSingleSelectFieldOptionVariables creates proper variables", func(t *testing.T) {
		input := &DeleteSingleSelectFieldOptionInput{
			OptionID: gql.ID("option-id"),
		}

		variables := BuildDeleteSingleSelectFieldOptionVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})
}

func TestFieldDataTypes(t *testing.T) {
	t.Run("All field data types defined", func(t *testing.T) {
		assert.Equal(t, "TEXT", string(ProjectV2FieldDataTypeText))
		assert.Equal(t, "NUMBER", string(ProjectV2FieldDataTypeNumber))
		assert.Equal(t, "DATE", string(ProjectV2FieldDataTypeDate))
		assert.Equal(t, "SINGLE_SELECT", string(ProjectV2FieldDataTypeSingleSelect))
		assert.Equal(t, "ITERATION", string(ProjectV2FieldDataTypeIteration))
	})
}

func TestValidSingleSelectColors(t *testing.T) {
	t.Run("All valid colors returned", func(t *testing.T) {
		colors := ValidSingleSelectColors()

		assert.Len(t, colors, 8)
		assert.Contains(t, colors, SingleSelectColorGray)
		assert.Contains(t, colors, SingleSelectColorRed)
		assert.Contains(t, colors, SingleSelectColorOrange)
		assert.Contains(t, colors, SingleSelectColorYellow)
		assert.Contains(t, colors, SingleSelectColorGreen)
		assert.Contains(t, colors, SingleSelectColorBlue)
		assert.Contains(t, colors, SingleSelectColorPurple)
		assert.Contains(t, colors, SingleSelectColorPink)
	})
}

func TestFieldInfo(t *testing.T) {
	t.Run("FieldInfo structure", func(t *testing.T) {
		info := FieldInfo{
			ID:       "field-id",
			Name:     "Priority",
			DataType: ProjectV2FieldDataTypeSingleSelect,
		}

		assert.Equal(t, "field-id", info.ID)
		assert.Equal(t, "Priority", info.Name)
		assert.Equal(t, ProjectV2FieldDataTypeSingleSelect, info.DataType)
	})

	t.Run("FieldOptionInfo structure", func(t *testing.T) {
		description := "High priority option"
		option := FieldOptionInfo{
			ID:          "option-id",
			Name:        "High",
			Color:       SingleSelectColorRed,
			Description: &description,
		}

		assert.Equal(t, "option-id", option.ID)
		assert.Equal(t, "High", option.Name)
		assert.Equal(t, SingleSelectColorRed, option.Color)
		assert.Equal(t, "High priority option", *option.Description)
	})
}
