package chat

import (
	"context"
	"time"

	"journeyhub/ent"
	"journeyhub/ent/schema/pulid"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type subscriptionsLogging struct {
	logger log.Logger
	Subscriptions
}

func NewSubscriptionsLogging(logger log.Logger, s Subscriptions) Subscriptions {
	return &subscriptionsLogging{logger, s}
}

func (s *subscriptionsLogging) PublishMessageCreatedEvent(
	ctx context.Context,
	roomID pulid.ID,
	messageID pulid.ID,
) (subject string, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "PublishMessageCreatedEvent",
			"subject", subject,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.PublishMessageCreatedEvent(ctx, roomID, messageID)
}

func (s *subscriptionsLogging) SubscribeToMessageCreatedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (ch <-chan *ent.MessageEdge, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToMessageCreatedEvent",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.SubscribeToMessageCreatedEvent(ctx, roomID)
}

func (s *subscriptionsLogging) PublishMessageUpdatedEvent(
	ctx context.Context,
	roomID pulid.ID,
	messageID pulid.ID,
) (subject string, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "PublishMessageUpdatedEvent",
			"subject", subject,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.PublishMessageUpdatedEvent(ctx, roomID, messageID)
}

func (s *subscriptionsLogging) SubscribeToMessageUpdatedEvent(
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
	return s.Subscriptions.SubscribeToMessageUpdatedEvent(ctx, roomID)
}

func (s *subscriptionsLogging) PublishMessageDeletedEvent(
	ctx context.Context,
	roomID pulid.ID,
	messageID pulid.ID,
) (subject string, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "PublishMessageDeletedEvent",
			"subject", subject,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.PublishMessageDeletedEvent(ctx, roomID, messageID)
}

func (s *subscriptionsLogging) SubscribeToMessageDeletedEvent(
	ctx context.Context,
	roomID pulid.ID,
) (ch <-chan pulid.ID, err error) {
	defer func(begin time.Time) {
		level.Debug(s.logger).Log(
			"method", "SubscribeToMessageDeletedEvent",
			"roomID", roomID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Subscriptions.SubscribeToMessageDeletedEvent(ctx, roomID)
}
