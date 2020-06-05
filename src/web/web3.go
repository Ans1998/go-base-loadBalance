package main

import "net/http"

type webHandler struct {}
func (* webHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	writer.Write([]byte("<h1>web3</h1>"))
}
func main()  {
	http.ListenAndServe(":9093", &webHandler{})
}