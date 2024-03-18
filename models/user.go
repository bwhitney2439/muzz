package models

import (
	"gorm.io/gorm"
)

type Location struct {
	Latitude  float64
	Longitude float64
}
type User struct {
	gorm.Model
	Email               string   `json:"email"`
	Password            string   `json:"password"`
	Name                string   `json:"name"`
	Gender              string   `json:"gender"`
	Age                 uint8    `json:"age"`
	Location            Location `gorm:"embedded" json:"location,omitempty"`
	AttractivenessScore int      `json:"attractivenessScore,omitempty"`
}

type Swipe struct {
	gorm.Model
	UserID       uint   `json:"userId"`       // The ID of the user who swiped
	TargetUserID uint   `json:"targetUserId"` // The ID of the user being swiped on
	Preference   string `json:"preference"`   // YES or NO
}
