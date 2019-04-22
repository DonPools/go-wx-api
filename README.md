# 微信公众号开发SDK

go-wx-api是对微信公众号API的封装，可以当作SDK使用。可以同时启用多个公众号。

## 编译例子
 1. 该函数包已经使用go modules发布，需要golang 1.11.x及以上版本
 1. 请参考[go-wx-apps](https://github.com/rosbit/go-wx-apps)，那里包含了例程和工具程序

## 使用方法1: (单一公众号)

以下是一个简单的例子，用于说明使用go-wx-api的主要执行步骤。更详细的例子参考[go-wx-apps](https://github.com/rosbit/go-wx-apps)

```go
package main

import (
	"github.com/rosbit/go-wx-api/conf"
	"github.com/rosbit/go-wx-api"
	"net/http"
	"fmt"
)

const (
	token     = "微信公众号的token"
	appId     = "微信公众号appId"
	appSecret = "微信公众号的secret"
	aesKey    = "" //安全模式 使用的AESKey，如果是 明文传输，该串为空
	
	listenPort = 7070   // 服务侦听的端口号，请根据微信公众号管理端的服务器配置正确设置
	service    = "/wx"  // 微信公众号管理端服务器配置中URL的路径部分

	workerNum = 3 // 处理请求的并发数
)

func main() {
	// 步骤1. 设置配置参数
	if err := wxconf.SetParams(token, appId, appSecret, aesKey); err != nil {
		fmt.Printf("failed to set params: %v\n", err)
		return
	}

	// 步骤2. 初始化SDK
	wxapi.InitWxAPI(workerNum, os.Stdout)

	// 步骤2.5 设置签名验证的中间件。由于net/http不支持中间件，省去该步骤
	// signatureChecker := wxapi.NewWxSignatureChecker(wxconf.WxParams.Token, 0, []string{service})
	// <middleWareContainer>.Use(signatureChecker)

	// 步骤3. 设置http路由，启动http服务
	http.HandleFunc(service, wxapi.Echo)     // 用于配置
	http.HandleFunc(service, wxapi.Request)  // 用于实际执行公众号请求，和wxapi.Echo只能使用一个。
	                                         // 可以使用支持高级路由功能的web框架同时设置，参考 github.com/rosbit/go-wx-api/samples/wx-echo-server
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}
```

## 使用方法: (多个公众号)

以下代码仅仅为同时启用公众号的示例:

```go
package main

import (
	"github.com/rosbit/go-wx-api/conf"
	"github.com/rosbit/go-wx-api"
	"net/http"
	"fmt"
)

type WxConf struct {
	token string
	appId string
	appSecret string
	aesKey string
	workerNum int
	service string
}

var (
	listenPort = 7070   // 服务侦听的端口号，请根据微信公众号管理端的服务器配置正确设置
	wxServices = []WxConf{
		WxConf{
			token: "微信公众号1的token",
			appId: "微信公众号1的appId",
			appSecret: "微信公众号的1secret",
			aesKey: "",      // 安全模式 使用的AESKey，如果是 明文传输，该串为空
			workerNum: 3,    // 处理请求的并发数
			service: "/wx1", // 微信公众号管理端服务器配置中URL的路径部分
		},
		WxConf{
			token: "微信公众号2的token",
			appId: "微信公众号2的appId",
			appSecret: "微信公众号2的secret",
			aesKey: "",      // 安全模式 使用的AESKey，如果是 明文传输，该串为空
			workerNum: 3,    // 处理请求的并发数
			service: "/wx2", // 微信公众号管理端服务器配置中URL的路径部分
		},
		// 其它服务号
	}
)

func main() {
	// 对于每一个公众号执行
	for _, conf := range wxServices {
		// 步骤1. 设置配置参数
		wxParams, err := wxconf.NewWxParams(conf.token, conf.appId, conf.appSecret, conf.aesKey)
		if err != nil {
			fmt.Printf("failed to set params: %v\n", err)
			return
		}

		// 步骤2. 初始化SDK
		wxService := wxapi.InitWxAPIWithParams(wxParams, conf.workerNum, os.Stdout)

		// 步骤2.5 设置签名验证的中间件。由于net/http不支持中间件，省去该步骤
		// signatureChecker := wxapi.NewWxSignatureChecker(wxParams.Token, 0, []string{conf.service})
		// <middleWareContainer>.Use(signatureChecker)

		// 步骤3. 设置http路由，启动http服务
		http.HandleFunc(conf.service, wxService.Echo)     // 用于配置
		http.HandleFunc(conf.service, wxService.Request)  // 用于实际执行公众号请求，和wxService.Echo只能使用一个。
		                                                  // 可以使用支持高级路由功能的web框架同时设置
	}

	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}

```

## 其它
 1. 该函数包可以处理文本消息、用户关注/取消关注事件、菜单点击事件
 2. 其它消息、事件可以根据需要扩充
