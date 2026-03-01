package mapper

import (
	"github.com/zouipo/yumsday/backend/internal/dto"
	"github.com/zouipo/yumsday/backend/internal/model"
)

// ToUserDtoNoPassword maps a User model to a UserDto without the password field.
func ToUserDtoNoPassword(user *model.User) *dto.UserDto {
	return &dto.UserDto{
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
func FromNewUserDtoToUser(newUserDto *dto.NewUserDto) *model.User {
	return &model.User{
		Username: newUserDto.Username,
		Password: newUserDto.Password,
		AppAdmin: newUserDto.AppAdmin,
		Avatar:   newUserDto.Avatar,
		Language: newUserDto.Language,
		AppTheme: newUserDto.AppTheme,
	}
}

// ToModelFromUserDto maps a UserDto to a User model (omits password field).
func FromUserDtoToUser(userDto *dto.UserDto) *model.User {
	return &model.User{
		ID:        userDto.ID,
		Username:  userDto.Username,
		AppAdmin:  userDto.AppAdmin,
		CreatedAt: userDto.CreatedAt,
		Avatar:    userDto.Avatar,
		Language:  userDto.Language,
		AppTheme:  userDto.AppTheme,
	}
}
