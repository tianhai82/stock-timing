package openapi

type TigerResponse struct {
	Code    int
	Message string
	Data    string
}

func (r *TigerResponse) IsSuccess() bool {
	return r.Code == 0
}

func (r *TigerResponse) ParseResponseContent(response map[string]interface{}) map[string]interface{} {
	if code, ok := response["code"]; ok {
		r.Code = int(code.(float64))
	}
	if message, ok := response["message"]; ok {
		r.Message = message.(string)
	}
	if data, ok := response["data"]; ok {
		r.Data = data.(string)
	}
	return response
}
