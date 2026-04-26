package migrations

import (
	"context"

	"github.com/fntsky/ddl_guard/internal/entity"
	"xorm.io/xorm"
)

func addExamTable(ctx context.Context, x *xorm.Engine) error {
	return x.Context(ctx).Sync2(&entity.Exam{})
}

func addGradeTables(ctx context.Context, x *xorm.Engine) error {
	return x.Context(ctx).Sync2(&entity.FinalGrade{}, &entity.DailyScore{})
}
