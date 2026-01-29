
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apphttp "test-messenger/internal/http"
	"test-messenger/internal/http/handler"
	"test-messenger/internal/models"
	"test-messenger/internal/repository"
	"test-messenger/internal/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp(t *testing.T) http.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Chat{}, &models.Message{})
	require.NoError(t, err)

	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	chatService := service.NewChatService(*chatRepo, *messageRepo)
	messageService := service.NewMessageService(*chatRepo, *messageRepo)

	chatHandler := handler.NewChatHandler(chatService)
	messageHandler := handler.NewMessageHandler(messageService)

	router := apphttp.NewRouter(chatHandler, messageHandler)
	return router.Handler()
}

func TestCreateChat(t *testing.T) {
	app := setupTestApp(t)

	body := map[string]string{
		"title": "Test Chat",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chats", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	}

	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.NotZero(t, resp.ID)
	assert.Equal(t, "Test Chat", resp.Title)
}

func TestCreateChat_ValidationError(t *testing.T) {
	app := setupTestApp(t)

	body := map[string]string{
		"title": "   ",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chats", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
