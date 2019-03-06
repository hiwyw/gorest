package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/zdnscloud/gorest/types"
)

func CreateHandler(apiContext *types.APIContext) *types.APIError {
	handler := apiContext.Schema.Handler
	if handler == nil {
		return types.NewAPIError(types.NotFound, "no found schema handler")
	}

	yamlContent, object, err := parseCreateBody(apiContext)
	if err != nil {
		return err
	}

	result, err := handler.Create(object, yamlContent)
	if err != nil {
		return err
	}

	if obj, ok := result.(types.Object); ok {
		obj.SetCreationTimestamp(time.Now().String())
	}
	addResourceLinks(apiContext, result)
	WriteResponse(apiContext, http.StatusCreated, result)
	return nil
}

func DeleteHandler(apiContext *types.APIContext) *types.APIError {
	handler := apiContext.Schema.Handler
	if handler == nil {
		return types.NewAPIError(types.NotFound, "no found schema handler")
	}

	obj, err := getObject(apiContext, getSchemaStructVal(apiContext))
	if err != nil {
		return err
	}

	obj.SetID(apiContext.ID)
	if err = handler.Delete(obj); err != nil {
		return err
	}

	WriteResponse(apiContext, http.StatusOK, nil)
	return nil
}

func UpdateHandler(apiContext *types.APIContext) *types.APIError {
	handler := apiContext.Schema.Handler
	if handler == nil {
		return types.NewAPIError(types.NotFound, "no found schema handler")
	}

	val := getSchemaStructVal(apiContext)
	if err := decodeBody(apiContext.Request, val); err != nil {
		return err
	}

	object, err := getObject(apiContext, val)
	if err != nil {
		return err
	}

	object.SetID(apiContext.ID)
	result, err := handler.Update(object)
	if err != nil {
		return err
	}

	addResourceLinks(apiContext, result)
	WriteResponse(apiContext, http.StatusOK, result)
	return nil
}

func ListHandler(apiContext *types.APIContext) *types.APIError {
	handler := apiContext.Schema.Handler
	if handler == nil {
		return types.NewAPIError(types.NotFound, "no found schema handler")
	}

	var result interface{}
	obj, err := getObject(apiContext, getSchemaStructVal(apiContext))
	if err != nil {
		return err
	}

	if apiContext.ID == "" {
		collection := &types.Collection{
			Type:         "collection",
			ResourceType: apiContext.Schema.ID,
			Data:         handler.List(obj),
		}
		addCollectionLinks(apiContext, collection)
		result = collection
	} else {
		obj.SetID(apiContext.ID)
		result = handler.Get(obj)
		addResourceLinks(apiContext, result)
	}

	WriteResponse(apiContext, http.StatusOK, result)
	return nil
}

func ActionHandler(apiContext *types.APIContext, action *types.Action) *types.APIError {
	handler := apiContext.Schema.Handler
	if handler == nil {
		return types.NewAPIError(types.NotFound, "no found schema handler")
	}

	var params map[string]interface{}
	if err := decodeBody(apiContext.Request, &params); err != nil {
		return err
	}

	obj, err := getObject(apiContext, getSchemaStructVal(apiContext))
	if err != nil {
		return err
	}

	obj.SetID(apiContext.ID)
	result, err := handler.Action(obj, apiContext.Action, params)
	if err != nil {
		return err
	}

	WriteResponse(apiContext, http.StatusOK, result)
	return nil
}

func getSchemaStructVal(apiContext *types.APIContext) interface{} {
	val := apiContext.Schema.StructVal
	valPtr := reflect.New(val.Type())
	valPtr.Elem().Set(val)
	return valPtr.Interface()
}

func decodeBody(req *http.Request, params interface{}) *types.APIError {
	reqBody, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return types.NewAPIError(types.InvalidBodyContent,
			fmt.Sprintf("Failed to read request body: %v", err.Error()))
	}

	err = json.Unmarshal(reqBody, params)
	if err != nil {
		return types.NewAPIError(types.InvalidBodyContent,
			fmt.Sprintf("Failed to parse request body: %v", err.Error()))
	}

	return nil
}

func parseCreateBody(apiContext *types.APIContext) ([]byte, types.Object, *types.APIError) {
	var (
		params struct {
			Yaml string `json:"yaml_"`
		}
		content []byte
	)

	reqBody, err := ioutil.ReadAll(apiContext.Request.Body)
	defer apiContext.Request.Body.Close()
	if err != nil {
		return nil, nil, types.NewAPIError(types.InvalidBodyContent,
			fmt.Sprintf("Failed to read request body: %v", err.Error()))
	}

	if err := json.Unmarshal(reqBody, &params); err != nil {
		return nil, nil, types.NewAPIError(types.InvalidBodyContent,
			fmt.Sprintf("Failed to parse request body: %v", err.Error()))
	}

	if params.Yaml != "" {
		content, err = base64.StdEncoding.DecodeString(params.Yaml)
		if err != nil {
			return nil, nil, types.NewAPIError(types.InvalidBodyContent,
				fmt.Sprintf("Failed to parse request body: %v", err.Error()))
		}
	}

	val := getSchemaStructVal(apiContext)
	if err := json.Unmarshal(reqBody, val); err != nil {
		return nil, nil, types.NewAPIError(types.InvalidBodyContent,
			fmt.Sprintf("Failed to parse request body: %v", err.Error()))
	}

	obj, apiErr := getObject(apiContext, val)
	return content, obj, apiErr
}

func getObject(apiContext *types.APIContext, val interface{}) (types.Object, *types.APIError) {
	if obj, ok := val.(types.Object); ok {
		obj.SetType(apiContext.Schema.ID)
		obj.SetParent(apiContext.Parent)
		return obj, nil
	} else {
		return nil, types.NewAPIError(types.NotFound, fmt.Sprintf("no found resource schema"))
	}
}
