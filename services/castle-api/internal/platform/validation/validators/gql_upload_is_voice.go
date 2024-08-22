package validators

import (
	"slices"

	"github.com/99designs/gqlgen/graphql"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var voiceContentTypes = []string{
	"audio/ogg",
}

func handleGqlUploadIsVoiceTag(fl validator.FieldLevel) bool {
	upload := fl.Field().Interface().(*graphql.Upload)
	if upload == nil {
		return true
	}
	return slices.Contains(voiceContentTypes, upload.ContentType)
}

func RegisterGqlUploadIsVoiceTag(val *validator.Validate, trans ut.Translator) error {
	if err := val.RegisterValidation("gql_upload_is_voice", handleGqlUploadIsVoiceTag, true); err != nil {
		return err
	}

	if err := val.RegisterTranslation(
		"gql_upload_is_voice",
		trans,
		registerTranslator("gql_upload_is_voice", "{0} must be in (audio/ogg) format"),
		translate,
	); err != nil {
		return err
	}

	return nil
}
