package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

type OIDCData struct {
	Sub string
}

func OIDCHtml(writer http.ResponseWriter, request *http.Request) {
	// 解析指定文件生成模板对象
	tem, err := template.ParseFiles("html/oidc.html")
	if err != nil {
		fmt.Println("读取文件失败,err", err)
		return
	}
	// 利用给定数据渲染模板，并将结果写入w
	data := OIDCData{
		Sub: request.Header.Get("sub"),
	}
	tem.Execute(writer, data)
}
