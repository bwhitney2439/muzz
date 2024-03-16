package database

import (
	"fmt"
	"log"

	"github.com/bwhitney2439/muzz/models"
	"github.com/bwhitney2439/muzz/types"
	"github.com/bwhitney2439/muzz/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	dburl = utils.GoDotEnvVariable("DB_URL")
	db    *gorm.DB
)

func Connect() {
	var err error
	dbConn := sqlite.Open(dburl)
	db, err = gorm.Open(dbConn, &gorm.Config{})
	if err != nil {
		log.Fatal()
	}

	fmt.Println("Connected with Database", dburl)
	err = db.AutoMigrate(&models.User{}, &models.Swipe{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database migrated successfully")

}

func InsertUser(user *models.User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByEmail(email string) (*models.User, bool, error) {
	user := new(models.User)
	err := db.Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}

func GetUser(userResponse *types.UserResponse, userID uint) error {
	err := db.Model(&models.User{}).Where("id = ?", userID).Select("id", "name", "age", "gender").Scan(userResponse).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUsers(usersResponse *[]types.UserResponse, excludeUserID float64) error {
	err := db.Model(&models.User{}).Where("id <> ?", uint(excludeUserID)).Select("id", "name", "age", "gender").Scan(usersResponse).Error
	if err != nil {
		return err
	}
	return nil
}

func GetPotentialMatches(userID uint) ([]models.User, error) {
	var users []models.User

	subQuery := db.Model(&models.Swipe{}).Select("target_user_id").Where("user_id = ?", userID)

	err := db.Model(&models.User{}).
		Where("id <> ?", userID).
		Where("id NOT IN (?)", subQuery).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func SwipeAction(userID, targetUserID uint, preference string) (matched bool, matchedUser *types.UserResponse, err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return false, nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	defer func() {
		if err == nil {
			err = tx.Commit().Error // Capture commit error
		}
		if err != nil {
			tx.Rollback() // Only rollback if not already done and there's an error
		}
	}()

	swipe := models.Swipe{
		UserID:       userID,
		TargetUserID: targetUserID,
		Preference:   preference,
	}
	if err := tx.Create(&swipe).Error; err != nil {
		return false, nil, err
	}

	if preference != "YES" {
		return false, nil, nil // Early return if preference is not YES, no need to rollback explicitly due to defer
	}

	var count int64
	err = tx.Model(&models.Swipe{}).
		Where("user_id = ? AND target_user_id = ? AND preference = 'YES'", targetUserID, userID).
		Count(&count).Error
	if err != nil {
		return false, nil, err
	}

	if count > 0 {
		fmt.Println("Match found!")
		matchedUser = new(types.UserResponse)
		err := tx.Model(&models.User{}).Where("id = ?", userID).Select("id", "name", "age", "gender").Scan(matchedUser).Error
		if err != nil {
			return true, nil, err
		}

		return true, matchedUser, nil
	}

	return false, nil, nil
}
