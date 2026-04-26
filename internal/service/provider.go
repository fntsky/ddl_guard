package service

import (
	"github.com/fntsky/ddl_guard/internal/base/OTP"
	baseauth "github.com/fntsky/ddl_guard/internal/base/auth"
	ai "github.com/fntsky/ddl_guard/internal/service/ai"
	authsvc "github.com/fntsky/ddl_guard/internal/service/auth"
	"github.com/fntsky/ddl_guard/internal/service/ddl"
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
)
