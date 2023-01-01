package openapi

import (
	"encoding/json"

	"github.com/tianhai82/stock-timing/tiger/model"
)

type OpenApiRequest struct {
	Method    string
	Timestamp string
	BizModel  interface{}
}

func (r *OpenApiRequest) GetParams() map[string]string {
	params := make(map[string]string)
	params[model.P_METHOD] = r.Method
	params[model.P_TIMESTAMP] = r.Timestamp
	if r.BizModel != nil {
		bizMap, ok := r.BizModel.(map[string]interface{})
		if ok {
			version, ok := bizMap[model.P_VERSION]
			if ok {
				params[model.P_VERSION] = version.(string)
			}
		}
		bizContent, err := json.Marshal(r.BizModel)
		if err != nil {
			panic("fail to marshal BizModel: " + err.Error())
		}
		params[model.P_BIZ_CONTENT] = string(bizContent)
	}
	return params
}
