package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	RoleSysadmin = "super"
	RoleAdmin    = "admin"
	RoleMember   = "member"
)

const (
	StatusValid = iota
	StatusDisabled
)

// LoginInfo 登录信息
type LoginInfo struct {
	SessionToken string    `json:"session_token,omitempty"`
	LastAlive    time.Time `json:"last_alive,omitempty"`
	IP           string    `json:"ip,omitempty"`
}

// Profile 用户资料
type Profile struct {
	Remark  string `json:"remark,omitempty"`
	Address string `json:"address,omitempty"`
	Email   string `json:"email,omitempty"`
}

// User 用户信息
type User struct {
	gorm.Model
	Name     string                `json:"name,omitempty" gorm:"unique"`
	NickName string                `json:"nick_name"`
	Password string                `json:"password,omitempty"`
	Role     string                `json:"role"`
	Status   int                   `json:"status"`
	Profile  Profile               `json:"profile,omitempty" gorm:"embedded"`
	Session  map[string]*LoginInfo `json:"session,omitempty" gorm:"-"`
}

func NewUser(name, password string) (u *User) {
	u = new(User)
	u.Name = name
	u.NickName = name
	u.UpdatePassword(password)
	u.Role = RoleMember
	return
}

func (u *User) UpdatePassword(newPassword string) (err error) {
	var bytes []byte
	if bytes, err = bcrypt.GenerateFromPassword([]byte(newPassword), 9); err != nil {
		return
	}
	u.Password = string(bytes)
	return
}

func (u *User) VerifyPassword(password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return
}
