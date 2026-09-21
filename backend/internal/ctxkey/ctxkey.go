package ctxkey

import (
	"fmt"
	"strings"
)

// Types used to key value in contexts without collisions
// https://pkg.go.dev/context#WithValue: "The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context."

const (
	categoryKeyStr = "category"
	codeKeyStr     = "code"
	idKeyStr       = "id"
	nameKeyStr     = "name"
	priceKeyStr    = "price"
	sessionKeyStr  = "session"
	statusKeyStr   = "status"
	userKeyStr     = "user"
	usernameKeyStr = "username"
	weightKeyStr   = "weight"
)

type CtxKey interface {
	fmt.Stringer
}

type CategoryCtxKey struct{}
type CodeCtxKey struct{}
type IdCtxKey struct{}
type NameCtxKey struct{}
type PriceCtxKey struct{}
type SessionCtxKey struct{}
type StatusCtxKey struct{}
type UserCtxKey struct{}
type UsernameCtxKey struct{}
type WeightCtxKey struct{}

// Implementation of Stringer interface for each key type

func (k CategoryCtxKey) String() string { return categoryKeyStr }
func (k CodeCtxKey) String() string     { return codeKeyStr }
func (k IdCtxKey) String() string       { return idKeyStr }
func (k NameCtxKey) String() string     { return nameKeyStr }
func (k PriceCtxKey) String() string    { return priceKeyStr }
func (k SessionCtxKey) String() string  { return sessionKeyStr }
func (k StatusCtxKey) String() string   { return statusKeyStr }
func (k UserCtxKey) String() string     { return userKeyStr }
func (k UsernameCtxKey) String() string { return usernameKeyStr }
func (k WeightCtxKey) String() string   { return weightKeyStr }

func FromString(str string) CtxKey {
	switch strings.ToLower(str) {

	// "standard" key names from const block
	case categoryKeyStr:
		return &CategoryCtxKey{}
	case codeKeyStr:
		return &CodeCtxKey{}
	case idKeyStr:
		return &IdCtxKey{}
	case nameKeyStr:
		return &NameCtxKey{}
	case priceKeyStr:
		return &PriceCtxKey{}
	case sessionKeyStr:
		return &SessionCtxKey{}
	case statusKeyStr:
		return &StatusCtxKey{}
	case userKeyStr:
		return &UserCtxKey{}
	case usernameKeyStr:
		return &UsernameCtxKey{}
	case weightKeyStr:
		return &WeightCtxKey{}

	// additional key that are not covered by Stringer implementations
	case "userid":
		return &IdCtxKey{}

	default:
		panic(fmt.Errorf("unhandled context key: %v", str))
	}
}
