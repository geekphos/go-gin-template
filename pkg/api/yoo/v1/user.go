package v1

type LoginRequest struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=6,max=14"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6,max=14"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=14"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"email,required"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=14"`
}
