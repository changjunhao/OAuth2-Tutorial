package handle

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 模拟授权码、令牌等数据存储
var (
	codeMap         = make(map[string]string)
	codeScopeMap    = make(map[string][]string)
	TokenMap        = make(map[string]string)
	TokenScopeMap   = make(map[string][]string)
	refreshTokenMap = make(map[string]string)
	appMap          = make(map[string]string)
	reqidMap        = make(map[string]string)
)

func init() {
	//模拟第三方软件注册之后的数据库存储
	appMap["app_id"] = "APPID_RABBIT"
	appMap["app_secret"] = "APPSECRET_RABBIT"
	appMap["redirect_uri"] = "http://localhost:8080/AppServlet"
	appMap["scope"] = "today history"
}

// OauthHandle /** 模拟【授权服务】
func OauthHandle(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		doGet(writer, request)
	} else {
		doPost(writer, request)
	}
}

func doGet(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	responseType := query.Get("response_type")
	redirectUri := query.Get("redirect_uri")
	appId := query.Get("app_id")
	scope := query.Get("scope")

	fmt.Println("8081 GET responseType: " + responseType)

	if appMap["app_id"] != appId {
		return
	}

	if appMap["redirect_uri"] != redirectUri {
		return
	}

	// 验证第三方软件请求的权限范围是否与当时注册的权限范围一致
	if !checkScope(scope) {
		//超出注册的权限范围
		return
	}

	//生成页面reqid
	reqid := strconv.FormatInt(time.Now().UnixMilli(), 10)
	reqidMap[reqid] = reqid //保存该reqid值

	request.Header.Set("reqid", reqid)
	request.Header.Set("response_type", responseType)
	request.Header.Set("redirect_uri", redirectUri)
	request.Header.Set("app_id", appId)

	// 跳转到授权页面
	u, _ := url.Parse("http://localhost:8081/approve.html")
	proxy := httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.URL = u
		},
	}

	proxy.ServeHTTP(writer, request) // 授权码流程的【第一次】重定向

	//至此颁发授权码code的准备工作完毕
}

func doPost(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("start accept post req, generate access_toen")
	reqType := request.PostFormValue("reqType")

	grantType := request.PostFormValue("grant_type")
	appId := request.PostFormValue("app_id")
	appSecret := request.PostFormValue("app_secret")

	responseType := request.PostFormValue("response_type")
	redirectUri := request.PostFormValue("redirect_uri")
	//scope := request.PostFormValue("scope")

	//处理用户点击approve按钮动作
	if reqType == "approve" {
		reqid := request.PostFormValue("reqid") //假设一定能够获取到值

		if _, ok := reqidMap[reqid]; !ok {
			return
		}

		if responseType == "code" {
			rscope := request.Form["rscope"]
			if !checkRScope(rscope) { //验证权限范围，对又验证一次
				//超出注册的权限范围
				fmt.Println("out of scope ...")
				return
			}

			code := generateCode(appId, "USERTEST") //模拟登陆用户为USERTEST

			codeScopeMap[code] = rscope //授权范围与授权码做绑定

			oauthUrl, _ := url.Parse(redirectUri)
			params := oauthUrl.Query()
			params.Add("code", code)
			oauthUrl.RawQuery = params.Encode()                // 构造第三方软件的回调地址，并重定向到该地址
			writer.Header().Set("Location", oauthUrl.String()) // 授权码流程的【第二次】重定向
			writer.WriteHeader(302)
		}
	}

	//处理授权码流程中的 颁发访问令牌 环节
	if grantType == "authorization_code" {
		if appId != appMap["app_id"] {
			io.WriteString(writer, "app_id is not available")
			return
		}

		if appSecret != appMap["app_secret"] {
			io.WriteString(writer, "app_secret is not available")
			return
		}

		code := request.PostFormValue("code")

		fmt.Println("code", code)
		if !isExistCode(code) { //验证code值
			return
		}
		delete(codeMap, code) //授权码一旦被使用，须要立即作废

		fmt.Println("start generate access_toen")
		accessToken := generateAccessToken(appId, "USERTEST") //生成访问令牌access_token的值
		TokenScopeMap[accessToken] = codeScopeMap[code]       //授权范围与访问令牌绑定

		refreshToken := generateRefreshToken(appId, "USERTEST") //生成刷新令牌refresh_token的值

		// TODO: 2020/2/28 将accessToken和refreshToken做绑定 ，将refreshToken和codeScopeMap做绑定
		io.WriteString(writer, accessToken+"|"+refreshToken)
	} else if grantType == "refresh_token" {
		if appId != "APPIDTEST" {
			io.WriteString(writer, "app_id is not available")
			return
		}

		if appSecret != "APPSECRETTEST" {
			io.WriteString(writer, "app_secret is not available")
			return
		}

		refreshToken := request.PostFormValue("refresh_token")

		if _, ok := refreshTokenMap[refreshToken]; !ok { //该refresh_token值不存在
			return
		}

		appStr := refreshTokenMap["refresh_token"]
		if !strings.HasPrefix(appStr, appId+"|"+"USERTEST") { //该refresh_token值 不是颁发给该 第三方软件的
			return
		}

		accessToken := generateAccessToken(appId, "USERTEST") //生成访问令牌access_token的值

		// TODO: 2020/2/28 删除旧的access_token 、删除旧的refresh_token、生成新的refresh_token

		io.WriteString(writer, accessToken)
	}
}

func checkRScope(rscope []string) bool {
	scope := ""

	for i := 0; i < len(rscope); i++ {
		scope = scope + rscope[i]
	}

	return strings.Contains(strings.Replace(appMap["scope"], " ", "", -1), scope) //简单模拟权限验证
}

func generateRefreshToken(appId string, user string) string {
	uid, _ := uuid.NewUUID()
	refreshToken := uid.String()

	refreshTokenMap[refreshToken] = appId + "|" + user + "|" + strconv.FormatInt(time.Now().UnixMilli(), 10) //在这一篇章我们仅作为演示用，实际这应该是一个全局数据库,并且有有效期

	return refreshToken
}

func generateAccessToken(appId string, user string) string {
	uid, _ := uuid.NewUUID()

	accessToken := uid.String()

	expiresIn := "1" //1天时间过期

	TokenMap[accessToken] = appId + "|" + user + "|" + strconv.FormatInt(time.Now().UnixMilli(), 10) + "|" + expiresIn //在这一篇章我们仅作为演示用，实际这应该是一个全局数据库,并且有有效期

	return accessToken
}

func isExistCode(code string) bool {
	_, ok := codeMap[code]
	return ok
}

func checkScope(scope string) bool {
	fmt.Println("appMap size: " + string(rune(len(appMap))))
	fmt.Println("appMap scope: " + appMap["scope"])
	fmt.Println("scope: " + scope)

	return strings.Contains(appMap["scope"], scope) // 简单模拟权限验证
}

func generateCode(appId string, user string) string {
	strb := ""
	for i := 0; i < 8; i++ {
		strb = strb + strconv.Itoa(rand.Intn(8))
	}

	code := strb

	// 在这一篇章我们仅作为演示用，实际这应该是一个全局内存数据库，有效期官方建议是10分钟
	codeMap[code] = appId + "|" + user + "|" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	return code
}
