package dto

type CreateUser struct {
	Mobile       string
	PasswordHash string
	RoleId       int64
}

type UpdateUser struct {
	Mobile       string
	PasswordHash string
	RoleId       int64
}

type User struct {
	BaseModel
	Mobile       string
	PasswordHash string
	RoleId       int64
}

type UserResponse struct {
	Mobile       string `json:"mobile"`
	PasswordHash string
	Role         Role
	RoleId       int64
}
