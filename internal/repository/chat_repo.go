package repository

import (
	"context"
	"log/slog"

	"gorm.io/gorm"

	"test-messenger/internal/models"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *models.Chat) error {
	err := r.db.WithContext(ctx).Create(chat).Error
	if err != nil {
		slog.Error("failed to create chat", "error", err)
		return err
	}

	slog.Info("chat created", "chat_id", chat.ID)
	return nil
}

func (r *ChatRepository) GetByID(ctx context.Context, id uint) (*models.Chat, error) {
	var chat models.Chat

	err := r.db.WithContext(ctx).First(&chat, id).Error
	if err != nil {
		slog.Warn("failed to get chat by id", "chat_id", id, "error", err)
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) GetByIDWithMessages(
	ctx context.Context,
	id uint,
	limit int,
) (*models.Chat, error) {
	var chat models.Chat

	err := r.db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc").Limit(limit)
		}).
		First(&chat, id).
		Error

	if err != nil {
		slog.Warn(
			"failed to get chat with messages",
			"chat_id", id,
			"limit", limit,
			"error", err,
		)
		return nil, err
	}

	slog.Debug(
		"chat loaded with messages",
		"chat_id", chat.ID,
		"messages_count", len(chat.Messages),
	)

	return &chat, nil
}

func (r *ChatRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&models.Chat{}, id).Error
	if err != nil {
		slog.Error("failed to delete chat", "chat_id", id, "error", err)
		return err
	}

	slog.Info("chat deleted", "chat_id", id)
	return nil
}
