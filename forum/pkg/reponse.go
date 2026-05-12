package pkg

import (
	"encoding/json"
	"net/http"
)

func DecodeJson(req *http.Request, data any) error {
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		return err
	}
	return nil
}

func EncodeJson(reader http.ResponseWriter, data any) error {
	err := json.NewEncoder(reader).Encode(&data)

	if err != nil {
		return err
	}
	return nil
}
