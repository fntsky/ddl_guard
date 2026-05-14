package migrations

import (
	"context"

	"github.com/fntsky/ddl_guard/internal/entity"
	"xorm.io/xorm"
)

func restructureGradeTables(ctx context.Context, x *xorm.Engine) error {
	session := x.Context(ctx)

	// Drop old tables
	if _, err := session.Exec("DROP TABLE IF EXISTS daily_scores"); err != nil {
		return err
	}
	if _, err := session.Exec("DROP TABLE IF EXISTS final_grades"); err != nil {
		return err
	}

	// Create new tables with updated schema
	if err := session.Sync2(&entity.FinalGrade{}, &entity.QuizScore{}, &entity.HomeworkScore{}); err != nil {
		return err
	}

	return nil
}