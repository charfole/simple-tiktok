package common

import "errors"

var (
	ErrorDBMigrateFail   = errors.New("数据库迁移错误")
	ErrorSQLFalse        = errors.New("SQL执行错误")
	ErrorCreateUserFalse = errors.New("新建用户错误")
	ErrorTokenFalse      = errors.New("token不正确")
	ErrorExpired         = errors.New("token已过期")
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNameEmpty   = errors.New("用户名为空")
	ErrorUserNameInvalid = errors.New("用户名长度应少于32位！")
	ErrorPasswordEmpty   = errors.New("密码为空")
	ErrorPasswordInvalid = errors.New("密码长度应至少为6位且不超过32位！")
)
