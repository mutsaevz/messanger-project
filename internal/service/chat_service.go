package service

import (
	"context"
	"errors"
	"slices"
	"time"

	"gorm.io/gorm"

	"test-messenger/internal/dto"
	"test-messenger/internal/models"
	"test-messenger/internal/repository"
)

var (
	ErrChatNotFound = errors.New("chat not found")
)

type ChatService struct {
	chats    repository.ChatRepository
	messages repository.MessageRepository
}

func NewChatService(
	chats repository.ChatRepository,
	messages repository.MessageRepository,
) *ChatService {
	return &ChatService{
		chats:    chats,
		messages: messages,
	}
}

func (s *ChatService) CreateChat(
	ctx context.Context,
	req *dto.CreateChatRequest,
) (*dto.ChatResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	chat := &models.Chat{
		Title: req.Title,
	}

	if err := s.chats.Create(ctx, chat); err != nil {
		return nil, err
	}

	return &dto.ChatResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *ChatService) GetChat(
	ctx context.Context,
	chatID uint,
	limit int,
) (*models.Chat, error) {
	chat, err := s.chats.GetByIDWithMessages(ctx, chatID, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, err
	}

	slices.Reverse(chat.Messages)

	return chat, nil
}

func (s *ChatService) DeleteChat(ctx context.Context, chatID uint) error {
	_, err := s.chats.GetByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrChatNotFound
		}
		return err
	}

	return s.chats.Delete(ctx, chatID)
}
