package handler

import (
	"fmt"
	"net/http"
	"net/url"
)

/**
 * 模拟【第三方软件的首页】
 * 浏览器输入 http://localhost:8080/AppIndexServlet-ch03
 */

// 8080: 三方软件，
// 8081：授权服务，8081：受保护资源服务 为了演示方便将授权服务和受保护资源服务放在同一个服务上面

const OauthUrl = "http://localhost:8081/OauthServlet?reqType=oauth"
const RedirectUrl = "http://localhost:8080/AppServlet"

func AppIndexHandle(writer http.ResponseWriter, request *http.Request) {
	// 授权码许可流程，DEMO CODE
	fmt.Println("app index ...")
	oauthUrl, _ := url.Parse(OauthUrl)
	params := oauthUrl.Query()
	params.Add("response_type", "code")
	params.Add("redirect_uri", RedirectUrl)
	params.Add("app_id", "APPID_RABBIT")
	params.Add("scope", "today history")
	oauthUrl.RawQuery = params.Encode()                // 构造请求授权的URl
	writer.Header().Set("Location", oauthUrl.String()) // 授权码流程的【第一次】重定向
	writer.WriteHeader(302)
}
