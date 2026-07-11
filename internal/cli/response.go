// Package cli defines the structured JSON response contract that every easygzh
// command emits. Agents (and scripts) consume this to make decisions.
package cli

// Response is the standard JSON envelope. Every command returns one.
type Response struct {
	Success     bool        `json:"success"`
	Code        string      `json:"code"`
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	NextActions []string    `json:"next_actions,omitempty"`
	Error       string      `json:"error,omitempty"`
}

// OK returns a success Response with the given code and data.
func OK(code string, data interface{}, next ...string) Response {
	return Response{Success: true, Code: code, Data: data, NextActions: next}
}

// Fail returns a failure Response.
func Fail(code, msg string) Response {
	return Response{Success: false, Code: code, Error: msg}
}
