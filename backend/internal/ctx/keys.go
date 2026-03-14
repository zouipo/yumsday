package ctx

// Create unique types (no collision)
// https://pkg.go.dev/context#WithValue: "The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context."
type SessionCtxKey struct{}
type UserCtxKey struct{}
