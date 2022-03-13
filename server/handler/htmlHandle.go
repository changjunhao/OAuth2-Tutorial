package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

type ApproveData struct {
	Regid        string
	ResponseType string
	RedirectUri  string
	AppId        string
}

func ApproveHtml(writer http.ResponseWriter, request *http.Request) {
	// 解析指定文件生成模板对象
	tem, err := template.ParseFiles("html/approve.html")
	if err != nil {
		fmt.Println("读取文件失败,err", err)
		return
	}
	// 利用给定数据渲染模板，并将结果写入w
	data := ApproveData{
		Regid:        request.Header.Get("reqid"),
		RedirectUri:  request.Header.Get("redirect_uri"),
		ResponseType: request.Header.Get("response_type"),
		AppId:        request.Header.Get("app_id"),
	}
	tem.Execute(writer, data)
}
