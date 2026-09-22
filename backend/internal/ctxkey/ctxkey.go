package ctxkey

import (
	"fmt"
	"strings"
)

// Types used to key value in contexts without collisions
// https://pkg.go.dev/context#WithValue: "The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context."

const (
	categoryStr = "category"
	codeStr     = "code"
	idStr       = "id"
	nameStr     = "name"
	priceStr    = "price"
	sessionStr  = "session"
	statusStr   = "status"
	userStr     = "user"
	usernameStr = "username"
	weightStr   = "weight"
)

// Same as fmt.Stringer
type Key interface {
	fmt.Stringer
}

type Category struct{}
type Code struct{}
type Id struct{}
type Name struct{}
type Price struct{}
type Session struct{}
type Status struct{}
type User struct{}
type Username struct{}
type Weight struct{}

// Implementation of Stringer interface for each key type

func (k Category) String() string { return categoryStr }
func (k Code) String() string     { return codeStr }
func (k Id) String() string       { return idStr }
func (k Name) String() string     { return nameStr }
func (k Price) String() string    { return priceStr }
func (k Session) String() string  { return sessionStr }
func (k Status) String() string   { return statusStr }
func (k User) String() string     { return userStr }
func (k Username) String() string { return usernameStr }
func (k Weight) String() string   { return weightStr }

func FromString(str string) Key {
	switch strings.ToLower(str) {

	// "standard" key names from const block
	case categoryStr:
		return &Category{}
	case codeStr:
		return &Code{}
	case idStr:
		return &Id{}
	case nameStr:
		return &Name{}
	case priceStr:
		return &Price{}
	case sessionStr:
		return &Session{}
	case statusStr:
		return &Status{}
	case userStr:
		return &User{}
	case usernameStr:
		return &Username{}
	case weightStr:
		return &Weight{}
	default:
		panic(fmt.Errorf("unhandled context key: %v", str))
	}
}
