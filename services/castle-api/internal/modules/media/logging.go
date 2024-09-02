package media

import (
	"context"
	"time"

	"journeyhub/graph/model"
	"journeyhub/internal/platform/config"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) UploadMessageFiles(
	ctx context.Context,
	prefix string,
	files []*model.UploadMessageFile,
) ([]*UploadInfo, error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "UploadMessageFiles",
			"host", s.Service.Config().Host,
			"ssl", s.Service.Config().Ssl,
			"count", len(files),
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.UploadMessageFiles(ctx, prefix, files)
}

func (s *loggingService) UploadFile(
	ctx context.Context,
	prefix string,
	file *graphql.Upload,
) (*UploadInfo, error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "UploadFile",
			"host", s.Service.Config().Host,
			"ssl", s.Service.Config().Ssl,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.UploadFile(ctx, prefix, file)
}

func (s *loggingService) Config() (
	config config.S3Config,
) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "Config",
			"host", s.Service.Config().Host,
			"ssl", s.Service.Config().Ssl,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Config()
}
