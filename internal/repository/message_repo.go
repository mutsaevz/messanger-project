

package repository

import (
	"context"
	"log/slog"

	"gorm.io/gorm"

	"test-messenger/internal/models"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	err := r.db.WithContext(ctx).Create(message).Error
	if err != nil {
		slog.Error(
			"failed to create message",
			"chat_id", message.ChatID,
			"error", err,
		)
		return err
	}

	slog.Info(
		"message created",
		"message_id", message.ID,
		"chat_id", message.ChatID,
	)
	return nil
}

func (r *MessageRepository) GetLastByChatID(
	ctx context.Context,
	chatID uint,
	limit int,
) ([]models.Message, error) {
	var messages []models.Message

	err := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at desc").
		Limit(limit).
		Find(&messages).
		Error

	if err != nil {
		slog.Warn(
			"failed to get messages by chat id",
			"chat_id", chatID,
			"limit", limit,
			"error", err,
		)
		return nil, err
	}

	slog.Debug(
		"messages loaded",
		"chat_id", chatID,
		"messages_count", len(messages),
	)

	return messages, nil
}

func (r *MessageRepository) DeleteByChatID(ctx context.Context, chatID uint) error {
	err := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Delete(&models.Message{}).
		Error

	if err != nil {
		slog.Error(
			"failed to delete messages by chat id",
			"chat_id", chatID,
			"error", err,
		)
		return err
	}

	slog.Info(
		"messages deleted by chat id",
		"chat_id", chatID,
	)
	return nil
}