package Utils

import (
	"sync"
)

type UserPassInfo struct {
	Username string `json:"username"`
	Password string `json:"password"` //保存哈希
}

type UserManager struct {
	UserList map[string]UserPassInfo
	lock     sync.RWMutex
}

func (u *UserManager) New() *UserManager {
	u.UserList = make(map[string]UserPassInfo)
	return u
}

func InitUserManager() *UserManager {
	userManager := UserManager{}
	userManager.New()
	return &userManager
}

func (u *UserManager) IsExist(username string) bool {
	u.lock.Lock()
	defer u.lock.Unlock()
	_, ok := u.UserList[username]
	return ok
}

func (u *UserManager) Check(username string, password string) bool {
	if u.IsExist(username) {
		u.lock.Lock()
		defer u.lock.Unlock()
		hash := GetStringMD5(password)
		passIn := u.UserList[username].Password
		return hash == passIn
	} else {
		return false
	}
}

func (u *UserManager) AddUser(username string, password string) bool {
	if u.IsExist(username) {
		return false
	} else {
		userInfo := UserPassInfo{
			Username: username,
			Password: "",
		}
		passHash := GetStringMD5(password)
		userInfo.Password = passHash
		u.UserList[username] = userInfo
		return true
	}
}
