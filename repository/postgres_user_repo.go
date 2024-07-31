package repository

import (
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/rosberry/go-pagination"
	"gorm.io/gorm"
)

type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (repo *PostgresUserRepo) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	// if err := repo.db.Preload("Photo").Where("username = ?, is_active = ?", username, true).Error; err != nil{
	if err := repo.db.Preload("Photo").Where("username = ?", username).First(&user).Error; err != nil{
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresUserRepo) CreateUser(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *PostgresUserRepo) UpdateUser(user *models.User) error {
	return repo.db.Save(user).Error
}

func (repo *PostgresUserRepo) SetPreviewPhoto(userID uint, photoID uint) error {
	tx := repo.db.Begin()
	if err := tx.Model(&models.Photo{}).Where("user_id = ? AND id != ? AND is_preview = ?", userID, photoID, true).Update("is_preview", false).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&models.Photo{}).Where("id = ? AND user_id = ?", photoID, userID).Update("is_preview", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (repo *PostgresUserRepo) SaveLocation(username string, lat float32, lon float32) error {
	var user models.User
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	user.Lat = lat
	user.Lon = lon

	return repo.db.Save(&user).Error
}

func (repo *PostgresUserRepo) GetUsersWhoLikedMe(userID uint) ([]models.User, error) {
    var usersIdsWhichLikedMe []uint
    if err := repo.db.Model(&models.UserInteraction{}).Where("target_id = ?", userID).Pluck("user_id", &usersIdsWhichLikedMe).Error; err != nil {
        return nil, err
    }

    var usersWhichLikedMe []models.User
    if err := repo.db.Preload("Photo").Where("id IN (?)", usersIdsWhichLikedMe).Find(&usersWhichLikedMe).Error; err != nil {
        return nil, err
    }

    return usersWhichLikedMe, nil
}

func (repo *PostgresUserRepo) GetUsersList(userID uint, role string, paginator *pagination.Paginator) ([]models.User, error) {
	var curUser models.User
	config.DB.First(&curUser, userID)

	var users []models.User

	var interactedIDs []uint
	config.DB.Model(&models.UserInteraction{}).Where("user_id = ?", userID).Pluck("target_id", &interactedIDs)

	ageLower := curUser.Age - 3
	ageUpper := curUser.Age + 100
	gender := "male"
	if curUser.Sex == "male" {
		gender = "female"
	} else {
		gender = "male"
	}

	q := config.DB.Preload("Photo", "is_preview = ?", true).Model(&models.User{}).
		Where("role = ?", role).
		Where("id != ?", userID).
		Where("age BETWEEN ? and ?", ageLower, ageUpper).
		Where("sex = ?", gender).
		Where("id NOT IN (?)", interactedIDs)

	err := paginator.Find(q, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *PostgresUserRepo) AddUserInteraction(interaction *models.UserInteraction) error {
	return repo.db.Create(interaction).Error
}

func (repo *PostgresUserRepo) GetUserInteractionsCount(userID uint) (int64, error) {
	var userInteractionsCount int64
	if err := config.DB.Model(&models.UserInteraction{}).Where("user_id = ?", userID).Count(&userInteractionsCount).Error; err != nil {
		return 0, nil
	}
	return userInteractionsCount, nil
}