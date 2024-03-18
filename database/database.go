package database

import (
	"fmt"
	"log"
	"sort"
	"strings"

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

func GetPotentialMatches(userID uint, age *uint, gender, orderBy string) (*[]types.UserResponse, error) {

	user := new(models.User)

	err := db.Model(&models.User{}).Where("id = ?", userID).Scan(user).Error
	if err != nil {
		return nil, err
	}

	userLong := user.Location.Longitude
	userLat := user.Location.Latitude

	// get list of userids where the user has already swiped
	subQuery := db.Model(&models.Swipe{}).Select("target_user_id").Where("user_id = ?", userID)

	// get list of users who have not been swiped by the user using subquery
	query := db.Model(&models.User{}).Where("id <> ?", userID).Where("id NOT IN (?)", subQuery).Select("id", "name", "age", "gender", "latitude", "longitude", "attractiveness_score")

	if age != nil {
		query = query.Where("age = ?", *age)
	}

	if gender != "" {
		query = query.Where("LOWER(gender) = ?", strings.ToLower(gender))
	}

	if orderBy == "attractiveness_score" {
		query = query.Order("attractiveness_score desc")
	}

	matchedUsers := new([]types.UserResponse)

	err = query.Scan(matchedUsers).Error
	if err != nil {
		return nil, err
	}

	// Calculate the distance from the user.
	for i := range *matchedUsers {
		(*matchedUsers)[i].DistanceFromMe = utils.Haversine(userLat, userLong, (*matchedUsers)[i].Location.Latitude, (*matchedUsers)[i].Location.Longitude)
	}
	if orderBy == "distance" {

		sort.Slice(*matchedUsers, func(i, j int) bool {
			return (*matchedUsers)[i].DistanceFromMe < (*matchedUsers)[j].DistanceFromMe
		})
	}

	return matchedUsers, nil
}

func BeginTransaction(db *gorm.DB) *gorm.DB {
	return db.Begin()
}

func CheckExistingSwipe(tx *gorm.DB, userID, targetUserID uint) (bool, error) {
	var count int64
	tx.Model(&models.Swipe{}).Where("user_id = ? AND target_user_id = ?", userID, targetUserID).Count(&count)
	return count > 0, nil
}

func CreateSwipe(tx *gorm.DB, userID, targetUserID uint, preference string) error {
	swipe := models.Swipe{
		UserID:       userID,
		TargetUserID: targetUserID,
		Preference:   preference,
	}
	return tx.Create(&swipe).Error
}

func CheckForMatch(tx *gorm.DB, userID, targetUserID uint) (bool, error) {
	var count int64
	err := tx.Model(&models.Swipe{}).
		Where("user_id = ? AND target_user_id = ? AND preference = 'YES'", targetUserID, userID).
		Count(&count).Error
	return count > 0, err
}

func GetMatchedUser(tx *gorm.DB, userID uint) (*types.UserResponse, error) {
	matchedUser := new(types.UserResponse)
	err := tx.Model(&models.User{}).Where("id = ?", userID).Select("id", "name", "age", "gender").Scan(matchedUser).Error
	return matchedUser, err
}

func UpdateAttractivenessScore(tx *gorm.DB, targetUserID uint, preference string) error {
	var increment int
	if preference == "YES" {
		increment = 1
	} else {
		increment = 0
	}

	err := tx.Model(&models.User{}).Where("id = ?", targetUserID).
		Update("attractiveness_score", gorm.Expr("attractiveness_score + ?", increment)).Error

	return err
}

func SwipeAction(userID, targetUserID uint, preference string) (bool, *types.UserResponse, error) {
	tx := BeginTransaction(db)
	if tx.Error != nil {
		return false, nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	exists, err := CheckExistingSwipe(tx, userID, targetUserID)
	if err != nil {
		return false, nil, err
	}
	if exists {
		return false, nil, fmt.Errorf("swipe already exists")
	}

	if err := CreateSwipe(tx, userID, targetUserID, preference); err != nil {
		return false, nil, err
	}

	if preference != "YES" {
		return false, nil, nil
	}

	err = UpdateAttractivenessScore(tx, targetUserID, preference)
	if err != nil {
		return false, nil, err
	}

	matched, err := CheckForMatch(tx, userID, targetUserID)
	if err != nil || !matched {
		return matched, nil, err
	}

	matchedUser, err := GetMatchedUser(tx, targetUserID)
	return matched, matchedUser, err
}
