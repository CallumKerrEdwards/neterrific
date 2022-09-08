package neterrific

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidRequest = errors.New("invalid request body")
	vdr               *validator.Validate
)

type ValidationError struct {
	Errors []string
}

func SendValidationError(w http.ResponseWriter, err error, status int) {
	s := status
	if s == 0 {
		s = http.StatusUnprocessableEntity
	}

	if err, ok := err.(ValidationError); ok {
		SendJSON(w, s, Payload{
			"status": s,
			"errors": err.Errors,
		})
	}

	SendHTTPJSONError(w, s, err)
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

func ParseAndValidate(r *http.Request, v any) error {
	if err := ParseBody(r, v); err != nil {
		return err
	}

	if err := Validator().Struct(v); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var ve ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			ve.Errors = append(ve.Errors, err.Error())
		}

		return ve
	}

	return nil
}

func ParseBody(r *http.Request, v any) error {
	b, err := io.ReadAll(r.Body)
	if err != nil || len(b) == 0 {
		return ErrInvalidRequest
	}

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

func Validator() *validator.Validate {
	if vdr == nil {
		vdr = validator.New()
	}

	return vdr
}
