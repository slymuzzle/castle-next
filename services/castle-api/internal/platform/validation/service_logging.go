package validation

import (
	"reflect"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type serviceLogging struct {
	logger log.Logger
	Service
}

func NewServiceLogging(logger log.Logger, s Service) Service {
	return &serviceLogging{logger, s}
}

func (s *serviceLogging) ValidateStruct(st interface{}) (fErrs []validator.FieldError) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "ValidateStruct",
			"took", time.Since(begin),
			"struct", reflect.TypeOf(st).String(),
			"err", fErrs,
		)
	}(time.Now())
	return s.Service.ValidateStruct(st)
}

func (s *serviceLogging) ValidateGqlStruct(st interface{}) (fErrs gqlerror.List) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "ValidateGqlStruct",
			"struct", reflect.TypeOf(st).String(),
			"took", time.Since(begin),
			"err", fErrs.Error(),
		)
	}(time.Now())
	return s.Service.ValidateGqlStruct(st)
}
