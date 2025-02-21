package errdefs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserMetricsAlreadyExists = newError(codes.AlreadyExists, "user metrics already exists")
	ErrFieldMask                = newError(codes.Internal, "error applying field mask")
	ErrInternal                 = newError(codes.Internal, "internal error")
	ErrUserNotConfirmed         = newError(codes.Unauthenticated, "user not confirmed")
	ErrUserNotFound             = newError(codes.NotFound, "user not found")
	ErrSessionAlreadyExists     = newError(codes.AlreadyExists, "session already exists")
	ErrEmailAlreadyExists       = newError(codes.AlreadyExists, "email already exists")
	ErrTokenExpired             = newError(codes.InvalidArgument, "token expired")
	ErrGeneratingPassword       = newError(codes.Internal, "password hash generation failed")
	ErrSendEmail                = newError(codes.Internal, "failed to send email")
	ErrDatabaseInternal         = newError(codes.Internal, "internal database error")
)

var (
	// InvalidRequest used when grpc request data is invalid.
	InvalidRequest       = newError(codes.InvalidArgument, "invalid request")
	InvalidCredentials   = newError(codes.InvalidArgument, "invalid credentials")
	InvalidUserId        = newError(codes.InvalidArgument, "invalid user id")
	InvalidEmail         = newError(codes.InvalidArgument, "invalid email")
	InvalidOperationType = newError(codes.InvalidArgument, "invalid operation type")
	InvalidToken         = newError(codes.InvalidArgument, "invalid token")
	InvalidLinkBase      = newError(codes.InvalidArgument, "invalid link base")
	InvalidSigningMethod = newError(codes.InvalidArgument, "invalid token signing method")
)

var (
	RegistrationSuccess    = newError(codes.OK, "registration success")
	TokenGenerationSuccess = newError(codes.OK, "token generation success")
	TokenValidationSuccess = newError(codes.OK, "token validation success")
	LinkGenerationSuccess  = newError(codes.OK, "link generation success")
	ConfirmUserSuccess     = newError(codes.OK, "confirmation user success")
	SendEmailSuccess       = newError(codes.OK, "send email success")
	AuthenticationSuccess  = newError(codes.OK, "authentication success")
	UpdUsrMetricsSuccess   = newError(codes.OK, "update user metrics success")
)

// newError returns status error created from code and message.
func newError(code codes.Code, msg string) error {
	return status.Error(code, msg)
}
