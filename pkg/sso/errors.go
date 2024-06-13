package sso

import (
	"fmt"

	sssov1 "go.pilab.hu/pilab-cloud/ga-pi/gen/pilab/ssso/v1"
	"go.pilab.hu/pilab-cloud/ga-pi/pkg/sso/codes"
)

type AuthError struct {
	Code    codes.Code
	Message string
	Params  any
}

func newAuthError(code codes.Code, label string) *AuthError {
	return &AuthError{
		Code:    code,
		Message: label,
	}
}

// Error returns the error message.
func (e *AuthError) Error() string {
	if e.Params != nil {
		return fmt.Sprintf("auth error: %s code: %d params: %v", e.Message, e.Code, e.Params)
	}

	return fmt.Sprintf("auth error: %s code: %d", e.Message, e.Code)
}

// ToResponse returns the error as a LoginResponse.
func (e *AuthError) ToResponse() *sssov1.TokenResponse {
	return &sssov1.TokenResponse{
		Response: &sssov1.TokenResponse_Error{
			Error: &sssov1.ErrorMessage{
				Code:    (int32)(e.Code),
				Message: e.Message,
			},
		},
	}
}

// WithParams sets the params field of the error.
func (e *AuthError) WithParams(params any) *AuthError {
	e.Params = params
	return e
}

var (
	ErrNoTenantOrClient   = newAuthError(codes.NoTenantOrClient, "no_tenant_or_client")
	ErrUserNotFound       = newAuthError(codes.UserNotFound, "user_not_found")
	ErrInvalidCredentials = newAuthError(codes.InvalidCredentials, "invalid_credentials")
	ErrLocked             = newAuthError(codes.Locked, "account_locked")
	ErrPasswordExpired    = newAuthError(codes.PasswordExpired, "password_expired")
	// Deprecated: use ErrRequestMFA instead.
	ErrRequestMFA      = newAuthError(codes.RequestMFA, "request_mfa")
	ErrTooManyAttempts = newAuthError(codes.TooManyAttempts, "too_many_attempts")
	ErrAlreadyExists   = newAuthError(codes.AlreadyExists, "already_exists")
)
