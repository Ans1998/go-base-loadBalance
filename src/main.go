package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"proxy/util"
)

type ProxyHandler struct {}
func (* ProxyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	if request.URL.Path=="/favicon.ico" { // 谷歌会访问一个图标文件，我们不做处理
		return
	}
	url, _:=url.Parse(util.LoadBalanceInit.RoundRobinByWeight().Host)
	proxy:=httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(writer, request)
}

func main() {
	http.ListenAndServe(":8080", &ProxyHandler{})
}
