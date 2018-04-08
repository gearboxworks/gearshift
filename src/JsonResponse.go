package main

import "encoding/json"

type JsonResponse interface {
	ToByteArray() []byte
}

type OkJsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}
type FailJsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
}

// NewJsonResponse() is a constructor for JsonResponse
func NewJsonResponse(s, m, b string) JsonResponse {
	var jr JsonResponse
	js := new(interface{})
	json.Unmarshal([]byte(b),&js)
	if s == "ok" {
		jr = &OkJsonResponse{
			Status:  s,
			Message: m,
			Body:    js,
		}
	} else {
		jr = &FailJsonResponse{
			Status:  s,
			Message: m,
		}
	}
	return jr
}

func (jr OkJsonResponse) ToByteArray() []byte {
	return jsonResponseToByteArray(jr)
}
func (jr FailJsonResponse) ToByteArray() []byte {
	return jsonResponseToByteArray(jr)
}
func jsonResponseToByteArray(jr JsonResponse) []byte {
	var err error
	ba, err := json.Marshal(jr)
	if err != nil {
		return []byte("")
	}
	return ba
}

