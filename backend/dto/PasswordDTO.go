package dto

type PasswordChange struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required, min=7"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
