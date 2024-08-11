package models

import (
	"gorm.io/gorm"
)


type Chat struct {
	gorm.Model
	User1ID uint `gorm:"not null" json:"user1_id"`
	User2ID uint `gorm:"not null" json:"user2_id"`

	User1 User `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE;" json:"user1,omitempty"`
	User2 User `gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE;" json:"user2,omitempty"`
}

func (Chat) TableName() string {
	return "chats"
}


type Message struct {
	gorm.Model
	ChatID     uint   `gorm:"not null" json:"chat_id"`
	SenderID   uint   `gorm:"not null" json:"sender_id"`
	ReceiverID uint   `gorm:"not null" json:"receiver_id"`
	Content    string `gorm:"type:text;not null" json:"content"`
	IsRead     bool   `gorm:"default:false" json:"is_read"`

	Chat     Chat `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE;" json:"chat,omitempty"`
	Sender   User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE;" json:"sender,omitempty"`
	Receiver User `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE;" json:"receiver,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}
