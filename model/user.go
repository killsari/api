package model

import (
	"github.com/killsari/api/common"
	"github.com/sirupsen/logrus"
	"time"
)

type UserModel struct {
}

type UserDao struct {
	Id    uint64    `db:"f_id"`
	Tel   string    `db:"f_tel"`
	Ctime time.Time `db:"f_ctime"`
	Utime time.Time `db:"f_utime"`
}

type UserShow struct {
	Tel string `json:"tel"`
}

const (
	UserOptionRegLogin = "reg_login"
)

func (m UserModel) Get(id uint64) (dao *UserDao, err error) {
	dao = &UserDao{}
	err = common.DB.Get(dao, "select * from t_user where f_id = ?", id)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (m UserModel) GetByTel(tel string) (dao *UserDao, err error) {
	dao = &UserDao{}
	err = common.DB.Get(dao, "select * from t_user where f_tel = ?", tel)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (m UserModel) Put(userDao *UserDao) (err error) {
	// 用户基本信息
	result, err := common.DB.Exec("insert into t_user (f_tel) "+
		"values (?) "+
		"on duplicate key update "+
		"f_tel=values(f_tel) ",
		userDao.Tel)
	if err != nil {
		return
	}
	if userDao.Id == 0 {
		var lastInsertId int64
		lastInsertId, err = result.LastInsertId()
		if err != nil {
			logrus.Error(err)
			return
		}
		userDao.Id = uint64(lastInsertId)
	}
	return
}

func (m UserModel) Show(dao *UserDao, uid uint64, brief bool) (show *UserShow, err error) {
	show = &UserShow{}
	show.Tel = dao.Tel
	// brief
	return
}

func (m UserModel) ShowById(id uint64, uid uint64, brief bool) (show *UserShow, err error) {
	show = &UserShow{}
	dao, err := UserModel{}.Get(id)
	if err != nil {
		return
	}
	show, err = m.Show(dao, uid, brief)
	if err != nil {
		return
	}
	return
}
