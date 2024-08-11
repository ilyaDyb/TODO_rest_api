package repository

import (
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/utils"
	"gorm.io/gorm"
)

type PostgresChatRepo struct {
	db *gorm.DB
}

func NewPostgresChatRepo(db *gorm.DB) *PostgresChatRepo {
	return &PostgresChatRepo{db: db}
}


func (repo *PostgresChatRepo) CreateChat(user *models.Chat) error {
	return repo.db.Create(user).Error
}

func (repo *PostgresChatRepo) GetAllChats() (*[]models.Chat, error) {
	var chats []models.Chat
	if err := repo.db.Model(models.Chat{}).
	Preload("User1.Photo", "is_preview = ?", true).
	Preload("User2.Photo", "is_preview = ?", true).
	Find(&chats).Error; err != nil {
		return nil, err
	}
	return &chats, nil
}

func (repo *PostgresChatRepo) GetChatByUsernames(username1, username2 string) (*models.Chat, error) {
	var chat models.Chat
	err := repo.db.Model(&models.Chat{}).
		Joins("JOIN users u1 ON u1.id = chats.user1_id").
		Joins("JOIN users u2 ON u2.id = chats.user2_id").
		Where("(u1.username = ? AND u2.username = ?) OR (u1.username = ? AND u2.username = ?)", username1, username2, username2, username1).
		Preload("User1.Photo", "is_preview = ?", true).
		Preload("User2.Photo", "is_preview = ?", true).
		First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}


func (repo *PostgresChatRepo) GetMessagesByIDChat(chatID uint) (*[]models.Message, error) {
	var messages []models.Message
	if err := repo.db.Model(&models.Message{}).Where("chat_id = ?", chatID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return &messages, nil
}

// func (repo *PostgresChatRepo) GetChatsForSpecUser(userID uint) (*[]models.Chat, error) {
// 	var chats []models.Chat
// 	if err := repo.db.Model(&models.Chat{}).
// 	Preload("User1.Photo", "is_preview = ?", true).
// 	Preload("User2.Photo", "is_preview = ?", true).
// 	Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&chats).Error; err != nil {
// 		return nil, err
// 	}
// 	return &chats, nil
// }

// func (repo *PostgresChatRepo) GetChatsForSpecUser(userID uint) ([]struct {
// 	User              models.User
// 	LastMessage       string
// 	IsLastMessageRead bool
// 	LastMessageTime   time.Time
// }, error) {
// 	var results []struct {
// 		User              models.User
// 		LastMessage       string
// 		IsLastMessageRead bool
// 		LastMessageTime   time.Time
// 	}
	
// 	query := `
// 		SELECT 
// 			u.*,
// 			last_messages.content AS last_message,
// 			last_messages.is_read AS is_last_message_read,
// 			last_messages.created_at AS last_message_time
// 		FROM 
// 			users u
// 		JOIN (
// 			SELECT 
// 				c.id AS chat_id,
// 				CASE
// 					WHEN c.user1_id = ? THEN c.user2_id
// 					ELSE c.user1_id
// 				END AS other_user_id,
// 				messages.content,
// 				messages.is_read,
// 				messages.created_at AS last_message_time
// 			FROM 
// 				chats c
// 			JOIN (
// 				SELECT 
// 					m1.chat_id, 
// 					m1.content, 
// 					m1.is_read, 
// 					m1.created_at
// 				FROM 
// 					messages m1
// 				WHERE 
// 					(m1.chat_id, m1.created_at) IN (
// 						SELECT 
// 							m2.chat_id, 
// 							MAX(m2.created_at) 
// 						FROM 
// 							messages m2
// 						GROUP BY 
// 							m2.chat_id
// 					)
// 			) messages ON messages.chat_id = c.id
// 			WHERE 
// 				c.user1_id = ? OR c.user2_id = ?
// 		) last_messages ON u.id = last_messages.other_user_id
// 		ORDER BY 
// 			last_messages.last_message_time DESC
// 	`
	
// 	err := repo.db.Raw(query, userID, userID, userID).Scan(&results).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return results, nil
// }

func (repo *PostgresChatRepo) GetUserChats(userID uint) (*[]utils.ChatsListResponse, error) {
	var results []utils.ChatsListResponse

	query := `
	SELECT users.firstname, users.lastname, users.username, photos.url, messages.content AS last_message, messages.is_read, chats.id AS chat_id
		FROM chats
		JOIN users ON (users.id = chats.user1_id OR users.id = chats.user2_id)
		LEFT JOIN photos ON photos.user_id = users.id AND photos.is_preview = true
		LEFT JOIN 
			messages ON messages.chat_id = chats.id AND messages.created_at = (
			SELECT MAX(created_at) 
			FROM messages 
			WHERE chat_id = chats.id
		)
		WHERE (chats.user1_id = ? OR chats.user2_id = ?) AND users.id != ?;
	`
	err := repo.db.Raw(query, userID, userID, userID).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}