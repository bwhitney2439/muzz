package types

import "github.com/bwhitney2439/muzz/models"

type UserResponse struct {
	ID                  uint            `json:"id"`
	Name                string          `json:"name"`
	Age                 uint8           `json:"age"`
	Gender              string          `json:"gender"`
	Location            models.Location `gorm:"embedded" json:"-"`
	DistanceFromMe      float64         `json:"distanceFromMe"`
	AttractivenessScore int             `json:"attractivenessScore"`
}
