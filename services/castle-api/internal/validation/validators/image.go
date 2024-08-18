package validators

import (
	"slices"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
)

type Upload = graphql.Upload

var ImageContentTypes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
	"image/tiff",
}

func IsValidImage(contentType string) bool {
	return slices.Contains(ImageContentTypes, contentType)
}

func GqlUploadImageValidator(sl validator.StructLevel) {
	upload := sl.Current().Interface().(Upload)

	if !IsValidImage(upload.ContentType) {
		sl.ReportError(upload.ContentType, "ContentType", "", "", "")
	}
}
