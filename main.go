package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type User struct {
	Id          uint64    `xorm:"pk autoincr"` // pk=primary key, autoincr=オートインクリメント
	Name        string    `xorm:"varchar(40)"`
	Description *string   // stringポインタにすれば空文字やnilを書き込める
	EquipId     uint16    `xorm:"smallint"` // 型対応表: https://gist.github.com/suzujun/6a6db0019f1c14c0ec9ff0ea4a495451
	Version     uint64    `xorm:"version"`  // カラム名は変わらないが自動入力される
	CreatedAt   time.Time `xorm:"created"`  // カラム名は変わらないが自動入力される
	UpdatedAt   time.Time `xorm:"updated"`  // カラム名は変わらないが自動入力される
	DeletedAt   time.Time `xorm:"deleted"`  // deletedがあれば論理削除、なければ物理削除
}

func (u User) IsValid() bool {
	return u.Id != 0
}

func (User) TableName() string {
	return "user"
}

func Migrate(engine *xorm.Engine) {
	_ = engine.Sync2(new(User))
}

func Create(engine *xorm.Engine, user *User) {
	affected, err := engine.Insert(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	if affected == 0 {
		fmt.Println("No Affected")
		return
	}
}

func Read(engine *xorm.Engine, id uint64) User {
	var res User

	isExist, _ := engine.Where("id=?", id).Get(&res)
	if !isExist {
		fmt.Println("Not Found")
		res = User{
			Id:          0,
			Name:        "",
			Description: nil,
			EquipId:     0,
			Version:     0,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			DeletedAt:   time.Time{},
		}
	}

	return res
}

func UpdateName(engine *xorm.Engine, id uint64, newName string) {
	newUser := Read(engine, id)
	if !newUser.IsValid() {
		return
	}

	newUser.Name = newName

	result, err := engine.Where("id=?", id).Update(&newUser)
	if err != nil {
		fmt.Println(err)
	}
	if result == 0 {
		fmt.Println("Not Found")
	}
}

func UpdateDescription(engine *xorm.Engine, id uint64, description string) {
	newUser := Read(engine, id)
	if !newUser.IsValid() {
		return
	}

	newUser.Description = &description

	result, err := engine.Where("id=?", id).Update(&newUser)
	if err != nil {
		fmt.Println(err)
	}
	if result == 0 {
		fmt.Println("Not Found")
	}
}

func Delete(engine *xorm.Engine, id uint64) {
	user := User{}
	result, err := engine.Where("id=?", id).Delete(&user)
	if err != nil {
		fmt.Println(err)
	}
	if result == 0 {
		fmt.Println("Not Found")
	}
}

func DeleteAll(engine *xorm.Engine) {
	var user []User
	_ = engine.Find(&user)
	for _, v := range user {
		_, err := engine.Where("id=?", v.Id).Delete(User{})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Print(engine *xorm.Engine, id uint64) {
	user := Read(engine, id)
	if !user.IsValid() {
		return
	}

	fmt.Println(user)
}

func PrintAll(engine *xorm.Engine) {
	var user []User
	_ = engine.Find(&user)
	for _, v := range user {
		Print(engine, v.Id)
	}
}

func main() {
	engine, err := xorm.NewEngine("mysql", "root:root@tcp([127.0.0.1]:3306)/hoge?charset=utf8mb4&parseTime=true")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()

	// 実行されたSQL文を標準出力
	// engine.ShowSQL(true)

	Migrate(engine)

    fmt.Println("before:")
    PrintAll(engine)

	// description := new(string)
	// *description = "説明"
	// user := User{
	// 	Name:        "太郎3",
	// 	Description: description,
	// 	EquipId:     100,
	// }
	// Create(engine, &user)

	// UpdateName(engine, 2, "花子")
	// UpdateName(engine, 2, "") // フィールドがゼロ値の場合は更新されない

	// UpdateDescription(engine, 2, "説明文書き換えテスト")

	// Delete(engine, 3)
	// DeleteAll(engine)

	fmt.Println("after:")
	PrintAll(engine)
}
