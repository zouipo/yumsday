package mappers

import (
	"fmt"
	"testing"
	"time"

	"github.com/zouipo/yumsday/backend/internal/dtos"
	"github.com/zouipo/yumsday/backend/internal/models"
	"github.com/zouipo/yumsday/backend/internal/models/enums"
)

/*** DATA ***/

var creationTime = time.Now()

var user = models.User{
	ID:        1,
	Username:  "testuser",
	Password:  "securepassword",
	AppAdmin:  false,
	CreatedAt: creationTime,
	Avatar:    enums.Avatar1,
	Language:  enums.English,
	AppTheme:  enums.Light,
}

var userDtoNoPassword = dtos.UserDto{
	ID:        1,
	Username:  "testuser",
	AppAdmin:  false,
	CreatedAt: creationTime,
	Avatar:    enums.Avatar1,
	Language:  enums.English,
	AppTheme:  enums.Light,
}

var newUserDto = dtos.NewUserDto{
	Username: "testuser",
	Password: "securepassword",
	AppAdmin: false,
	Avatar:   enums.Avatar1,
	Language: enums.English,
	AppTheme: enums.Light,
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

func compareUserDtos(actual, expected *dtos.UserDto) (bool, string) {
	if actual.ID != expected.ID {
		return false, fmt.Errorf("ID mismatch: actual %d !=  expected %d", actual.ID, expected.ID).Error()
	}
	if actual.Username != expected.Username {
		return false, fmt.Errorf("Username mismatch: actual %s != expected %s", actual.Username, expected.Username).Error()
	}
	if actual.AppAdmin != expected.AppAdmin {
		return false, fmt.Errorf("AppAdmin mismatch: actual %v != expected %v", actual.AppAdmin, expected.AppAdmin).Error()
	}
	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		return false, fmt.Errorf("CreatedAt mismatch: actual %v != expected %v", actual.CreatedAt, expected.CreatedAt).Error()
	}
	if actual.Avatar != expected.Avatar {
		return false, fmt.Errorf("Avatar mismatch: actual %v != expected %v", actual.Avatar, expected.Avatar).Error()
	}
	if actual.Language != expected.Language {
		return false, fmt.Errorf("Language mismatch: actual %v != expected %v", actual.Language, expected.Language).Error()
	}
	if actual.AppTheme != expected.AppTheme {
		return false, fmt.Errorf("AppTheme mismatch: actual %v != expected %v", actual.AppTheme, expected.AppTheme).Error()
	}
	return true, ""
}

func compareUserNoID(actual, expected *models.User) (bool, string) {
	if actual.Username != expected.Username {
		return false, fmt.Errorf("Username mismatch: actual %s != expected %s", actual.Username, expected.Username).Error()
	}
	if actual.Password != expected.Password {
		return false, fmt.Errorf("Password mismatch: actual %s != expected %s", actual.Password, expected.Password).Error()
	}
	if actual.AppAdmin != expected.AppAdmin {
		return false, fmt.Errorf("AppAdmin mismatch: actual %v != expected %v", actual.AppAdmin, expected.AppAdmin).Error()
	}
	if !actual.CreatedAt.Equal(*new(time.Time)) {
		return false, fmt.Errorf("CreatedAt mismatch: actual %v != expected %v", actual.CreatedAt, expected.CreatedAt).Error()
	}
	if actual.Avatar != expected.Avatar {
		return false, fmt.Errorf("Avatar mismatch: actual %v != expected %v", actual.Avatar, expected.Avatar).Error()
	}
	if actual.Language != expected.Language {
		return false, fmt.Errorf("Language mismatch: actual %v != expected %v", actual.Language, expected.Language).Error()
	}
	if actual.AppTheme != expected.AppTheme {
		return false, fmt.Errorf("AppTheme mismatch: actual %v != expected %v", actual.AppTheme, expected.AppTheme).Error()
	}
	return true, ""
}

func compareUserNoPassword(actual, expected *models.User) (bool, string) {
	if actual.ID != expected.ID {
		return false, fmt.Errorf("ID mismatch: actual %d !=  expected %d", actual.ID, expected.ID).Error()
	}
	if actual.Username != expected.Username {
		return false, fmt.Errorf("Username mismatch: actual %s != expected %s", actual.Username, expected.Username).Error()
	}
	if actual.AppAdmin != expected.AppAdmin {
		return false, fmt.Errorf("AppAdmin mismatch: actual %v != expected %v", actual.AppAdmin, expected.AppAdmin).Error()
	}
	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		return false, fmt.Errorf("CreatedAt mismatch: actual %v != expected %v", actual.CreatedAt, expected.CreatedAt).Error()
	}
	if actual.Avatar != expected.Avatar {
		return false, fmt.Errorf("Avatar mismatch: actual %v != expected %v", actual.Avatar, expected.Avatar).Error()
	}
	if actual.Language != expected.Language {
		return false, fmt.Errorf("Language mismatch: actual %v != expected %v", actual.Language, expected.Language).Error()
	}
	if actual.AppTheme != expected.AppTheme {
		return false, fmt.Errorf("AppTheme mismatch: actual %v != expected %v", actual.AppTheme, expected.AppTheme).Error()
	}
	return true, ""
}
