package tgc

import (
	"context"

	"github.com/divyam234/teldrive/internal/logger"
	"github.com/gotd/td/telegram"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func RunWithAuth(ctx context.Context, client *telegram.Client, token string, f func(ctx context.Context) error) error {
	return client.Run(ctx, func(ctx context.Context) error {
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return err
		}

		if token == "" {
			if !status.Authorized {
				return errors.Errorf("not authorized. please login first")
			}
			logger.Logger.Info("User Session",
				zap.Int64("id", status.User.ID),
				zap.String("username", status.User.Username))
		} else {
			if !status.Authorized {
				logger.Logger.Info("creating bot session")
				_, err := client.Auth().Bot(ctx, token)
				if err != nil {
					return err
				}
				status, _ = client.Auth().Status(ctx)
				logger.Logger.Info("Bot Session",
					zap.Int64("id", status.User.ID),
					zap.String("username", status.User.Username))
			}
		}

		return f(ctx)
	})
}
