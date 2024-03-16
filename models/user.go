package models

import (
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      uint8  `json:"age"`
}

type Swipe struct {
	gorm.Model
	UserID       uint   // The ID of the user who swiped
	TargetUserID uint   // The ID of the user being swiped on
	Preference   string // YES or NO
}
