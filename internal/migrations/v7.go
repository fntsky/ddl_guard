package migrations

import (
	"context"

	"xorm.io/xorm"
)

func fixUserNullableFields(ctx context.Context, x *xorm.Engine) error {
	if _, err := x.Context(ctx).Exec("UPDATE \"user\" SET email = NULL WHERE email = ''"); err != nil {
		return err
	}
	if _, err := x.Context(ctx).Exec("UPDATE \"user\" SET phone = NULL WHERE phone = ''"); err != nil {
		return err
	}
	return nil
}
