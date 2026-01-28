package mappers

import (
	"github.com/zouipo/yumsday/backend/internal/dtos"
	"github.com/zouipo/yumsday/backend/internal/models"
)

// Maps a User model to a UserDtoNoPassword DTO
func ToUserDtoNoPassword(user *models.User) *dtos.UserDto {
	return &dtos.UserDto{
		ID:        user.ID,
		Username:  user.Username,
		AppAdmin:  user.AppAdmin,
		CreatedAt: user.CreatedAt,
		Avatar:    user.Avatar,
		Language:  user.Language,
		AppTheme:  user.AppTheme,
	}
}

// Maps a NewUserDto DTO to a User model (only received from the client when creating a new user)
func ToModelFromNewUserDto(newUserDto *dtos.NewUserDto) *models.User {
	return &models.User{
		ID:       newUserDto.ID,
		Username: newUserDto.Username,
		Password: newUserDto.Password,
		AppAdmin: newUserDto.AppAdmin,
		Avatar:   newUserDto.Avatar,
		Language: newUserDto.Language,
		AppTheme: newUserDto.AppTheme,
	}
}

func ToModelFromUserDto(userDto *dtos.UserDto) *models.User {
	return &models.User{
		ID:        userDto.ID,
		Username:  userDto.Username,
		AppAdmin:  userDto.AppAdmin,
		CreatedAt: userDto.CreatedAt,
		Avatar:    userDto.Avatar,
		Language:  userDto.Language,
		AppTheme:  userDto.AppTheme,
	}
}
