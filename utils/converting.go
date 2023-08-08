package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetDataFromBody(w http.ResponseWriter, r *http.Request, str interface{}) error {
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()

	err := json.Unmarshal(body, str)
	if err != nil {
		NewJsonError(w, http.StatusInternalServerError, err.Error())
		return fmt.Errorf("cant unpack payload %w", err)
	}
	return nil
}

func PrepareDataForSend(w http.ResponseWriter, r *http.Request, str interface{}) error {
	returnableData, err := json.Marshal(str)
	if err != nil {
		NewJsonError(w, http.StatusInternalServerError, err.Error())
		return fmt.Errorf("cant pack payload %w", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(returnableData)
	return nil
}
