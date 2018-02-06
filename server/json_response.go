package main

type jsonResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	ResponseCode int    `json:"response"`
	Data         interface{}
}
