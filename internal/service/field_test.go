package service

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
)

func TestFieldService(t *testing.T) {
	t.Run("NewFieldService creates new service", func(t *testing.T) {
		client := api.NewClient("test-token")
		service := NewFieldService(client)

		assert.NotNil(t, service)
		assert.IsType(t, &FieldService{}, service)
	})

	t.Run("CreateField with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewFieldService(client)

		ctx := context.Background()
		input := CreateFieldInput{
			ProjectID: "test-project-id",
			Name:      "Priority",
			DataType:  graphql.ProjectV2FieldDataTypeText,
		}

		field, err := service.CreateField(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, field)
		assert.Contains(t, err.Error(), "failed to create field")
	})

	t.Run("UpdateField with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewFieldService(client)

		ctx := context.Background()
		newName := "Updated Priority"
		input := UpdateFieldInput{
			FieldID: "test-field-id",
			Name:    &newName,
		}

		field, err := service.UpdateField(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, field)
		assert.Contains(t, err.Error(), "failed to update field")
	})

	t.Run("DeleteField with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewFieldService(client)

		ctx := context.Background()
		input := DeleteFieldInput{
			FieldID: "test-field-id",
		}

		err := service.DeleteField(ctx, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete field")
	})

	t.Run("CreateFieldOption with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewFieldService(client)

		ctx := context.Background()
		input := CreateFieldOptionInput{
			FieldID: "test-field-id",
			Name:    "High",
			Color:   "RED",
		}

		option, err := service.CreateFieldOption(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, option)
		assert.Contains(t, err.Error(), "failed to create field option")
	})

	t.Run("UpdateFieldOption with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewFieldService(client)

		ctx := context.Background()
		newName := "Very High"
		input := UpdateFieldOptionInput{
			OptionID: "test-option-id",
			Name:     &newName,
		}

		option, err := service.UpdateFieldOption(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, option)
		assert.Contains(t, err.Error(), "failed to update field option")
	})

	t.Run("DeleteFieldOption with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewFieldService(client)

		ctx := context.Background()
		input := DeleteFieldOptionInput{
			OptionID: "test-option-id",
		}

		err := service.DeleteFieldOption(ctx, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete field option")
	})

	t.Run("GetProjectFields with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewFieldService(client)

		ctx := context.Background()
		fields, err := service.GetProjectFields(ctx, "testowner", 1, false)

		assert.Error(t, err)
		assert.Nil(t, fields)
		assert.Contains(t, err.Error(), "failed to get project")
	})
}

func TestFieldValidation(t *testing.T) {
	t.Run("ValidateFieldName accepts valid names", func(t *testing.T) {
		err := ValidateFieldName("Priority")
		assert.NoError(t, err)

		err = ValidateFieldName("Status Category")
		assert.NoError(t, err)
	})

	t.Run("ValidateFieldName rejects empty names", func(t *testing.T) {
		err := ValidateFieldName("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field name cannot be empty")

		err = ValidateFieldName("   ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field name cannot be empty")
	})

	t.Run("ValidateFieldName rejects long names", func(t *testing.T) {
		longName := strings.Repeat("a", 101)
		err := ValidateFieldName(longName)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot exceed 100 characters")
	})

	t.Run("ValidateFieldType accepts valid types", func(t *testing.T) {
		dataType, err := ValidateFieldType("text")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2FieldDataTypeText, dataType)

		dataType, err = ValidateFieldType("NUMBER")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2FieldDataTypeNumber, dataType)

		dataType, err = ValidateFieldType("Date")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2FieldDataTypeDate, dataType)

		dataType, err = ValidateFieldType("single_select")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2FieldDataTypeSingleSelect, dataType)

		dataType, err = ValidateFieldType("ITERATION")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2FieldDataTypeIteration, dataType)
	})

	t.Run("ValidateFieldType rejects invalid types", func(t *testing.T) {
		_, err := ValidateFieldType("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid field type: invalid")
	})

	t.Run("ValidateColor accepts valid colors", func(t *testing.T) {
		err := ValidateColor("red")
		assert.NoError(t, err)

		err = ValidateColor("BLUE")
		assert.NoError(t, err)

		err = ValidateColor("Green")
		assert.NoError(t, err)
	})

	t.Run("ValidateColor rejects invalid colors", func(t *testing.T) {
		err := ValidateColor("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid color: invalid")
	})

	t.Run("NormalizeColor converts to uppercase", func(t *testing.T) {
		assert.Equal(t, "RED", NormalizeColor("red"))
		assert.Equal(t, "BLUE", NormalizeColor("Blue"))
		assert.Equal(t, "GREEN", NormalizeColor("GREEN"))
	})
}

func TestFieldFormatting(t *testing.T) {
	t.Run("FormatFieldDataType formats correctly", func(t *testing.T) {
		assert.Equal(t, "Text", FormatFieldDataType(graphql.ProjectV2FieldDataTypeText))
		assert.Equal(t, "Number", FormatFieldDataType(graphql.ProjectV2FieldDataTypeNumber))
		assert.Equal(t, "Date", FormatFieldDataType(graphql.ProjectV2FieldDataTypeDate))
		assert.Equal(t, "Single Select", FormatFieldDataType(graphql.ProjectV2FieldDataTypeSingleSelect))
		assert.Equal(t, "Iteration", FormatFieldDataType(graphql.ProjectV2FieldDataTypeIteration))
	})

	t.Run("FormatColor formats correctly", func(t *testing.T) {
		assert.Equal(t, "Red", FormatColor("RED"))
		assert.Equal(t, "Blue", FormatColor("BLUE"))
		assert.Equal(t, "Green", FormatColor("green"))
		assert.Equal(t, "", FormatColor(""))
	})
}

func TestFieldInfo(t *testing.T) {
	t.Run("FieldInfo structure", func(t *testing.T) {
		info := FieldInfo{
			ID:          "field-id",
			Name:        "Priority",
			DataType:    graphql.ProjectV2FieldDataTypeSingleSelect,
			ProjectID:   "project-id",
			ProjectName: "Test Project",
		}

		assert.Equal(t, "field-id", info.ID)
		assert.Equal(t, "Priority", info.Name)
		assert.Equal(t, graphql.ProjectV2FieldDataTypeSingleSelect, info.DataType)
		assert.Equal(t, "project-id", info.ProjectID)
		assert.Equal(t, "Test Project", info.ProjectName)
	})

	t.Run("FieldOptionInfo structure", func(t *testing.T) {
		description := "High priority option"
		option := FieldOptionInfo{
			ID:          "option-id",
			Name:        "High",
			Color:       "RED",
			Description: &description,
		}

		assert.Equal(t, "option-id", option.ID)
		assert.Equal(t, "High", option.Name)
		assert.Equal(t, "RED", option.Color)
		assert.Equal(t, "High priority option", *option.Description)
	})
}
