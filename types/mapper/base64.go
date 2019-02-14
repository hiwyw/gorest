package mapper

import (
	"encoding/base64"
	"strings"

	"github.com/zdnscloud/gorest/types"
	"github.com/zdnscloud/gorest/types/convert"
	"github.com/zdnscloud/gorest/types/values"
)

type Base64 struct {
	Field            string
	IgnoreDefinition bool
	Separator        string
}

func (m Base64) FromInternal(data map[string]interface{}) {
	if v, ok := values.RemoveValue(data, strings.Split(m.Field, m.getSep())...); ok {
		str := convert.ToString(v)
		if str == "" {
			return
		}

		newData, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			//TODO add log
		}

		values.PutValue(data, string(newData), strings.Split(m.Field, m.getSep())...)
	}
}

func (m Base64) ToInternal(data map[string]interface{}) error {
	if v, ok := values.RemoveValue(data, strings.Split(m.Field, m.getSep())...); ok {
		str := convert.ToString(v)
		if str == "" {
			return nil
		}

		newData := base64.StdEncoding.EncodeToString([]byte(str))
		values.PutValue(data, newData, strings.Split(m.Field, m.getSep())...)
	}

	return nil
}

func (m Base64) ModifySchema(s *types.Schema, schemas *types.Schemas) error {
	if !m.IgnoreDefinition {
		if err := ValidateField(m.Field, s); err != nil {
			return err
		}
	}

	return nil
}

func (m Base64) getSep() string {
	if m.Separator == "" {
		return "/"
	}
	return m.Separator
}
