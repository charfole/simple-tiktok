/*
与dao层中的数据库进行交互，负责数据库中某张表的增删改查等基本操作
通过gorm，将type与MySQL的某张表进行映射，通过操作type来操作数据表
*/
package models

import "github.com/charfole/simple-tiktok/dao"

// Todo Model 默认映射到当前数据库中的todos表
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Status bool   `json:"status"`
}

// 创建一个todo：直接传入一个结构体进行创建，有错误则返回
func CreateATodo(todo *Todo) (err error) {
	err = dao.DB.Create(&todo).Error
	return
}

// 获取所有todo：传入todo结构体指针数组todoList，获取当前todos表中的所有数据
// 同时更新todoList，有错误则返回
func GetAllTodo() (todoList []*Todo, err error) {
	if err = dao.DB.Find(&todoList).Error; err != nil {
		return nil, err
	}
	return
}

// 获取某个todo：通过id获取某个todo事项，有错误则返回
func GetATodo(id string) (todo *Todo, err error) {
	todo = new(Todo)
	if err = dao.DB.Debug().Where("id=?", id).First(todo).Error; err != nil {
		return nil, err
	}
	return
}

// 更新某个todo：传入一个todo结构体指针，更新todos表中对应的todo项
func UpdateATodo(todo *Todo) (err error) {
	err = dao.DB.Save(todo).Error
	return
}

// 删除某个todo：传入一个id，删除对应id的todo项
func DeleteATodo(id string) (err error) {
	err = dao.DB.Where("id=?", id).Delete(&Todo{}).Error
	return
}
