package service

import (
	"github.com/fntsky/ddl_guard/internal/base/OTP"
	baseauth "github.com/fntsky/ddl_guard/internal/base/auth"
	ai "github.com/fntsky/ddl_guard/internal/service/ai"
	authsvc "github.com/fntsky/ddl_guard/internal/service/auth"
	"github.com/fntsky/ddl_guard/internal/service/ddl"
	exam "github.com/fntsky/ddl_guard/internal/service/exam"
	finalgrade "github.com/fntsky/ddl_guard/internal/service/final_grade"
	homeworkscore "github.com/fntsky/ddl_guard/internal/service/homework_score"
	quizscore "github.com/fntsky/ddl_guard/internal/service/quiz_score"
	"github.com/fntsky/ddl_guard/internal/service/user"
	"github.com/fntsky/ddl_guard/internal/service/wechat"
	"github.com/google/wire"
)

var ProviderSetService = wire.NewSet(
	baseauth.NewTokenService,
	otp.NewSMTPEmailOTP,
	authsvc.NewAuthService,
	ai.NewAIProvider,
	ddl.NewDDLService,
	exam.NewExamService,
	user.NewUserService,
	wechat.NewWechatService,
	finalgrade.NewFinalGradeService,
	quizscore.NewQuizScoreService,
	homeworkscore.NewHomeworkScoreService,
)