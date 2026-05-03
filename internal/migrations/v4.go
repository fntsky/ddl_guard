package migrations

import (
	"context"

	"xorm.io/xorm"
)

func updateDDLRemindFields(ctx context.Context, x *xorm.Engine) error {
	session := x.Context(ctx)

	if _, err := session.Exec("ALTER TABLE ddl DROP COLUMN IF EXISTS early_remind_time"); err != nil {
		return err
	}

	if _, err := session.Exec("ALTER TABLE ddl DROP COLUMN IF EXISTS remind_sent"); err != nil {
		return err
	}

	if _, err := session.Exec("ALTER TABLE ddl ADD COLUMN IF NOT EXISTS subject TEXT"); err != nil {
		return err
	}

	if _, err := session.Exec("ALTER TABLE ddl ADD COLUMN IF NOT EXISTS remind_24h BOOLEAN NOT NULL DEFAULT false"); err != nil {
		return err
	}

	if _, err := session.Exec("ALTER TABLE ddl ADD COLUMN IF NOT EXISTS remind_2h BOOLEAN NOT NULL DEFAULT false"); err != nil {
		return err
	}

	return nil
}
