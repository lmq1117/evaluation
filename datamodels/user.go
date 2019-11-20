package datamodels

import (
	"encoding/base64"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);unique_index"`
	Password string `gorm:"type:varchar(255)"`
}

// IsValid can do some very very simple "low-level" data validations.
func (u User) IsValid() bool {
	return u.ID > 0
}

// GeneratePassword will generate a hashed password for us based on the
// user's input.
func GeneratePassword(userPassword string) (string, error) {
	fmt.Println(userPassword)
	if hashedByte, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(hashedByte), err
	}
}

// ValidatePassword will check if passwords are matched.
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
