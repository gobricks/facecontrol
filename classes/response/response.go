package response

// BaseResponse is a basic JSON response
type BaseResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// TokenSignResponse is a basic JSON response with JWT token
type TokenSignResponse struct {
	BaseResponse
	Token string `json:"token"`
}

// UserResponse is a basic response with user credentials
type UserResponse struct {
	BaseResponse
	User interface{} `json:"user"`
}
