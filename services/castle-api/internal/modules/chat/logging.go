package chat

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"

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

func (s *loggingService) SendMessage(
	ctx context.Context,
	input SendMessageInput,
) (msg *ent.MessageEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SendMessage",
			"targetUserID", input.TargetUserID,
			"replyTo", input.ReplyTo,
			"content", input.Content,
			"filesCount", len(input.Files),
			"message", msg.Node,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendMessage(ctx, input)
}

func (s *loggingService) SubscribeToMessageAddedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (ch <-chan *ent.MessageEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToMessageAddedEvent",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToMessageAddedEvent(ctx, roomID)
}

func (s *loggingService) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	input UpdateMessageInput,
) (msg *ent.MessageEdge, err error) {
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

func (s *loggingService) SubscribeToMessageUpdatedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (ch <-chan *ent.MessageEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToMessageUpdatedEvent",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToMessageUpdatedEvent(ctx, roomID)
}

func (s *loggingService) DeleteMessage(
	ctx context.Context,
	messageID pulid.ID,
) (msg *ent.MessageEdge, err error) {
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

func (s *loggingService) SubscribeToMessageDeletedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (ch <-chan *ent.MessageEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToMessageDeletedEvent",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToMessageDeletedEvent(ctx, roomID)
}
