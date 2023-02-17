package common

import "errors"

var (
	ErrorDBMigrateFail      = errors.New("数据库迁移错误")
	ErrorSQLFalse           = errors.New("SQL执行错误")
	ErrorCreateUserFalse    = errors.New("新建用户错误")
	ErrorTokenFalse         = errors.New("token不正确")
	ErrorExpired            = errors.New("token已过期")
	ErrorUserExist          = errors.New("用户已存在")
	ErrorUserNotFound       = errors.New("用户不存在")
	ErrorUserNameEmpty      = errors.New("用户名为空")
	ErrorUserNameInvalid    = errors.New("用户名长度应少于32位！")
	ErrorPasswordEmpty      = errors.New("密码为空")
	ErrorPasswordInvalid    = errors.New("密码长度应至少为6位且不超过32位！")
	ErrorFullPossibility    = errors.New("账号或密码出错！")
	ErrorNullPointer        = errors.New("空指针异常")
	ErrorPasswordFalse      = errors.New("密码错误！")
	ErrorVideoList          = errors.New("获取视频列表出错！")
	ErrorVideoUpload        = errors.New("上传视频出错！")
	ErrorCOSUpload          = errors.New("COS服务出错！")
	ErrorVideoDBCreateFalse = errors.New("视频数据创建失败！")
	ErrorRelationExit       = errors.New("关注已存在")
	ErrorRelationNull       = errors.New("关注不存在")
)
