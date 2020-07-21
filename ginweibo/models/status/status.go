package status

import (
	"ginweibo/models/database"
	"ginweibo/models"

	"fmt"
	userModel "ginweibo/models/user"
	"strconv"

	"github.com/lexkong/log"
)

const tableName = "statuses"

// Status 微博
type Status struct {
	models.BaseModel
	Content string `gorm:"column:context;type:text;not null"`
	UserID  uint   `gorm:"column:user_id;not null" sql:"index"` // 一对多，关联 User Model
}

// TableName 表名
func (Status) TableName() string {
	return tableName
}

func (s *Status) Create() (err error) {
	if err = database.DB.Create(&s).Error; err != nil {
		log.Warnf("微博创建失败: %v", err)
		return err
	}
	return nil
}

func Delete(id int) (err error) {
	status := &Status{}
	status.BaseModel.ID = uint(id)
	if err = database.DB.Delete(&status).Error; err != nil {
		log.Warnf("微博删除失败: %v", err)
		return err
	}
	return nil
}

func Get(id int) (*Status, error) {
	s := &Status{}
	d := database.DB.First(&s, id)
	return s, d.Error
}

// 获取指定用户们的微博数量
func GetByUsersStatusesCount(ids []uint) (int, error) {
	sqlStr := fmt.Sprintf("select count(*) from %s where deleted_at is null and user_id in (", tableName)
	l := len(ids) - 1
	for i, v := range ids {
		sqlStr = sqlStr + strconv.Itoa(int(v))
		if i < l {
			sqlStr = sqlStr + ","
		}
	}
	sqlStr = sqlStr + ")"
	count := 0
	d := database.DB.Raw(sqlStr).Count(&count)
	return count, d.Error
}

// 获取指定用户们的微博
func GetByUsersStatuses(ids []uint, offset, limit int) ([]*Status, error) {
	status := make([]*Status, 0)
	sqlStr := fmt.Sprintf("select * from %s where deleted_at is null and user_id in (", tableName)
	l := len(ids) - 1
	for i, v := range ids {
		sqlStr = sqlStr + strconv.Itoa(int(v))
		if i < l {
			sqlStr = sqlStr + ","
		}
	}
	sqlStr = sqlStr + fmt.Sprintf(") order by `created_at` desc limit %d offset %d", limit, offset)
	d := database.DB.Raw(sqlStr).Scan(&status)
	return status, d.Error
}

// 通过 status_id 获取该微博的所有者
func GetUser(statusID int) (*userModel.User, error) {
	s, err := Get(statusID)
	if err != nil {
		return nil, err
	}
	u, err := userModel.Get(int(s.UserID))
	if err != nil {
		return nil, err
	}
	return u, nil
}

// 获取该用户的所有微博
func GetUserAllStatus(userID int) ([]*Status, error) {
	status := make([]*Status, 0)
	err := database.DB.Where("user_id = ?", userID).Order("id desc").Find(&status).Error
	if err != nil {
		return status, err
	}
	return status, nil
}

// 获取该用户的微博 (分页)
func GetUserStatus(userID, offset, limit int) ([]*Status, error) {
	status := make([]*Status, 0)
	err := database.DB.Where("user_id = ?", userID).Offset(
		offset).Limit(limit).Order("id desc").Find(&status).Error
	if err != nil {
		return status, err
	}
	return status, nil
}

// 获取该用户的所有微博 的 count
func GetUserAllStatusCount(userID int) (count int, err error) {
	err = database.DB.Model(&Status{}).Where("user_id = ?", userID).Count(&count).Error
	return
}
