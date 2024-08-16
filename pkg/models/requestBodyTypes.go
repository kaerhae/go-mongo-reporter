package models

type CreateUser struct {
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Permission Permission `json:"permission"`
}

type CreateGuestUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserPasswordChange struct {
	UserID      string `json:"userId"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserPermissionUpdate struct {
	UserID      string     `json:"userId"`
	Permissions Permission `json:"permissions"`
}
