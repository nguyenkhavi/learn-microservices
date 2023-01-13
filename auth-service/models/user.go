package models

import (
	"errors"
	"html"
	"strings"

	"nkvi/auth-service/utils/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func (u *User) SaveUser() (*User, error) {

	var err error
	err = GetDB().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, string, error) {

	var err error

	u := User{}
	err = GetDB().Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", "", err
	}

	access_token, err := token.GenerateToken(u.ID)

	refresh_token, err := token.GenerateRefreshToken(u.ID)

	if err != nil {
		return "", "", err
	}

	ut := UserToken{UserId: u.ID, Token: refresh_token}
	err = ut.SaveUserToken()

	if err != nil {
		return "", "", err
	}

	return access_token, refresh_token, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := GetDB().First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func RefreshToken(refreshToken string) (string, string, error) {

	var ID, err = token.ExtractRefreshTokenID(refreshToken)
	if err != nil {
		return "", "", err
	}

	ut := UserToken{}
	err = GetDB().Model(UserToken{}).Where(map[string]interface{}{"user_id": ID, "deleted_at": nil, "token": refreshToken}).Take(&ut).Debug().Error

	if err != nil {
		return "", "", err
	}
	access_token, err := token.GenerateToken(ID)
	refresh_token, err := token.GenerateRefreshToken(ID)

	if err != nil {
		return "", "", err
	}

	newUT := UserToken{UserId: ID, Token: refresh_token}

	ut.SoftDeleteUserToken()
	newUT.SaveUserToken()
	return access_token, refresh_token, nil
}

func Logout(refreshToken string) error {

	var ID, err = token.ExtractRefreshTokenID(refreshToken)
	if err != nil {
		return nil
	}

	ut := UserToken{}
	err = GetDB().Model(UserToken{}).Where("user_id = ?", ID).Take(&ut).Error
	if err != nil {
		return err
	}

	ut.SoftDeleteUserToken()

	return nil
}

func Block(userId int) error {
	err := GetDB().Where("user_id = ?", userId).Delete(&UserToken{}).Error
	if err != nil {
		return err
	}

	return nil
}
