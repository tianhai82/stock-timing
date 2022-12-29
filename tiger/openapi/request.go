package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/tianhai82/stock-timing/tiger/model"
)

type OpenApiRequest struct {
	Method   string
	BizModel interface{}
}

func (r *OpenApiRequest) GetParams() map[string]string {
	params := make(map[string]string)
	params[model.P_METHOD] = r.Method
	if r.BizModel != nil {
		v, ok := r.BizModel.(interface{ ToOpenApiDict() map[string]interface{} })
		if !ok {
			panic(fmt.Errorf("biz_model must implement ToOpenApiDict method"))
		}
		bizContent, err := json.Marshal(v.ToOpenApiDict())
		if err != nil {
			panic(err)
		}
		params[params.P_BIZ_CONTENT] = string(bizContent)
	}
	return params
}
