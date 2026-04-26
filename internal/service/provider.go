package service

import (
	"github.com/fntsky/ddl_guard/internal/base/OTP"
	baseauth "github.com/fntsky/ddl_guard/internal/base/auth"
	ai "github.com/fntsky/ddl_guard/internal/service/ai"
	authsvc "github.com/fntsky/ddl_guard/internal/service/auth"
	dailyscore "github.com/fntsky/ddl_guard/internal/service/daily_score"
	"github.com/fntsky/ddl_guard/internal/service/ddl"
	finalgrade "github.com/fntsky/ddl_guard/internal/service/final_grade"
	"github.com/fntsky/ddl_guard/internal/service/exam"
	"github.com/fntsky/ddl_guard/internal/service/user"
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
	finalgrade.NewFinalGradeService,
	dailyscore.NewDailyScoreService,
)
