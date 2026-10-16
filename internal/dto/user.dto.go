package dto

import "mime/multipart"

type UserUpdateRequest struct {
	Id        string                `form:"id"`
	FullName  string                `form:"fullname"`
	Phone     string                `form:"phone"`
	Email     string                `form:"email" binding:"omitempty,email"`
	AvatarUrl *multipart.FileHeader `form:"avatar_url"`
}

type UserUpdateResponse struct {
	Id        string  `json:"id"`
	FullName  string  `json:"fullname"`
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	AvatarUrl *string `json:"avatar_url"`
}
