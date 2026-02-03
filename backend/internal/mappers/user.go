package mappers

import (
	"github.com/zouipo/yumsday/backend/internal/dtos"
	"github.com/zouipo/yumsday/backend/internal/models"
)

// ToUserDtoNoPassword maps a User model to a UserDto without the password field.
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

// ToModelFromNewUserDto maps a NewUserDto to a User model (used when creating a new user).
func FromNewUserDtoToUser(newUserDto *dtos.NewUserDto) *models.User {
	return &models.User{
		Username: newUserDto.Username,
		Password: newUserDto.Password,
		AppAdmin: newUserDto.AppAdmin,
		Avatar:   newUserDto.Avatar,
		Language: newUserDto.Language,
		AppTheme: newUserDto.AppTheme,
	}
}

// ToModelFromUserDto maps a UserDto to a User model (omits password field).
func FromUserDtoToUser(userDto *dtos.UserDto) *models.User {
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
