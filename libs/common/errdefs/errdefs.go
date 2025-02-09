package errdefs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal             = newError(codes.Internal, "internal error")
	ErrUserNotConfirmed     = newError(codes.Unauthenticated, "user not confirmed")
	ErrUserNotFound         = newError(codes.NotFound, "user not found")
	ErrSessionAlreadyExists = newError(codes.AlreadyExists, "session already exists")
	ErrEmailAlreadyExists   = newError(codes.AlreadyExists, "email already exists")
	ErrTokenExpired         = newError(codes.InvalidArgument, "token expired")
	ErrGeneratingPassword   = newError(codes.Internal, "password hash generation failed")
	ErrSendEmail            = newError(codes.Internal, "failed to send email")
	ErrDatabaseInternal     = newError(codes.Internal, "internal database error")
)

var (
	InvalidCredentials   = newError(codes.InvalidArgument, "invalid credentials")
	InvalidUserId        = newError(codes.InvalidArgument, "invalid user id")
	InvalidEmail         = newError(codes.InvalidArgument, "invalid email")
	InvalidOperationType = newError(codes.InvalidArgument, "invalid operation type")
	InvalidToken         = newError(codes.InvalidArgument, "invalid token")
	InvalidLinkBase      = newError(codes.InvalidArgument, "invalid link base")
	InvalidSigningMethod = newError(codes.InvalidArgument, "invalid token signing method")
)

// newError returns status error created from code and message.
func newError(code codes.Code, msg string) error {
	return status.Error(code, msg)
}
