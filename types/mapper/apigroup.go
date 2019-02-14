package mapper

import (
	"github.com/zdnscloud/gorest/types"
	"github.com/zdnscloud/gorest/types/convert"
)

type APIGroup struct {
	apiVersion string
	kind       string
}

func (a *APIGroup) FromInternal(data map[string]interface{}) {
}

func (a *APIGroup) ToInternal(data map[string]interface{}) error {
	_, ok := data["apiVersion"]
	if !ok && data != nil {
		data["apiVersion"] = a.apiVersion
	}

	_, ok = data["kind"]
	if !ok && data != nil {
		data["kind"] = a.kind
	}

	return nil
}

func (a *APIGroup) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	a.apiVersion = schema.Version.Group + "/" + schema.Version.Version
	a.kind = convert.Capitalize(schema.ID)

	return nil
}
