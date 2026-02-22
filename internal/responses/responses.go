package responses

import (
	"encoding/json"
	"net/http"
)

func SetJsonBody(w http.ResponseWriter, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)

	if err != nil {
		return err
	}

	return nil
}
