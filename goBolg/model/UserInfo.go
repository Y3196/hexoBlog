package model

import (
	"time"

	"gorm.io/gorm"
)

// UserInfo represents a user entity as stored in the database.
type UserInfo struct {
	ID         int       `gorm:"primaryKey;autoIncrement"` // User ID
	Email      string    `gorm:"not null"`                 // Email address of the user
	Nickname   string    `gorm:"size:255"`                 // User's nickname
	Avatar     string    `gorm:"size:255"`                 // URL to the user's avatar
	Intro      string    `gorm:"size:1024"`                // Short user introduction
	WebSite    string    `gorm:"size:255"`                 // Personal website URL
	IsDisable  int       `gorm:"default:0"`                // Indicates if the user is disabled (0 for false, 1 for true)
	CreateTime time.Time `gorm:"autoCreateTime"`           // Timestamp of creation
	UpdateTime time.Time `gorm:"autoUpdateTime"`           // Timestamp of the last update

}

// Here's how you could potentially set up this model with auto migration in your main application or during setup.
// db is assumed to be a *gorm.DB instance connected to your database.

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&UserInfo{})
	if err != nil {
		return err
	}
	return nil
}
func (UserInfo) TableName() string {
	return "tb_user_info" // 指定正确的表名
}
