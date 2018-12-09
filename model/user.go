package model

import (
	"doc-manager/core/mongo"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

const (
	userCollection = "user"
)


type User struct {
	Id			bson.ObjectId 	`json:"_id" bson:"_id"`
	UserId 		string			`json:"user_id" bson:"user_id"`
	Account 	string  		//`json:"account"`  //登录账号（唯一）
	Password 	string
	Mail 		string			//
	Phone 		string
	Sex 		string
	RealName	string			`json:"real_name" bson:"real_name"`
	Title 		string  //职位
}

func (u User) FindAll() ([]User,error) {
	var results []User
	err := mongo.FindAll(db,userCollection,nil,nil,&results)
	return results,err
}

func (u User) Insert() error {
	if mongo.IsExist(db,userCollection,bson.M{"account": u.Account}) {
		return errors.New("account 已存在")
	}
	if mongo.IsExist(db,userCollection,bson.M{"mail": u.Mail}) {
		return errors.New("mail 已存在")
	}
	if u.UserId == "" {
		u.Id = bson.NewObjectId()
		u.UserId = u.Id.Hex()
	}
	return mongo.Insert(db,userCollection,u)
}

func (u User) Update(id string) error {
	if !strings.EqualFold(u.UserId,id) {
		return errors.New("userId not equal id")
	}
	return mongo.Update(db,userCollection,bson.M{"_id": bson.ObjectIdHex(id)},u)
}

func (u User) Remove(id string) error {
	return mongo.Remove(db,userCollection,bson.M{"_id": bson.ObjectIdHex(id)})
}

func (u User) FindById(id string) (*User, error) {
	err := mongo.FindOne(db,userCollection,bson.M{"_id": bson.ObjectIdHex(id)},nil,&u)
	return &u,err
}
func (u User)UserLogin(account, password string)(*User,error)  {
	err := mongo.FindOne(db, userCollection,bson.M{"account": account,"password":password},nil,&u)
	return &u,err
}
