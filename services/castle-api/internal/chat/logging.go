package chat

import (
	"context"
	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"
	"time"

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
	currentUserID pulid.ID,
	targetUserID pulid.ID,
	replyTo pulid.ID,
	content string,
) (msg *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SendMessage",
			"currentUserID", currentUserID,
			"targetUserID", targetUserID,
			"replyTo", replyTo,
			"content", content,
			"message", msg,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendMessage(ctx, currentUserID, targetUserID, replyTo, content)
}

func (s *loggingService) SubscribeToMessageAddedEvent(roomID pulid.ID) (ch <-chan *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToMessageAddedEvent",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToMessageAddedEvent(roomID)
}

func (s *loggingService) UpdateMessage(
	ctx context.Context,
	messageID pulid.ID,
	content string,
) (msg *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "UpdateMessage",
			"messageID", messageID,
			"content", content,
			"message", msg,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UpdateMessage(ctx, messageID, content)
}

func (s *loggingService) SubscribeToMessageUpdatedEvent(roomID pulid.ID) (ch <-chan *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToMessageUpdatedEvent",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToMessageUpdatedEvent(roomID)
}

func (s *loggingService) DeleteMessage(
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

func (s *loggingService) SubscribeToMessageDeletedEvent(roomID pulid.ID) (ch <-chan *ent.Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToMessageDeletedEvent",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToMessageDeletedEvent(roomID)
}
