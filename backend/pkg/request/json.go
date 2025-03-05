package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	return decodeJSON(w, r, target, false)
}

func DecodeJSONStrict(w http.ResponseWriter, r *http.Request, target any) error {
	return decodeJSON(w, r, target, true)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any, disallowUnknownFields bool) error {
	var (
		maxBytes            = 1_048_576
		syntaxErr           *json.SyntaxError
		unmarshalTypeErr    *json.UnmarshalTypeError
		invalidUnmarshalErr *json.InvalidUnmarshalError
	)

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	if disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(target)

	if err != nil {
		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("malformed json body at char %d", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed json body")
		case errors.As(err, &unmarshalTypeErr):
			if unmarshalTypeErr.Field != "" {
				return fmt.Errorf("incorrect json key type at: %q", unmarshalTypeErr.Field)
			}
			return fmt.Errorf("json contains incorrect value type at char number %d", unmarshalTypeErr.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("json body cannot be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("unknown json key %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body exceeds maximum allowed byte size %d", maxBytes)
		case errors.As(err, &invalidUnmarshalErr):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body contains multiple jsons, only single json allowed")
	}
	return nil
}
