package models

// User model gather every information about a user
type User struct {
	BaseModel
	Name     string    `json:"name" validate:"required,min=3,max=100" gorm:"not null"`
	Email    string    `json:"email" validate:"required,email" gorm:"unique;not null;index"`
	Image    string    `json:"image"`
	Role     string    `json:"role" validate:"required,oneof=admin user" gorm:"not null;default:'user';index"`
	Phone    string    `json:"phone"`
	Accounts []Account `json:"accounts" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
