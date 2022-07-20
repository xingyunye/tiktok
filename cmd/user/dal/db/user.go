package db

import (
	"douyin/v1/kitex_gen/user"
	"douyin/v1/pkg/constants"
	"douyin/v1/pkg/errno"

	"gorm.io/gorm"

	"context"
)

type User struct {
	gorm.Model
	ID            int64  `gorm:"primarykey" json:"user_id"`
	UserName      string `gorm:"not null" json:"user_name"`       // not null and repeat
	Password      string `gorm:"not null" json:"password"`        // md5加密后的密码
	FollowCount   int64  `gorm:"default:0" json:"follow_count"`   // 关注数
	FollowerCount int64  `gorm:"default:0" json:"follower_count"` // 粉丝数
}

func (u *User) TableName() string {
	return constants.UserTableName
}

// CreateUser 创建用户信息
func CreateUser(ctx context.Context, users []*User) error {
	return DB.WithContext(ctx).Create(users).Error
}

// QueryUser 通过username查询用户
func QueryUser(ctx context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("user_name = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// QueryUser 通过userid查询用户
func QueryUserById(ctx context.Context, userId int64) (User, error) {
	res := User{}
	if err := DB.WithContext(ctx).Where("id = ?", userId).First(&res).Error; err != nil {
		return User{}, err
	}
	return res, nil
}

// MGetUsers 根据一个Userid获取其关注粉丝数量
func MGetUsers(ctx context.Context, req *user.MGetUserRequest) ([]*User, error) {
	res := make([]*User, 0)
	if req.ToUserId < 1 || req.ActionType < constants.QueryUserInfo || req.ActionType > constants.QueryFollowerList {
		return nil, errno.ParamErr
	}
	if req.ActionType == constants.QueryUserInfo {
		if err := DB.WithContext(ctx).Where("id = ?", req.ToUserId).Find(&res).Error; err != nil {
			return nil, err
		}
		return res, nil
	} else if req.ActionType == constants.QueryFollowList {
		followIds, err := QueryFollowById(ctx, req.ToUserId)
		if err != nil {
			return nil, err
		}
		if err = DB.WithContext(ctx).Where("id in ?", followIds).Find(&res).Error; err != nil {
			return nil, err
		}
		return res, nil
	} else {
		followerIds, err := QueryFollowerById(ctx, req.ToUserId)
		if err != nil {
			return nil, err
		}
		if err = DB.WithContext(ctx).Where("id in ?", followerIds).Find(&res).Error; err != nil {
			return nil, err
		}
		return res, nil
	}
}

//QueryFollowRelation 查询是否关注了userId对应的用户
func QueryFollowRelation(ctx context.Context, users []*User, userId int64) ([]bool, error) {
	isFollowList := make([]bool, len(users))
	if userId == constants.NotLogin {
		for i := 0; i < len(users); i++ {
			isFollowList[i] = false
		}
	} else {
		for i, user := range users {
			var temp int64 = 0
			DB.WithContext(ctx).Model(&Follower{}).Where("user_id = ? and follower_id = ?", user.ID, userId).Count(&temp)
			if temp > 0 {
				isFollowList[i] = true
			} else {
				isFollowList[i] = false
			}
		}
	}
	return isFollowList, nil
}

//GetUserInfoList 根据userIDs获取用户信息列表
func GetUserInfoList(ctx context.Context, userIDs []int64) ([]*User, error) {
	var res []*User
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := DB.WithContext(ctx).Where("id in ?", userIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

//UpdateUser 关注或者取关后进行的数据库操作
func UpdateUser(ctx context.Context, req *user.UpdateUserRequest) error {
	if req.UserId == constants.NotLogin {
		return nil
	} //查询用户是否存在

	// 如果要关注 查询是否已经是关注状态
	//如果要取关 查询是否已经是取关状态
	var cnt int64 = 0
	if err := DB.WithContext(ctx).Model(&Follower{}).Where("user_id = ? and follower_id = ?", req.ToUserId, req.UserId).Count(&cnt).Error; err != nil {
		return err
	}

	if req.ActionType == constants.RelationAdd {
		if cnt > 0 {
			return nil
		}
	} else if req.ActionType == constants.RelationDel {
		if cnt == 0 {
			return nil
		}
	}

	//查询两个用户是否存在
	var user1 User
	var user2 User

	if err := DB.WithContext(ctx).Where("id = ?", req.UserId).First(&user1).Error; err != nil {
		return err
	}

	if err := DB.WithContext(ctx).Where("id = ?", req.ToUserId).First(&user2).Error; err != nil {
		return err
	}

	//使用事务封装
	return DB.Transaction(func(tx *gorm.DB) error {
		//先在Follow表中更改关注的关系
		//再在User表中更改follow_count与follower_count
		if req.ActionType == constants.RelationAdd {
			if err := tx.WithContext(ctx).Create(&Follower{UserID: req.ToUserId, FollowerID: req.UserId}).Error; err != nil {
				return err
			}
			user1.FollowCount += 1
			user2.FollowerCount += 1
		} else if req.ActionType == constants.RelationDel {
			if err := tx.WithContext(ctx).Where("user_id = ? and follower_id = ?", req.ToUserId, req.UserId).Delete(&Follower{}).Error; err != nil {
				return err
			}
			user1.FollowCount -= 1
			user2.FollowerCount -= 1
		}
		if err := tx.WithContext(ctx).Model(&user1).Select("follow_count").Updates(user1).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Model(&user2).Select("follower_count").Updates(user2).Error; err != nil {
			return err
		}
		return nil //没有错误则提交事务
	})
}
