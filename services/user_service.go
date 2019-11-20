package services

import (
	"encoding/base64"
	"evaluation/datamodels"
	"github.com/jinzhu/gorm"
)

var user datamodels.User

// UserService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type UserService interface {
	GetUserList(offset, limit int) []*datamodels.User
	Create(user datamodels.User) datamodels.User
	GetByID(id int64) (datamodels.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool)
}

// NewUserService returns the default user service.
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

type userService struct {
	db *gorm.DB
}

func (s *userService) GetUserList(offset, limit int) []*datamodels.User {
	var userList []*datamodels.User
	s.db.Find(&userList)
	return userList
}

func (s *userService) Create(user datamodels.User) datamodels.User {
	s.db.Create(&user)
	if user.ID > 0 {
		return user
	} else {
		return datamodels.User{}
	}
}

func (s *userService) GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool) {

	if username == "" || userPassword == "" {
		return user, false
	}
	s.db.Where("username = ?", username).First(&user)

	passwordBytes, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		//
	}

	if ok, _ := datamodels.ValidatePassword(userPassword, passwordBytes); ok {
		return user, true
	}
	return user, false

}

func (s *userService) GetByID(id int64) (datamodels.User, bool) {
	s.db.Where("id = ?", id).First(&user)
	if user.ID > 0 {
		return user, true
	} else {
		return user, false
	}
}
