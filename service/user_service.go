package service

import (
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/repository"
	"github.com/rosberry/go-pagination"
)

type UserService struct {
    repo repository.UserRepo
}
 
func NewUserService(repo repository.UserRepo) UserService {
    return UserService{repo: repo}
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
    return s.repo.GetUserByUsername(username)
}

func (s *UserService) GetUserByID(ID uint) (*models.User, error) {
    return s.repo.GetUserByID(ID)
}

func (s *UserService) CreateUser(user *models.User) error {
    return s.repo.CreateUser(user)
}

func (s *UserService) UpdateUser(user *models.User) error {
    return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(user *models.User) error {
    return s.repo.DeleteUser(user)
}


func (s *UserService) SetPreviewPhoto(userID uint, photoID uint) error {
    return s.repo.SetPreviewPhoto(userID, photoID)
}

func (s *UserService) SaveLocation(username string, lat float32, lon float32) error {
    return s.repo.SaveLocation(username, lat, lon)
}

func (s *UserService) GetUsersWhoLikedMe(userID uint) ([]models.User, error) {
    return s.repo.GetUsersWhoLikedMe(userID)
}

func (s *UserService) GetUsersList(userID uint, role string, paginator *pagination.Paginator) ([]models.User, error) {
    return s.repo.GetUsersList(userID, role, paginator)
}

func (s *UserService) AddUserInteraction(interaction *models.UserInteraction) error {
    return s.repo.AddUserInteraction(interaction)
}

func (s *UserService) GetUserInteractionsCount(userID uint) (int64, error) {
    return s.repo.GetUserInteractionsCount(userID)    
}

func (s *UserService) UserIsExists(username string, email string) (bool, error) {
    return s.repo.UserIsExists(username, email)
}

func (s *UserService) GetUserByHash(hash string) (*models.User, error) {
    return s.repo.GetUserByHash(hash)
}

func (s *UserService) GetAllUsers(limit int, page int) ([]models.User, error) {
    return s.repo.GetAllUsers(limit, page)
}

func (s *UserService) GetUsersCount() (int, error) {
    return s.repo.GetUsersCount()
}

func (s *UserService) IsExistsEmail(email string) (bool, error) {
    return s.repo.IsExistsEmail(email)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
    return s.repo.GetUserByEmail(email)
}