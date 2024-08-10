package utils

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Response struct {
	Status     string      `json:"status,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
	Name       string      `json:"name,omitempty"` //name of the error
	Message    string      `json:"message,omitempty"`
	Error      interface{} `json:"error,omitempty"` //for errors that occur even if request is successful
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Extra      interface{} `json:"extra,omitempty"`
}

// BuildResponse method is to inject data value to dynamic success response
func BuildSuccessResponse(code int, message string, data interface{}, pagination ...interface{}) Response {
	res := ResponseMessage(code, "success", "", message, nil, data, pagination, nil)
	return res
}

// BuildErrorResponse method is to inject data value to dynamic failed response
func BuildErrorResponse(code int, status string, message string, err interface{}, data interface{}, logger ...bool) Response {
	res := ResponseMessage(code, status, "", message, err, data, nil, nil)
	return res
}

// ResponseMessage method for the central response holder
func ResponseMessage(code int, status string, name string, message string, err interface{}, data interface{}, pagination interface{}, extra interface{}) Response {
	if pagination != nil && reflect.ValueOf(pagination).IsNil() {
		pagination = nil
	}

	if code == http.StatusInternalServerError {
		fmt.Println("internal server error", message, err, data)
		message = "internal server error"
		err = message
	}

	var errorField interface{}
    if err != nil {
        if e, ok := err.(error); ok {
            errorField = e.Error()
        } else {
            errorField = err
        }
    }
	res := Response{
		StatusCode: code,
		Name:       name,
		Status:     status,
		Message:    message,
		Error:      errorField,
		Data:       data,
		Pagination: pagination,
		Extra:      extra,
	}

	return res
}

func UnauthorisedResponse(code int, status string, name string, message string) Response {
	res := ResponseMessage(http.StatusUnauthorized, status, name, message, nil, nil, nil, nil)
	return res
}


func ValidationResponse(err error, validate *validator.Validate, obj interface{}) map[string]string {
	errs := err.(validator.ValidationErrors)
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	// Custom map to hold simplified error messages
	errorMessages := make(map[string]string)

	// Iterate over the validation errors and format them using JSON tags
	for _, e := range errs {
		// Get the JSON field name
		jsonFieldName := getJSONFieldName(obj, e.StructField())
		errorMessages[jsonFieldName] = getCustomErrorMessage(e)
	}

	return errorMessages
}
func getCustomErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required."
	case "email":
		return "Invalid email format."
	default:
		return "Invalid value."
	}
}
func getJSONFieldName(obj interface{}, fieldName string) string {
	t := reflect.TypeOf(obj)
	field, found := t.FieldByName(fieldName)
	if !found {
		return fieldName // Fallback to the struct field name if JSON tag is not found
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName // Fallback to the struct field name if JSON tag is not found
	}
	return jsonTag
}