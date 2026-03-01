package mapper

import (
	"fmt"
	"testing"
	"time"

	"github.com/zouipo/yumsday/backend/internal/dto"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/model/enum"
)

/*** DATA ***/

var creationTime = time.Now()

var avatar1 = enum.Avatar1

var user = model.User{
	ID:        1,
	Username:  "testuser",
	Password:  "securepassword",
	AppAdmin:  false,
	CreatedAt: creationTime,
	Avatar:    &avatar1,
	Language:  enum.English,
	AppTheme:  enum.Light,
}

var userDtoNoPassword = dto.UserDto{
	ID:        1,
	Username:  "testuser",
	AppAdmin:  false,
	CreatedAt: creationTime,
	Avatar:    &avatar1,
	Language:  enum.English,
	AppTheme:  enum.Light,
}

var newUserDto = dto.NewUserDto{
	Username: "testuser",
	Password: "securepassword",
	AppAdmin: false,
	Avatar:   &avatar1,
	Language: enum.English,
	AppTheme: enum.Light,
}

/*** TESTS ***/

func TestToUserDtoNoPassword(t *testing.T) {
	mappedDto := ToUserDtoNoPassword(&user)

	isIdentical, errMsg := compareUserDtos(mappedDto, &userDtoNoPassword)
	if !isIdentical {
		t.Errorf("ToUserDtoNoPassword mapping failed. %s", errMsg)
	}
}

func TestFromNewUserDtoToUser(t *testing.T) {
	mappedUser := FromNewUserDtoToUser(&newUserDto)

	isIdentical, errMsg := compareUserNoID(mappedUser, &user)
	if !isIdentical {
		t.Errorf("FromNewUserDtoToUser mapping failed. %s", errMsg)
	}
}

func TestFromUserDtoToUser(t *testing.T) {
	mappedUser := FromUserDtoToUser(&userDtoNoPassword)

	isIdentical, errMsg := compareUserNoPassword(mappedUser, &user)
	if !isIdentical {
		t.Errorf("FromUserDtoToUser mapping failed. %s", errMsg)
	}
}

/*** HELPERS ***/

func compareUserDtos(actual, expected *dto.UserDto) (bool, error) {
	if actual.ID != expected.ID {
		return false, fmt.Errorf("ID mismatch: actual %d !=  expected %d", actual.ID, expected.ID)
	}
	if actual.Username != expected.Username {
		return false, fmt.Errorf("Username mismatch: actual %s != expected %s", actual.Username, expected.Username)
	}
	if actual.AppAdmin != expected.AppAdmin {
		return false, fmt.Errorf("AppAdmin mismatch: actual'%v'!= expected %v", actual.AppAdmin, expected.AppAdmin)
	}
	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		return false, fmt.Errorf("CreatedAt mismatch: actual'%v'!= expected %v", actual.CreatedAt, expected.CreatedAt)
	}
	if actual.Avatar != expected.Avatar {
		return false, fmt.Errorf("Avatar mismatch: actual'%v'!= expected %v", actual.Avatar, expected.Avatar)
	}
	if actual.Language != expected.Language {
		return false, fmt.Errorf("Language mismatch: actual'%v'!= expected %v", actual.Language, expected.Language)
	}
	if actual.AppTheme != expected.AppTheme {
		return false, fmt.Errorf("AppTheme mismatch: actual'%v'!= expected %v", actual.AppTheme, expected.AppTheme)
	}
	return true, nil
}

func compareUserNoID(actual, expected *model.User) (bool, error) {
	if actual.Username != expected.Username {
		return false, fmt.Errorf("Username mismatch: actual %s != expected %s", actual.Username, expected.Username)
	}
	if actual.Password != expected.Password {
		return false, fmt.Errorf("Password mismatch: actual %s != expected %s", actual.Password, expected.Password)
	}
	if actual.AppAdmin != expected.AppAdmin {
		return false, fmt.Errorf("AppAdmin mismatch: actual'%v'!= expected %v", actual.AppAdmin, expected.AppAdmin)
	}
	if !actual.CreatedAt.Equal(*new(time.Time)) {
		return false, fmt.Errorf("CreatedAt mismatch: actual'%v'!= expected %v", actual.CreatedAt, expected.CreatedAt)
	}
	if (actual.Avatar == nil && expected.Avatar != nil) || (actual.Avatar != nil && expected.Avatar == nil) {
		return false, fmt.Errorf("Avatar mismatch: actual'%v'!= expected %v", actual.Avatar, expected.Avatar)
	}
	if actual.Avatar != nil && expected.Avatar != nil && *actual.Avatar != *expected.Avatar {
		return false, fmt.Errorf("Avatar mismatch: actual'%v'!= expected %v", *actual.Avatar, *expected.Avatar)
	}
	if actual.Language != expected.Language {
		return false, fmt.Errorf("Language mismatch: actual'%v'!= expected %v", actual.Language, expected.Language)
	}
	if actual.AppTheme != expected.AppTheme {
		return false, fmt.Errorf("AppTheme mismatch: actual'%v'!= expected %v", actual.AppTheme, expected.AppTheme)
	}
	return true, nil
}

func compareUserNoPassword(actual, expected *model.User) (bool, error) {
	if actual.ID != expected.ID {
		return false, fmt.Errorf("ID mismatch: actual %d !=  expected %d", actual.ID, expected.ID)
	}
	if actual.Username != expected.Username {
		return false, fmt.Errorf("Username mismatch: actual %s != expected %s", actual.Username, expected.Username)
	}
	if actual.AppAdmin != expected.AppAdmin {
		return false, fmt.Errorf("AppAdmin mismatch: actual'%v'!= expected %v", actual.AppAdmin, expected.AppAdmin)
	}
	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		return false, fmt.Errorf("CreatedAt mismatch: actual'%v'!= expected %v", actual.CreatedAt, expected.CreatedAt)
	}
	if (actual.Avatar == nil && expected.Avatar != nil) || (actual.Avatar != nil && expected.Avatar == nil) {
		return false, fmt.Errorf("Avatar mismatch: actual'%v'!= expected %v", actual.Avatar, expected.Avatar)
	}
	if actual.Avatar != nil && expected.Avatar != nil && *actual.Avatar != *expected.Avatar {
		return false, fmt.Errorf("Avatar mismatch: actual'%v'!= expected %v", *actual.Avatar, *expected.Avatar)
	}
	if actual.Language != expected.Language {
		return false, fmt.Errorf("Language mismatch: actual'%v'!= expected %v", actual.Language, expected.Language)
	}
	if actual.AppTheme != expected.AppTheme {
		return false, fmt.Errorf("AppTheme mismatch: actual'%v'!= expected %v", actual.AppTheme, expected.AppTheme)
	}
	return true, nil
}
