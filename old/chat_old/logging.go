package chat

import (
	"context"
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
	currentUserID ID,
	input SendMessageInput,
) (msg *Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SendMessage",
			"currentUserID", currentUserID,
			"targetUserID", input.TargetUserID,
			"replyTo", input.ReplyTo,
			"content", input.Content,
			"filesCount", len(input.Files),
			"message", msg,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendMessage(ctx, currentUserID, input)
}

func (s *loggingService) SubscribeToMessageAddedEvent(
	roomID ID,
) (ch <-chan *Message, err error) {
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
	currentUserID ID,
	messageID ID,
	input UpdateMessageInput,
) (msg *Message, err error) {
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
	return s.Service.UpdateMessage(ctx, currentUserID, messageID, input)
}

func (s *loggingService) SubscribeToMessageUpdatedEvent(
	roomID ID,
) (ch <-chan *Message, err error) {
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
	currentUserID ID,
	messageID ID,
) (msg *Message, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "DeleteMessage",
			"messageID", messageID,
			"message", msg,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteMessage(ctx, currentUserID, messageID)
}

func (s *loggingService) SubscribeToMessageDeletedEvent(
	roomID ID,
) (ch <-chan *Message, err error) {
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

func (s *loggingService) SubscribeToRoomsChangedEvent(
	currentUserID ID,
) (ch <-chan *Room, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToRoomsChangedEvent",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SubscribeToRoomsChangedEvent(currentUserID)
}
