package chat

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"journeyhub/graph/model"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type serviceLogging struct {
	logger log.Logger
	Service
}

func NewServiceLogging(logger log.Logger, s Service) Service {
	return &serviceLogging{logger, s}
}

func (s *serviceLogging) SendMessage(
	ctx context.Context,
	input model.SendMessageInput,
) (msg *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SendMessage",
			// "targetUserID", input.TargetUserID,
			// "replyTo", input.ReplyTo,
			// "content", input.Content,
			// "filesCount", len(input.Files),
			// "message", msg.Node,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendMessage(ctx, input)
}

func (s *serviceLogging) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	input model.UpdateMessageInput,
) (msg *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "UpdateMessage",
			"messageID", messageID,
			"content", input.Content,
			"message", msg,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UpdateMessage(ctx, messageID, input)
}

func (s *serviceLogging) DeleteMessage(
	ctx context.Context,
	messageID pulid.ID,
) (msg *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "DeleteMessage",
			"messageID", messageID,
			"message", msg,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteMessage(ctx, messageID)
}
