package response

import (
	"encoding/json"
	"net/http"
)

// Link is a HATEOAS link relation target.
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

// Links contains HATEOAS link relations.
type Links map[string]Link

// Envelope is the standard successful HTTP response shape.
type Envelope struct {
	Data  any   `json:"data,omitempty"`
	Links Links `json:"_links,omitempty"`
}

// WriteJSON writes a JSON response.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
