package wxapi

import (
	"io"
	"github.com/donpools/go-wx-api/auth"
	"github.com/donpools/go-wx-api/msg"
	"github.com/donpools/go-wx-api/log"
	"github.com/donpools/go-wx-api/conf"
)

type WxHandler struct {
	appIdHandler *wxauth.WxAppIdAuthHandler
	appMsgParser *wxmsg.WxAppIdMsgParser
}

var (
	defaultWxHandler *WxHandler
)

func InitWxAPI(workerNum int, logger io.Writer) {
	defaultWxHandler = InitWxAPIWithParams(nil, workerNum, logger)
}

func InitWxAPIWithParams(params *wxconf.WxParamsT, workerNum int, logger io.Writer) *WxHandler {
	wxlog.SetLogger(logger)
	appIdHandler := wxauth.StartAuthThreads(params, workerNum)
	appIdMsgParser := wxmsg.StartWxMsgParsers(params, workerNum)
	return &WxHandler{appIdHandler, appIdMsgParser}
}
