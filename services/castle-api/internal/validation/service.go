package validation

import (
	"journeyhub/internal/validation/validators"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Service interface {
	ValidateStruct(interface{}) []validator.FieldError
	ValidateGqlStruct(interface{}) gqlerror.List
}

type service struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewService() Service {
	validate := validator.New()

	validate.RegisterStructValidation(
		validators.GqlUploadImageValidator,
		validators.Upload{},
	)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	return &service{
		validate: validate,
		trans:    trans,
	}
}

func (s *service) ValidateStruct(st interface{}) []validator.FieldError {
	if err := s.validate.Struct(st); err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func (s *service) ValidateGqlStruct(st interface{}) gqlerror.List {
	if err := s.validate.Struct(st); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errList := gqlerror.List{}
		for _, err := range validationErrors {
			errList = append(errList, gqlerror.Errorf(
				err.Translate(s.trans),
			))
		}

		return errList
	}

	return nil
}
