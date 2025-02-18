package usrdatasvc

import (
	"context"
	"fmt"
	"github.com/kv1sidisi/shapeup/pkg/errdefs"
	"github.com/kv1sidisi/shapeup/pkg/proto/usrdatasvc/pb"
	"github.com/kv1sidisi/shapeup/services/usrdatasvc/internal/domain/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log/slog"
	"reflect"
	"strings"
)

type UsrDataSvc struct {
	log        *slog.Logger
	usrDataMgr UsrDataMgr
}

type UsrDataMgr interface {
	UpdBsUsrAttr(ctx context.Context, bsusrattr *models.BsUsrAttr) error
	GetById(ctx context.Context) (*models.BsUsrAttr, error)
}

func New(log *slog.Logger, usrDataMgr UsrDataMgr) *UsrDataSvc {
	return &UsrDataSvc{log, usrDataMgr}
}

// UpdUsr applies field mask to base user's data and updates base user data through database manager.
//
// Returns:
//
//   - Proto of base user's data if successful.
//
//   - Error if: database manager returns error. Applying field mask fails.
func (u *UsrDataSvc) UpdUsr(ctx context.Context, bsusrattr *usrdatasvc.BsUsrAttr, mask *fieldmaskpb.FieldMask) (updbsusrattr *usrdatasvc.BsUsrAttr, err error) {
	const op = "usrdatasvc.UpdUsr"

	log := u.log.With(slog.String("op", op))

	existingUser, err := u.usrDataMgr.GetById(ctx)
	if err != nil {
		return nil, err
	}

	updatedUser := protoToDomain(bsusrattr)

	if err := applyFieldMask(u.log, existingUser, updatedUser, mask); err != nil {
		log.Error("failed to apply mask")
		return nil, err
	}

	if err := u.usrDataMgr.UpdBsUsrAttr(ctx, updatedUser); err != nil {
		return nil, err
	}

	return domainToProto(updatedUser), nil
}

// applyFieldMask applies mask to data.
//
// Returns nil if successful.
//
// Return error if some field is invalid.
func applyFieldMask(logger *slog.Logger, existingPtr interface{}, updatePtr interface{}, mask *fieldmaskpb.FieldMask) error {
	const op = "usrdatasvc.applyFieldMask"
	log := logger.With(slog.String("op", op))

	existingVal := reflect.ValueOf(existingPtr).Elem()
	updateVal := reflect.ValueOf(updatePtr).Elem()

	for _, path := range mask.Paths {
		fieldName := snakeToCamel(path)
		field := existingVal.FieldByName(fieldName)
		if !field.IsValid() {
			log.Error(fmt.Sprintf("field %s not found", fieldName))
			return errdefs.ErrFieldMask
		}
		updateField := updateVal.FieldByName(fieldName)
		if !updateField.IsValid() {
			log.Error(fmt.Sprintf("field %s not found", fieldName))
			return errdefs.ErrFieldMask
		}
		if field.CanSet() {
			field.Set(updateVal)
		} else {
			log.Error(fmt.Sprintf("field %s can not be set", fieldName))
			return errdefs.ErrFieldMask
		}
	}
	return nil
}

func snakeToCamel(name string) string {
	parts := strings.Split(name, "_")
	caser := cases.Title(language.Und)

	for i, part := range parts {
		parts[i] = caser.String(part)
	}
	return strings.Join(parts, "")
}

func protoToDomain(pu *usrdatasvc.BsUsrAttr) *models.BsUsrAttr {
	return &models.BsUsrAttr{
		Name:      pu.GetName(),
		Height:    pu.GetHeight(),
		Weight:    pu.GetWeight(),
		Gender:    pu.GetGender(),
		BirthDate: pu.GetBirthDate(),
	}
}

func domainToProto(d *models.BsUsrAttr) *usrdatasvc.BsUsrAttr {
	return &usrdatasvc.BsUsrAttr{
		Name:      d.Name,
		Height:    d.Height,
		Weight:    d.Weight,
		Gender:    d.Gender,
		BirthDate: d.BirthDate,
	}
}
