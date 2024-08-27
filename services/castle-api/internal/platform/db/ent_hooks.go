package db

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"entgo.io/ent"
)

func LoggingHook(logger log.Logger) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			defer func() {
				level.Debug(logger).Log("op", m.Op(), "type", m.Type())
			}()
			return next.Mutate(ctx, m)
		})
	}
}
