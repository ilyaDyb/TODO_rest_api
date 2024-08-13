package repository

import (
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/utils"
)

type ChatRepo interface {
	//for admin controllers
	GetAllChats() (*[]models.Chat, error)

	//for other controllers
	CreateChat(chat *models.Chat) error 
	CreateMessage(message *models.Message) error

	GetChatByUsernames(username1, username2 string) (*models.Chat, error)
	GetMessagesByIDChat(chatID uint) (*[]models.Message, error)
	GetUserChats(userID uint) (*[]utils.ChatsListResponse, error)
	
	
	// GetChatsForSpecUser(userID uint) ([]struct {
	// 	User              models.User
	// 	LastMessage       string
	// 	IsLastMessageRead bool
	// 	LastMessageTime   time.Time
	// }, error)

}

// func NewChatRepo(db *gorm.DB) ChatRepo {
// 	return &chatRepo{db: db}
// }

// func (r *chatRepo) CreateChat(chat *models.Chat) error {
// 	return r.db.Create(chat).Error
// }

// func (r *chatRepo) GetAllChats() (*[]models.Chat, error) {
// 	return r.GetAllChats()
// }

// func (r *chatRepo) GetChat(user1ID, user2ID uint) (*models.Chat, error) {
// 	return r.GetChat(user1ID, user2ID)
// }