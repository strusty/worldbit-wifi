package admins

type Admins interface {
	Login(request LoginRequest) (*JWTResponse, error)
	ChangePassword(adminID string, request ChangePasswordRequest) error
}
