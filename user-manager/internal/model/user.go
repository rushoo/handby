package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

/*
	一种方法是自定义类型重写序列化方法，这里是在marshall前将结构体类型做了转换
*/

var Gender = map[bool]string{
	false: "Male",
	true:  "Female",
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  bool   `json:"sex"`
}

type UserStore struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

var UserCache map[string]*UserStore

func init() {
	UserCache = make(map[string]*UserStore)
	//os.Create("user.json")
}
func load(file string) error {
	log.Printf("从 %s 中加载用户信息", file)
	fileData, err := os.ReadFile(file)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(fileData) > 0 {
		//序列化空文件：unexpected end of JSON input
		return json.Unmarshal(fileData, &UserCache)
	}
	return nil
}
func write(user *User, file string) error {
	u := &UserStore{
		Name: user.Name,
		Age:  user.Age,
		Sex:  Gender[user.Sex],
	}
	UserCache[user.Name] = u
	data, _ := json.Marshal(UserCache)
	return os.WriteFile(file, data, 0644)
}

func (user *User) Get() (*UserStore, error) {
	err := load("user.json")
	if err != nil {
		return &UserStore{}, err
	}
	log.Println("将新的用户信息加入，然后一起持久化")
	if u, ok := UserCache[user.Name]; !ok {
		return &UserStore{}, errors.New("用户不存在")
	} else {
		return u, nil
	}
}
func (user *User) List() ([]*UserStore, error) {
	err := load("user.json")
	if err != nil {
		return nil, err
	}

	var users []*UserStore
	for _, u := range UserCache {
		users = append(users, u)
	}
	return users, nil

}
func (user *User) UpdateUser() error {
	err := load("user.json")
	if err != nil {
		return err
	}
	return write(user, "user.json")
}
func (user *User) CreateUser() error {
	if user.Age < 0 {
		log.Println("校验用户年龄：失败")
		return errors.New("age cannot be negative")
	}
	err := load("user.json")
	if err != nil {
		log.Println(err)
		return err
	}
	return write(user, "user.json")
}
func (user *User) DeleteUser() error {
	err := load("user.json")
	if err != nil {
		return err
	}
	log.Println("判断用户是否存在")
	if _, ok := UserCache[user.Name]; !ok {
		return errors.New(fmt.Sprintf("user %s not exist", user.Name))
	}
	delete(UserCache, user.Name)
	data, _ := json.Marshal(UserCache)
	return os.WriteFile("user.json", data, 0644)
}
func (user *User) CheckExist() (error, bool) {
	err := load("user.json")
	if err != nil {
		return err, false
	}
	_, ok := UserCache[user.Name]
	return nil, ok
}
