package repo

import (
	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/repo/ddl"
	"github.com/fntsky/ddl_guard/internal/repo/exam"
	"github.com/fntsky/ddl_guard/internal/repo/final_grade"
	"github.com/fntsky/ddl_guard/internal/repo/homework_score"
	"github.com/fntsky/ddl_guard/internal/repo/quiz_score"
	"github.com/fntsky/ddl_guard/internal/repo/session"
	"github.com/fntsky/ddl_guard/internal/repo/user"
	"github.com/google/wire"
)

var ProviderSetRepo = wire.NewSet(
	data.NewDB,
	data.NewRedisClient,
	data.NewData,
	ddl.NewDDLRepo,
	exam.NewExamRepo,
	session.NewSessionRepo,
	user.NewUserRepo,
	final_grade.NewFinalGradeRepo,
	quiz_score.NewQuizScoreRepo,
	homework_score.NewHomeworkScoreRepo,
)