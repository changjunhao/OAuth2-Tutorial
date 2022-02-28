package handle

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ProtectedHandle /** 模拟【受保护资源服务】
func ProtectedHandle(writer http.ResponseWriter, request *http.Request) {
	//省略验证代码

	accessToken := request.PostFormValue("token")

	fmt.Println("accessToken", accessToken)

	//根据当时授权的token对应的权限范围，做相应的处理动作
	//不同权限对应不同的操作
	scope := TokenScopeMap[accessToken]

	var sbuf []rune
	for i := 0; i < len(scope); i++ {
		sbuf = append(append(sbuf, []rune(scope[i])...), []rune("|")...)
	}

	if strings.Index(string(sbuf), "query") > 0 {
		queryGoods("")
	}

	if strings.Index(string(sbuf), "add") > 0 {
		addGoods("")
	}

	if strings.Index(string(sbuf), "del") > 0 {
		delGoods("")
	}

	//不同的用户对应不同的数据
	user := TokenMap[accessToken]
	fmt.Println("user", user)
	io.WriteString(writer, queryOrders(user))
}

func queryGoods(id string) string {
	return ""
}

func addGoods(goods string) bool {
	return true
}

func delGoods(id string) bool {
	return true
}

func queryOrders(user string) string {
	return ""
}
