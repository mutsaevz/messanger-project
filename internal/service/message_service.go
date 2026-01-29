package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"test-messenger/internal/dto"
	"test-messenger/internal/models"
	"test-messenger/internal/repository"
)

var (
	ErrChatNotFoundForMessage = errors.New("chat not found for message")
)

type MessageService struct {
	chats    repository.ChatRepository
	messages repository.MessageRepository
}

func NewMessageService(
	chats repository.ChatRepository,
	messages repository.MessageRepository,
) *MessageService {
	return &MessageService{
		chats:    chats,
		messages: messages,
	}
}

func (s *MessageService) CreateMessage(
	ctx context.Context,
	chatID uint,
	req *dto.CreateMessageRequest,
) (*dto.MessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	_, err := s.chats.GetByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatNotFoundForMessage
		}
		return nil, err
	}

	message := &models.Message{
		ChatID: chatID,
		Text:   req.Text,
	}

	if err := s.messages.Create(ctx, message); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		ID:        message.ID,
		ChatID:    message.ChatID,
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
	}, nil
}
