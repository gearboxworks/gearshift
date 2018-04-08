package main

import (
	"encoding/json"
)

type JsonResponse struct {
	Status string	`json:"status"`
	Message string  `json:"message"`
	Error error
}

// NewJsonResponse() is a constructor for JsonResponse
//
func NewJsonResponse( s, m string ) *JsonResponse {
	return &JsonResponse{
		Status:	s,
		Message: m,
	}
}

func (jr *JsonResponse) toString() (js string) {
	var err error
	ba, err := json.Marshal(jr)
	if err != nil {
		return ""
	}
	return string(ba)
}
