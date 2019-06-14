package gearshift

import "encoding/json"

type JsonResponse interface {
	ToByteArray() []byte
}

func jsonResponseToByteArray(jr JsonResponse) []byte {
	var err error
	ba, err := json.Marshal(jr)
	if err != nil {
		return []byte("")
	}
	return ba
}

type OkJsonResponse struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body"`
}

func (jr OkJsonResponse) ToByteArray() []byte {
	return jsonResponseToByteArray(jr)
}

// NewJsonResponse() is a constructor for JsonResponse
func NewOkJsonResponse(b string) JsonResponse {
	js := new(interface{})
	err := json.Unmarshal([]byte(b), &js)
	if err != nil {
		panic(err)
	}
	return &OkJsonResponse{
		Status: "ok",
		Body:   js,
	}
}

type FailJsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (jr FailJsonResponse) ToByteArray() []byte {
	return jsonResponseToByteArray(jr)
}

// NewJsonResponse() is a constructor for JsonResponse
func NewFailJsonResponse(m string) JsonResponse {
	return &FailJsonResponse{
		Status:  "fail",
		Message: m,
	}
}
