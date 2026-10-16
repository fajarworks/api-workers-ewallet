package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (u *UserController) UpdateUserProfile(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "token doesn't exist",
		})
		return
	}
	claims, ok := token.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "invalid token",
		})
		return
	}
	userId := claims.Id

	var payload dto.UserUpdateRequest
	if err := ctx.ShouldBindWith(&payload, binding.FormMultipart); err != nil {
		log.Println("[UpdateProfile] Bind error:", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var targetFilePath string
	var avatarPtr *string

	if payload.AvatarUrl != nil {
		const maxFileSize = 2 * 1024 * 1024
		if payload.AvatarUrl.Size > maxFileSize {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "file exceeds size limit (2mb)",
			})
			return
		}

		ext := strings.ToLower(filepath.Ext(payload.AvatarUrl.Filename))
		allowedExt := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		if !allowedExt[ext] {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "unsupported format. Please upload in .jpg, .jpeg, or .png format",
			})
			return
		}

		uploadDir := filepath.Join("public", "images", "profile")
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Println("[UpdateProfile] failed to create directory:", err.Error())
			ctx.JSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Message: "failed processing file",
			})
			return
		}

		filename := fmt.Sprintf("%s_%d%s", userId, time.Now().UnixNano(), ext)
		targetFilePath = filepath.Join(uploadDir, filename)
		if err := ctx.SaveUploadedFile(payload.AvatarUrl, targetFilePath); err != nil {
			log.Println("[UpdateProfile] save upload file err:", err.Error())
			ctx.JSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Message: "failed to upload file",
			})
			return
		}

		imageUrl := fmt.Sprintf("images/profile/%s", filename)
		avatarPtr = &imageUrl
	}

	account, err := u.service.UpdateProfileInfo(
		ctx.Request.Context(),
		userId,
		payload.Email,
		dto.UserUpdateResponse{
			Id:        userId,
			FullName:  payload.FullName,
			Phone:     payload.Phone,
			AvatarUrl: avatarPtr,
		},
	)
	if err != nil {
		log.Println("[UpdateProfile] service err:", err.Error())
		if targetFilePath != "" {
			if rmErr := os.Remove(targetFilePath); rmErr != nil {
				log.Println("[UpdateProfile] failed to remove file:", rmErr.Error())
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "failed to update profile",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "profile updated successfully",
		Data:    account,
	})
}
