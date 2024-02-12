package config

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

// Available field types.
const (
	FieldTypeString = "string"
	FieldTypeNumber = "number"
	FieldTypeBool   = "bool"
	FieldTypeObject = "object"
	FieldTypeList   = "list"
	FieldTypeAny    = "any"
)

const (
	FieldVisibilityPublic        = "public"
	FieldVisibilitySelf          = "self"
	FieldVisibilityPrivate       = "private"
	FieldVisibilityAuthenticated = "authenticated"
)

// FieldConfig describes how user-extra data looks like.
type FieldConfig struct {
	Type        string         `json:"type" hcl:"type,label"`
	Name        string         `json:"name" hcl:"name,label"`
	Visibility  string         `json:"visibility" hcl:"visibility,optional"`
	Writeable   *bool          `json:"writeable" hcl:"writeable,optional"`
	Description string         `json:"description" hcl:"description,optional"`
	DisplayName string         `json:"display_name" hcl:"display_name,optional"`
	Properties  []*FieldConfig `json:"property" hcl:"property,block"`
	ElementType *FieldConfig   `json:"element_type" hcl:"element_type,block"`
}

func (fc FieldConfig) Validate(data *structpb.Value) error {
	if data == nil {
		return nil
	}

	switch fc.Type {
	case FieldTypeBool:
		if _, ok := data.Kind.(*structpb.Value_BoolValue); !ok {
			return fmt.Errorf("invalid type: expected %q but got %T", "bool", data.Kind)
		}

	case FieldTypeNumber:
		if _, ok := data.Kind.(*structpb.Value_NumberValue); !ok {
			return fmt.Errorf("invalid type: expected %q but got %T", "number", data.Kind)
		}

	case FieldTypeString:
		if _, ok := data.Kind.(*structpb.Value_StringValue); !ok {
			return fmt.Errorf("invalid type: expected %q but got %T", "string", data.Kind)
		}

	case FieldTypeObject:
		ov, ok := data.Kind.(*structpb.Value_StructValue)
		if !ok {
			return fmt.Errorf("invalid type: expected %q but got %T", "object", data.Kind)
		}

		for _, propertyConfig := range fc.Properties {
			value := ov.StructValue.Fields[propertyConfig.Name]
			if err := propertyConfig.Validate(value); err != nil {
				return fmt.Errorf("%s: %w", propertyConfig.Name, err)
			}
		}

		for key := range ov.StructValue.Fields {
			var ok bool

			for _, p := range fc.Properties {
				if p.Name == key {
					ok = true

					break
				}
			}

			if !ok {
				return fmt.Errorf("%s: key not allowed", key)
			}
		}

	case FieldTypeList:
		lv, ok := data.Kind.(*structpb.Value_ListValue)
		if !ok {
			return fmt.Errorf("invalid type: expected %q but got %T", "list", data.Kind)
		}

		for idx, value := range lv.ListValue.Values {
			if err := fc.ElementType.Validate(value); err != nil {
				return fmt.Errorf("[%d]: %w", idx, err)
			}
		}

	case FieldTypeAny:
		// no validation for "any" fields

	default:
		return fmt.Errorf("invalid field type configuration")
	}

	return nil
}

func (fc *FieldConfig) ApplyVisibility(current string, value *structpb.Value) *structpb.Value {
	effectiveVisibility := getEffectiveVisibility(current, fc.Visibility)

	if effectiveVisibility != current {
		return nil
	}

	switch fc.Type {
	case FieldTypeObject:
		ov, ok := value.Kind.(*structpb.Value_StructValue)
		if !ok {
			return nil
		}

		for _, propertyConfig := range fc.Properties {
			propertyValue := ov.StructValue.Fields[propertyConfig.Name]

			if propertyValue == nil {
				continue
			}

			propertyValue = propertyConfig.ApplyVisibility(effectiveVisibility, propertyValue)
			if propertyValue == nil {
				delete(ov.StructValue.Fields, propertyConfig.Name)
			}
		}
	}

	return value
}

func (fc *FieldConfig) ValidateConfig(fieldVisibility string) error {
	// add some sense defaults
	if fc.Type == "" {
		fc.Type = FieldTypeString
	}

	if fc.Visibility == "" {
		fc.Visibility = fieldVisibility
	}

	if !isValidFieldType(fc.Type) {
		return fmt.Errorf("invalid field type %q", fc.Type)
	}

	if !isValidFieldVisibility(fc.Visibility) {
		return fmt.Errorf("invalid field visibility %q", fc.Visibility)
	}

	effectiveVisibility := getEffectiveVisibility(fieldVisibility, fc.Visibility)
	if effectiveVisibility != fc.Visibility {
		return fmt.Errorf("parent object has stronger visibility %q, %q does not take effect", fieldVisibility, fc.Visibility)
	}

	switch fc.Type {
	case FieldTypeList:
		if fc.ElementType == nil {
			return fmt.Errorf("elementType: not set")
		}

		if err := fc.ElementType.ValidateConfig(effectiveVisibility); err != nil {
			return fmt.Errorf("elementType: %w", err)
		}

	case FieldTypeObject:
		if len(fc.Properties) == 0 {
			return fmt.Errorf("properties: not set")
		}

		for _, cfg := range fc.Properties {
			if cfg == nil {
				return fmt.Errorf("properties: %s: not set", cfg.Name)
			}

			if err := cfg.ValidateConfig(effectiveVisibility); err != nil {
				return fmt.Errorf("properties: %s: %w", cfg.Name, err)
			}
		}
	}

	return nil
}

func isValidFieldVisibility(v string) bool {
	switch v {
	case FieldVisibilityAuthenticated,
		FieldVisibilitySelf,
		FieldVisibilityPrivate,
		FieldVisibilityPublic:
		return true
	default:
		return false
	}
}

func getEffectiveVisibility(previous string, next string) string {
	m := map[string]int{
		FieldVisibilityPrivate:       0,
		FieldVisibilitySelf:          1,
		FieldVisibilityAuthenticated: 2,
		FieldVisibilityPublic:        3,
	}

	previousN := m[previous]
	nextN := m[next]

	if previousN < nextN {
		return previous
	}

	return next
}

func isValidFieldType(v string) bool {
	switch v {
	case FieldTypeBool,
		FieldTypeList,
		FieldTypeNumber,
		FieldTypeObject,
		FieldTypeString:
		return true
	default:
		return false
	}
}
