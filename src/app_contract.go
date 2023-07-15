package main

type response struct {
	Error *string     `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
