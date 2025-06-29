package usrcreatesvc

import (
	"context"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	pbjwtsvc "github.com/kv1sidisi/shapeup/pkg/proto/jwtsvc/pb"
	pbsendsvc "github.com/kv1sidisi/shapeup/pkg/proto/sendsvc/pb"
	"github.com/kv1sidisi/shapeup/services/regsvc/cmd/grpccl"
	"github.com/kv1sidisi/shapeup/services/regsvc/cmd/grpccl/consts"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

const (
	confirmAccountLinkBase    = "http://localhost:8082/confirm_account?token="
	confirmationOperationType = "confirmation"
)

// UsrCreateSvc user creation service.
type UsrCreateSvc struct {
	log           *slog.Logger
	userSaver     UsrMgr
	sendingClient pbsendsvc.SendingClient
	jwtClient     pbjwtsvc.JWTClient
}

// UsrMgr manager for database.
type UsrMgr interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid []byte, err error)
	ConfirmAccount(
		ctx context.Context,
		uid []byte,
	) (err error)
	DeleteUser(
		ctx context.Context,
		uid []byte,
	) (err error)
}

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

// RegisterNewUser registers new user.
//
// Returns:
//
//   - user ID if operation successful.
//
//   - Error if: user with given username already exists.
//     Password hash generation fails.
//     Generation confirmation link fails.
func (r *UsrCreateSvc) RegisterNewUser(ctx context.Context, email, password string) ([]byte, error) {
	const op = "usrcreatesvc.RegisterNewUser"

	log := r.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	// Generate a hashed password from the provided password.
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generating password hash")
		return passHash, errdefs.ErrGeneratingPassword
	}
	log.Info("password hash generated")

	// Saving user in database
	uid, err := r.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		return uid, err
	}
	log.Info("user successfully saved: ", slog.Any("userId", uid))

	// Generating confirmation link
	linkGenResp, err := r.jwtClient.GenerateLink(ctx, &pbjwtsvc.GenerateLinkRequest{
		LinkBase:  confirmAccountLinkBase,
		Uid:       uid,
		Operation: confirmationOperationType,
	})
	if err != nil {
		log.Error("confirmation link generation failed", err)
		if err := r.userSaver.DeleteUser(ctx, uid); err != nil {
			return uid, err
		}
		log.Error("compensating move, user deleted")
		return uid, err
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

// ConfirmNewUser confirms account.
//
// Returns:
//   - user ID if operation successful.
//   - Error if: user does not exist. JWT token is invalid.
func (r *UsrCreateSvc) ConfirmNewUser(ctx context.Context, token string) (uid []byte, err error) {
	const op = "register.ConfirmAccount"

	log := r.log.With(
		slog.String("op", op),
		slog.String("token", token),
	)

	validationResp, err := r.jwtClient.ValidateToken(ctx, &pbjwtsvc.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		log.Error("failed to verify confirmation token")
		return uid, err
	}
	uid = validationResp.GetUid()
	log.Info("user confirmation token verified successfully: ", slog.Any("userId", uid))

	// Confirming user through database
	if err := r.userSaver.ConfirmAccount(ctx, uid); err != nil {
		return uid, err
	}

	return uid, nil
}
