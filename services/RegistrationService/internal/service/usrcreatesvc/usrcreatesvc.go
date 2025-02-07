package usrcreatesvc

import (
	"context"
	"github.com/kv1sidisi/shapeup/libs/common/errdefs"
	pbjwtsvc "github.com/kv1sidisi/shapeup/services/regsvc/api/grpccl/pb/jwtsvc"
	pbsendsvc "github.com/kv1sidisi/shapeup/services/regsvc/api/grpccl/pb/sendsvc"
	"github.com/kv1sidisi/shapeup/services/regsvc/cmd/grpccl"
	"github.com/kv1sidisi/shapeup/services/regsvc/cmd/grpccl/consts"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

const (
	confirmAccountLinkBase    = "http://localhost:8082/confirm_account?token="
	confirmationOperationType = "confirmation"
)

// UsrCreateSvc implementation of user creation service.
type UsrCreateSvc struct {
	log           *slog.Logger
	userSaver     UsrMgr
	sendingClient pbsendsvc.SendingClient
	jwtClient     pbjwtsvc.JWTClient
}

// UsrMgr interface defines the method for saving user information in database.
type UsrMgr interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
	ConfirmAccount(
		ctx context.Context,
		uid int64,
	) (err error)
	DeleteUser(
		ctx context.Context,
		uid int64,
	) (err error)
}

// New returns a new instance of UserCreation service.
func New(log *slog.Logger,
	userSaver UsrMgr,
	grpccl *grpccl.GRPCClients,
) *UsrCreateSvc {
	return &UsrCreateSvc{
		userSaver:     userSaver,
		log:           log,
		sendingClient: grpccl.Cl[consts.SendSvc].Client.(pbsendsvc.SendingClient),
		jwtClient:     grpccl.Cl[consts.JWTSvc].Client.(pbjwtsvc.JWTClient),
	}
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (r *UsrCreateSvc) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "usrcreatesvc.RegisterNewUser"

	log := r.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	// Generate a hashed password from the provided password.
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generating password hash")
		return 0, errdefs.ErrGeneratingPassword
	}
	log.Info("password hash generated")

	// Saving user in database
	uid, err := r.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		return 0, err
	}
	log.Info("user successfully saved: ", slog.Int64("userId", uid))

	// Generating confirmation link
	linkGenResp, err := r.jwtClient.GenerateLink(ctx, &pbjwtsvc.GenerateLinkRequest{
		LinkBase:  confirmAccountLinkBase,
		Uid:       uid,
		Operation: confirmationOperationType,
	})
	if err != nil {
		log.Error("confirmation link generation failed", err)
		if err := r.userSaver.DeleteUser(ctx, uid); err != nil {
			return 0, err
		}
		log.Error("compensating move, user deleted")
		return 0, err
	}
	log.Info("confirmation link generated successfully: ", linkGenResp.GetLink())

	// TODO: connect sending service
	// Sending confirmation link
	//_, err = r.sendingClient.SendEmail(ctx, &pbsendsvc.EmailRequest{
	//	Message: linkGenResp.GetLink(),
	//	Email:   email,
	//})
	//if err != nil {
	//	log.Error("failed to send confirmation link", err)
	//	return 0, err
	//}
	//log.Info("user confirmation link sent successfully")

	return uid, nil
}

// ConfirmNewUser confirms account
// If user does not exist returns error
func (r *UsrCreateSvc) ConfirmNewUser(ctx context.Context, token string) (uid int64, err error) {
	const op = "register.ConfirmAccount"

	log := r.log.With(
		slog.String("op", op),
		slog.String("token", token),
	)

	validationResp, err := r.jwtClient.ValidateAccessToken(ctx, &pbjwtsvc.ValidateAccessTokenRequest{
		Token: token,
	})
	if err != nil {
		log.Error("failed to verify confirmation token")
		return -1, err
	}
	uid = validationResp.GetUid()
	log.Info("user confirmation token verified successfully: ", slog.Int64("userId", uid))

	// Confirming user through database
	if err := r.userSaver.ConfirmAccount(ctx, uid); err != nil {
		return -1, err
	}

	return uid, nil
}
