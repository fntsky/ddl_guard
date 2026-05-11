package migrations

import (
	"context"

	"xorm.io/xorm"
)

func addPhoneColumn(ctx context.Context, x *xorm.Engine) error {
	session := x.Context(ctx)

	if _, err := session.Exec("ALTER TABLE \"user\" ADD COLUMN IF NOT EXISTS phone VARCHAR(20) NULL UNIQUE"); err != nil {
		return err
	}

	return nil
}
