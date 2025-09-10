package utils

import (
	"strings"

	id "github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
	indonesian := id.New()
	uni := ut.New(indonesian, indonesian)
	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validate, trans)

	return &Validator{
		validate: validate,
		trans:    trans,
	}
}

// Validasi dan kembalikan map error tertranslate
func (v *Validator) ValidateStruct(s interface{}) (errorMap map[string]string, isValid bool) {
	err := v.validate.Struct(s)
	if err == nil {
		return nil, err == nil
	}

	errs := err.(validator.ValidationErrors)
	out := make(map[string]string)
	for _, e := range errs {
		// Field ke lowercase jika mau
		field := strings.ToLower(e.Field())
		out[field] = e.Translate(v.trans)
	}
	return out, out == nil
}

func (v *Validator) ValidateVar(s interface{}, tag string) map[string]string {
	err := v.validate.Var(s, tag)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	out := make(map[string]string)
	for _, e := range errs {
		// Field ke lowercase jika mau
		out[e.Field()] = e.Translate(v.trans)
	}
	return out
}
