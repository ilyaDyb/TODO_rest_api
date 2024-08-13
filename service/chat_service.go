package service

import (
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/repository"
	"github.com/ilyaDyb/go_rest_api/utils"
)

type ChatService struct {
    repo repository.ChatRepo
}
 
func NewChatService(repo repository.ChatRepo) ChatService {
    return ChatService{repo: repo}
}

func (s *ChatService) CreateChat(chat *models.Chat) error {
	return s.repo.CreateChat(chat)
}

func (s *ChatService) CreateMessage(message *models.Message) error {
	return s.repo.CreateMessage(message)
}

func (s *ChatService) GetAllChats() (*[]models.Chat, error) {
	return s.repo.GetAllChats()
}

func (s *ChatService) GetChatByUsernames(username1, username2 string) (*models.Chat, error) {
	return s.repo.GetChatByUsernames(username1, username2)
}

func (s *ChatService) GetMessagesByIDChat(chatID uint) (*[]models.Message, error) {
	return s.repo.GetMessagesByIDChat(chatID)
}

func (s *ChatService) GetUserChats(userID uint) (*[]utils.ChatsListResponse, error) {
	return s.repo.GetUserChats(userID)
}