package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"sample/sample/manager"
	"sample/source"
	"sample/source/response"
	"sample/source/tool"
)

var svr *Service

type Service struct {
	logger   *logrus.Logger
	tool     *tool.Tool
	response *response.ResponseManager
}

func init() {
	svr = new(Service)
	svr.logger = source.GetLogger()
	svr.tool = tool.GetTool()
}

func (s *Service) Test(ctx *gin.Context) {
	manager.GetManager().Sample(ctx)
	s.response.ReturnSucceedResponseWithoutData(ctx)
}

//func (s *Service) Login(ctx *gin.Context) {
//	registerInformation := new(model.RegisterInformation)
//	err := ctx.Bind(registerInformation)
//	if err != nil {
//		s.logger.Errorln(err)
//		return
//	}
//	token, err := user.man.Login(ctx, registerInformation.Name, registerInformation.Email, registerInformation.VerificationCode)
//	if err != nil {
//		//s.logger.Errorln(registerInformation.Email, err)
//		user.responseManager.ReturnFailedResponse(ctx, s.tool.SprintfErr(err))
//		return
//	}
//	user.responseManager.ReturnSucceedResponse(ctx, map[string]string{"token": token})
//}

//func (s *Service) Login(ctx *gin.Context) {
//	token, err := man.Login(ctx.Query("phone"), ctx.Query("vCode"))
//	if err != nil {
//		s.res.NewResWithoutData(ctx, code.Failed, err.Error())
//		return
//	}
//	s.res.NewSucceedRes(ctx, map[string]string{"token": token})
//}

//func (s *Service) SendVerificationCode(ctx *gin.Context) {
//	email := ctx.Query("email")
//	err := user.man.SendVerificationCode(email)
//	if err != nil {
//		s.logger.Errorln(email, err)
//		user.responseManager.ReturnFailedResponse(ctx, s.tool.SprintfErr(err))
//		return
//	}
//	user.responseManager.ReturnSucceedResponseWithoutData(ctx, nil)
//}

//
//func (s *Service) GetUserInfo(ctx *gin.Context) {
//	phone := ctx.Query("phone")
//	id := ctx.Query("id")
//	var user *User
//	var err error
//	if len(phone) > 0 {
//		user, err = man.GetUserInfoByPhone(phone)
//	}
//	if len(id) > 0 {
//		user, err = man.GetUserInfoByID(id)
//	}
//	if err != nil {
//		s.res.NewRes(ctx, code.Failed, "get user info failed", user)
//		return
//	}
//	s.res.NewRes(ctx, 1, "ok", user)
//}

func GetService() *Service {
	return svr
}
