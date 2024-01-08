package config

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

type FieldType string

const (
	FieldTypeString = FieldType("string")
	FieldTypeNumber = FieldType("number")
	FieldTypeBool   = FieldType("bool")
	FieldTypeObject = FieldType("object")
	FieldTypeList   = FieldType("list")
)

type FieldVisibility string

const (
	FieldVisibilityPublic        = FieldVisibility("public")
	FieldVisibilitySelf          = FieldVisibility("self")
	FieldVisibilityPrivate       = FieldVisibility("private")
	FieldVisibilityAuthenticated = FieldVisibility("authenticated")
)

// FieldConfig describes how user-extra data looks like.
type FieldConfig struct {
	Name        string          `json:"name" hcl:",label"`
	Type        FieldType       `json:"type" hcl:"type,optional"`
	Visibility  FieldVisibility `json:"visibility" hcl:"visibility,optional"`
	Writeable   bool            `json:"writeable" hcl:"writeable,optional"`
	Description string          `json:"description" hcl:"description,optional"`
	DisplayName string          `json:"display_name" hcl:"display_name,optional"`
	Properties  []*FieldConfig  `json:"properties" hcl:"properties,block"`
	ElementType *FieldConfig    `json:"element_type" hcl:"element_type,block"`
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

	default:
		return fmt.Errorf("invalid field type configuration")
	}

	return nil
}

func (fc *FieldConfig) ApplyVisibility(current FieldVisibility, value *structpb.Value) *structpb.Value {
	effectiveVisilbity := getEffectiveVisibility(current, fc.Visibility)

	if effectiveVisilbity != current {
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

			propertyValue = propertyConfig.ApplyVisibility(effectiveVisilbity, propertyValue)
			if propertyValue == nil {
				delete(ov.StructValue.Fields, propertyConfig.Name)
			}
		}
	}

	return value
}

func (fc *FieldConfig) ValidateConfig(fieldVisiblity FieldVisibility) error {
	// add some sense defaults
	if fc.Type == "" {
		fc.Type = FieldTypeString
	}

	if fc.Visibility == "" {
		fc.Visibility = fieldVisiblity
	}

	if !isValidFieldType(fc.Type) {
		return fmt.Errorf("invalid field type %q", fc.Type)
	}

	if !isValidFieldVisiblity(fc.Visibility) {
		return fmt.Errorf("invalid field visibility %q", fc.Visibility)
	}

	effectiveVisibility := getEffectiveVisibility(fieldVisiblity, fc.Visibility)
	if effectiveVisibility != fc.Visibility {
		return fmt.Errorf("parent object has stronger visibility %q, %q does not take effect", fieldVisiblity, fc.Visibility)
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

		for key, cfg := range fc.Properties {
			if cfg == nil {
				return fmt.Errorf("properties: %s: not set", key)
			}

			if err := cfg.ValidateConfig(effectiveVisibility); err != nil {
				return fmt.Errorf("properties: %s: %w", key, err)
			}
		}
	}

	return nil
}

func isValidFieldVisiblity(v FieldVisibility) bool {
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

func getEffectiveVisibility(previous FieldVisibility, next FieldVisibility) FieldVisibility {
	m := map[FieldVisibility]int{
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

func isValidFieldType(v FieldType) bool {
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
