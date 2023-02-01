package mysql

import (
	"github.com/charfole/simple-tiktok/model"
)

// 创建一个todo：直接传入一个结构体进行创建，有错误则返回
func CreateATodo(todo *model.Todo) (err error) {
	err = DB.Create(&todo).Error
	return
}

// 获取所有todo：传入todo结构体指针数组todoList，获取当前todos表中的所有数据
// 同时更新todoList，有错误则返回
func GetAllTodo() (todoList []*model.Todo, err error) {
	if err = DB.Find(&todoList).Error; err != nil {
		return nil, err
	}
	return
}

// 获取某个todo：通过id获取某个todo事项，有错误则返回
func GetATodo(id string) (todo *model.Todo, err error) {
	todo = new(model.Todo)
	if err = DB.Debug().Where("id=?", id).First(todo).Error; err != nil {
		return nil, err
	}
	return
}

// 更新某个todo：传入一个todo结构体指针，更新todos表中对应的todo项
func UpdateATodo(todo *model.Todo) (err error) {
	err = DB.Save(todo).Error
	return
}

// 删除某个todo：传入一个id，删除对应id的todo项
func DeleteATodo(id string) (err error) {
	err = DB.Where("id=?", id).Delete(&model.Todo{}).Error
	return
}
