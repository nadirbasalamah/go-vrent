package users

// WrongUsernameOrPasswordError represents username or password error
type WrongUsernameOrPasswordError struct{}

// Error returns error for invalid login
func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}
