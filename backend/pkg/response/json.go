package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data any) error {
	return JSONWithHeaders(w, status, data, nil)
}

func JSONWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	res := map[string]any{
		"error":  nil,
		"data":   data,
		"status": status,
	}
	js, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		log.Printf("Error writing json response %v", err)
		return err
	}
	return nil
}

func ErrJSONWithHeaders(w http.ResponseWriter, status int, errMsg any, headers http.Header) error {
	response := map[string]any{
		"error":  errMsg,
		"data":   nil,
		"status": status,
	}
	js, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		log.Printf("Error writing json error response: %v", err)
		return err
	}
	return nil
}
