package dto

import (
	"time"
)

type LineUser struct {
	Id        string    `json:"id" bson:"id"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
