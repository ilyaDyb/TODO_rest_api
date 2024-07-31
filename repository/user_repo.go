package repository

import (
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/rosberry/go-pagination"
)

type UserRepo interface {
	GetUserByUsername(username string) (*models.User, error)
    CreateUser(user *models.User) error
    UpdateUser(user *models.User) error
    SetPreviewPhoto(userID uint, photoID uint) error
    SaveLocation(username string, lat float32, lon float32) error
    GetUsersWhoLikedMe(userID uint) ([]models.User, error)
    GetUsersList(userID uint, role string, paginator *pagination.Paginator) ([]models.User, error)
    AddUserInteraction(interaction *models.UserInteraction) error
    GetUserInteractionsCount(userID uint) (int64, error)
}