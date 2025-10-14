package usecase

import (
	"time"

	"dialog-service/internal/repository"

	"github.com/google/uuid"
)

type DialogUseCase struct {
	dialogRepo *repository.DialogRepository
}

func NewDialogUseCase(dialogRepo *repository.DialogRepository) *DialogUseCase {
	return &DialogUseCase{
		dialogRepo: dialogRepo,
	}
}

type SendMessageRequest struct {
	FromUserID string
	ToUserID   string
	Text       string
}

type DialogMessageResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func (uc *DialogUseCase) SendMessage(req *SendMessageRequest) error {
	// Validate required fields
	if req.Text == "" {
		return ErrInvalidData
	}

	// Create message
	message := &repository.DialogMessage{
		ID:         uuid.New().String(),
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Text:       req.Text,
		CreatedAt:  time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
	}

	return uc.dialogRepo.CreateMessage(message)
}

func (uc *DialogUseCase) GetMessagesBetweenUsers(userID1, userID2 string) ([]DialogMessageResponse, error) {
	// Get messages
	messages, err := uc.dialogRepo.GetMessagesBetweenUsers(userID1, userID2)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var response []DialogMessageResponse
	for _, msg := range messages {
		response = append(response, DialogMessageResponse{
			From: msg.FromUserID,
			To:   msg.ToUserID,
			Text: msg.Text,
		})
	}

	return response, nil
}
