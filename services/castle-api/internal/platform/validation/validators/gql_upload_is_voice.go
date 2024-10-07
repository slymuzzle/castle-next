package validators

import (
	"slices"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gabriel-vasile/mimetype"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/k0kubun/pp/v3"
)

var voiceContentTypes = []string{
	"audio/mp4",
	"video/mp4",
}

func handleGqlUploadIsVoiceTag(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}

	upload := fl.Field().Interface().(graphql.Upload)

	mtype, err := mimetype.DetectReader(upload.File)
	if err != nil {
		return false
	}

	pp.Print(upload)

	pp.Print(mtype)

	return slices.Contains(voiceContentTypes, mtype.String())
}

func RegisterGqlUploadIsVoiceTag(val *validator.Validate, trans ut.Translator) error {
	if err := val.RegisterValidation("gql_upload_is_voice", handleGqlUploadIsVoiceTag, true); err != nil {
		return err
	}

	if err := val.RegisterTranslation(
		"gql_upload_is_voice",
		trans,
		registerTranslator("gql_upload_is_voice", "{0} must be in (audio/mp4) format"),
		translate,
	); err != nil {
		return err
	}

	return nil
}
