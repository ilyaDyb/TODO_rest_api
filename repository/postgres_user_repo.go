package repository

import (
	"fmt"
	"log"
	"sort"

	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/utils"
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
	if err := repo.db.Preload("Photo").Where("username = ? AND is_active = ?", username, true).First(&user).Error; err != nil{
	// if err := repo.db.Preload("Photo").Where("username = ?", username).First(&user).Error; err != nil{
		return nil, err
	}
	// if user.ID == 0 {
	// 	return &user, fmt.Errorf("user's status is unactive")
	// }
	return &user, nil
}

func (repo *PostgresUserRepo) GetUserByID(ID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, ID).Error; err != nil {
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

func (repo *PostgresUserRepo) DeleteUser(user *models.User) error {
	return repo.db.Delete(user).Error
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

	gender := "male"
	if curUser.Sex == "male" {
		gender = "female"
	} else {
		gender = "male"
	}

	q := config.DB.Preload("Photo", "is_preview = ?", true).Model(&models.User{}).
		Where("role = ?", role).
		Where("id != ?", userID).
		Where("sex = ?", gender).
		Where("id NOT IN (?)", interactedIDs).Limit(10)

	err := paginator.Find(q, &users)
	if err != nil {
		return nil, err
	}
	scores := make(map[uint]float64)
	for _, u := range users {
		scores[u.ID] = utils.CalculateScore(curUser, u)
	}
	sort.Slice(users, func(i, j int) bool {
		return scores[users[i].ID] > scores[users[j].ID]
	})
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

func (repo *PostgresUserRepo) GetUserInteraction(userID, targetID uint) (*models.UserInteraction, error) {
	var userInteraction models.UserInteraction
	if err := config.DB.Model(&models.UserInteraction{}).Where("user_id = ? AND target_id = ?").First(&userInteraction).Error; err != nil {
		return nil, err
	}
	return &userInteraction, nil 
}

func (repo *PostgresUserRepo) UserIsExists(username string, email string) (bool, error){
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = ? OR email = ?)"
	var exists bool
	if err := config.DB.Raw(query, username, email).Scan(&exists).Error; err != nil {
		return exists, err
	}
	return exists, nil
}

func (repo *PostgresUserRepo) GetUserByHash(hash string) (*models.User, error) {
    var user models.User
    if err := repo.db.Where("confirmation_hash = ?", hash).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (repo *PostgresUserRepo) GetUsersCount() (int, error) {
    var count int64
    if err := repo.db.Model(&models.User{}).Count(&count).Error; err != nil {
        return 0, err
    }
    return int(count), nil
}

func (repo *PostgresUserRepo) GetAllUsers(limit int, offset int) ([]models.User, error) {
    var users []models.User
    if err := repo.db.Limit(limit).Offset(offset).Order("id DESC").Find(&users).Error; err != nil {
        return users, fmt.Errorf("err")
    }
    return users, nil
}

func (repo *PostgresUserRepo) IsExistsEmail(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)"
	log.Printf("Executing query: %s with email: %s\n", query, email)
	if err := config.DB.Raw(query, email).Scan(&exists).Error; err != nil {
		log.Printf("Error executing query: %v\n", err)
		return exists, err
	}
	log.Printf("Email exists: %v\n", exists)
	return exists, nil
}

func (repo *PostgresUserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	log.Printf("Finding user by email: %s\n", email)
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("Error finding user by email: %v\n", err)
		return nil, err
	}
	log.Printf("Found user: %v\n", user)
	return &user, nil
}