package database

import (
	"painter-server-new/models"
	"painter-server-new/tolog"
	"time"
)

func CreateFollow(followerID, followingID int) error {
	follow := models.FollowTable{
		FollowerID:  followerID,
		FollowingID: followingID,
		FollowTime:  time.Now(),
	}

	err := DbEngine.Create(&follow).Error
	if err != nil {
		tolog.Log().Infof("Error while create follow %e", err)
		return err
	}
	return nil
}

func DeleteFollow(followerID, followingID int) error {
	err := DbEngine.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.FollowTable{}).Error
	if err != nil {
		tolog.Log().Infof("Error while delete follow %e", err)
		return err
	}
	return nil
}

func GetFollowers(userID, limit, offset int) ([]models.FollowTable, error) {
	var followers []models.FollowTable
	result := DbEngine.Where("following_id = ?", userID).Limit(limit).Offset(offset).Find(&followers)
	if result.Error != nil {
		return nil, result.Error
	}
	return followers, nil
}

func GetFollowings(userID, limit, offset int) ([]models.FollowTable, error) {
	var followings []models.FollowTable
	result := DbEngine.Where("follower_id = ?", userID).Limit(limit).Offset(offset).Find(&followings)
	if result.Error != nil {
		return nil, result.Error
	}
	return followings, nil
}

func GetFollowingsUsers(followings []models.FollowTable) ([]models.UserTable, error) {
	var users []models.UserTable
	for _, following := range followings {
		var user models.UserTable
		result := DbEngine.Where("id = ?", following.FollowingID).First(&user)
		if result.Error != nil {
			return nil, result.Error
		}
		users = append(users, user)
	}
	return users, nil
}

func GetFollowerNumber(userID int) (int, error) {
	var count int64
	result := DbEngine.Model(&models.FollowTable{}).Where("following_id = ?", userID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

func GetFollowingNumber(userID int) (int, error) {
	var count int64
	result := DbEngine.Model(&models.FollowTable{}).Where("follower_id = ?", userID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

func CheckFollow(a, b int) (bool, error) {
	var count int64
	result := DbEngine.Model(&models.FollowTable{}).Where("follower_id = ? AND following_id = ?", b, a).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
