package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const OauthURl = "http://localhost:8081/OauthServlet"
const ProtectedURl = "http://localhost:8081/ProtectedServlet"

// AppServletHandle /** 模拟【第三方软件的Server端】
func AppServletHandle(writer http.ResponseWriter, request *http.Request) {
	//授权码许可流程，DEMO CODE

	query := request.URL.Query()
	fmt.Println(query)
	code := query.Get("code")

	params := url.Values{}
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")
	params.Add("app_id", "APPID_RABBIT")
	params.Add("app_secret", "APPSECRET_RABBIT")

	fmt.Println("start post code for token ...")
	response, err := http.PostForm(OauthURl, params)
	if err != nil {
		return
	}
	defer response.Body.Close()
	accessToken, _ := io.ReadAll(response.Body)

	fmt.Println("accessToken:" + string(accessToken))

	//使用 accessToken 请求受保护资源服务

	paramsMap := url.Values{}

	paramsMap.Add("app_id", "APPID_RABBIT")
	paramsMap.Add("app_secret", "APPSECRET_RABBIT")
	paramsMap.Add("token", string(accessToken))

	resp, _ := http.PostForm(ProtectedURl, paramsMap)
	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)

	fmt.Println(string(result))
}

func AppServletPasswordHandle(writer http.ResponseWriter, request *http.Request) {
	// 资源拥有者凭据许可流程，DEMO CODE
	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("app_id", "APPIDTEST")
	params.Add("app_secret", "APPSECRETTEST")
	params.Add("username", "USERNAMETEST")
	params.Add("password", "PASSWORDTEST")

	fmt.Println("start post code for token ...")
	response, err := http.PostForm(OauthURl, params)
	if err != nil {
		return
	}
	defer response.Body.Close()
	accessToken, _ := io.ReadAll(response.Body)

	fmt.Println("accessToken:" + string(accessToken))
}

// AppServletClientCredentialsHandle 第三方软件凭据许可流程，DEMO CODE
func AppServletClientCredentialsHandle(writer http.ResponseWriter, request *http.Request) {
	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("app_id", "APPIDTEST")
	params.Add("app_secret", "APPSECRETTEST")
	fmt.Println("start post code for token ...")
	response, err := http.PostForm(OauthURl, params)
	if err != nil {
		return
	}
	defer response.Body.Close()
	accessToken, _ := io.ReadAll(response.Body)

	fmt.Println("accessToken:" + string(accessToken))
}

// AppServletTokenHandle 隐式许可流程（模拟），DEMO CODE
func AppServletTokenHandle(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	accessToken := query.Get("access_token")

	fmt.Println("accessToken:" + string(accessToken))
}
