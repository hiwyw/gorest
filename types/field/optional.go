package field

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	requiredTag = "required="
	defaultTag  = "default="
)

func fieldParseOptional(f Field, kind reflect.Kind, restTags []string) error {
	hasRequired := false
	hasDefault := false
	for _, tag := range restTags {
		if strings.HasPrefix(tag, requiredTag) {
			requiredVal := strings.TrimPrefix(tag, requiredTag)
			hasRequired = true
			if hasDefault {
				return fmt.Errorf("require field cann't have default value")
			}

			if requiredVal == "no" || requiredVal == "false" {
				f.SetRequired(false)
			} else if requiredVal == "yes" || requiredVal == "true" {
				f.SetRequired(true)
			} else {
				return fmt.Errorf("invalid require value %s", requiredVal)
			}
		} else if strings.HasPrefix(tag, defaultTag) {
			if hasDefault {
				return fmt.Errorf("can only have one default value")
			}
			hasDefault = true
			if hasRequired {
				return fmt.Errorf("require field cann't have default value")
			}
			defaultValue := strings.TrimPrefix(tag, defaultTag)
			switch kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				def, err := strconv.ParseInt(defaultValue, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid default value %s for int", defaultValue)
				}
				f.SetDefault(def)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				def, err := strconv.ParseUint(defaultValue, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid default value %s for int", defaultValue)
				}
				f.SetDefault(def)
			case reflect.String:
				f.SetDefault(defaultValue)
			case reflect.Bool:
				if defaultValue == "true" {
					f.SetDefault(true)
				} else if defaultValue == "false" {
					f.SetDefault(false)
				} else {
					return fmt.Errorf("invalid default value %s for bool", defaultValue)
				}
			default:
				return fmt.Errorf("only primary type support default value")
			}
		}
	}
	return nil
}
