package model

import (
	"database/sql"
	"github.com/killsari/api/common"
	"github.com/sirupsen/logrus"
	"time"
)

type RegionModel struct {
}

type RegionDao struct {
	Id    uint64    `db:"f_id"`
	Code  uint64    `db:"f_code"`
	Name  string    `db:"f_name"`
	Ctime time.Time `db:"f_ctime"`
	Utime time.Time `db:"f_utime"`
}

type RegionShow struct {
	Code uint64 `json:"code"`
	Name string `json:"name"`
}

func (m RegionModel) Get(code uint64) (dao *RegionDao, err error) {
	dao = &RegionDao{}
	err = common.DB.Get(dao, "select * from t_region where f_code = ?", code)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (m RegionModel) Select() (daoArr *[]RegionDao, err error) {
	daoArr = &[]RegionDao{}
	err = common.DB.Select(daoArr, "select * from t_region order by f_id")
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (m RegionModel) Show(dao *RegionDao) (show []*RegionShow, err error) {
	show = []*RegionShow{}

	var showTmp []*RegionShow
	showTmp = append(showTmp, &RegionShow{
		Code: dao.Code,
		Name: dao.Name,
	})

	cityCode := dao.Code / 100 * 100
	provinceCode := dao.Code / 10000 * 10000

	if cityCode != dao.Code && cityCode != provinceCode {
		var cityDao *RegionDao
		cityDao, err = m.Get(cityCode)
		if err != nil {
			if err != sql.ErrNoRows {
				return
			}
		} else {
			showTmp = append(showTmp, &RegionShow{
				Code: cityDao.Code,
				Name: cityDao.Name,
			})
		}
	}

	if provinceCode != dao.Code && provinceCode != cityCode {
		var provinceDao *RegionDao
		provinceDao, err = m.Get(provinceCode)
		if err != nil {
			if err != sql.ErrNoRows {
				return
			}
		} else {
			showTmp = append(showTmp, &RegionShow{
				Code: provinceDao.Code,
				Name: provinceDao.Name,
			})
		}
	}

	len := len(showTmp)
	for i := 0; i < len; i++ {
		show = append(show, showTmp[len-1-i])
	}

	return
}

func (m RegionModel) ShowByCode(code uint64) (show []*RegionShow, err error) {
	show = []*RegionShow{}

	dao, err := m.Get(code)
	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = nil
		return
	}
	show, err = m.Show(dao)
	if err != nil {
		return
	}
	return
}
