package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// DecodeJSON decodes a request JSON body into dst.
func DecodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("request body is required")
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("decode json body: %w", err)
	}
	return nil
}
