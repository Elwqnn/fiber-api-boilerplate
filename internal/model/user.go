package model

// User model gather every information about a user
type User struct {
	BaseModel
	Name          string    `json:"name" gorm:"not null" validate:"required,min=3,max=100"`
	Email         string    `json:"email" gorm:"unique;not null" validate:"required,email"`
	Image         string    `json:"image"`
	Role          string    `json:"role" gorm:"not null;default:user" validate:"required,oneof=admin user"`
	Phone         string    `json:"phone"`
	Accounts      []Account `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
