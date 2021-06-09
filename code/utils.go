package account

import (
	"fmt"
	"net/http"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func SetCookie(data_info map[string]interface{}, response http.ResponseWriter, request *http.Request) {
	_, err := request.Cookie("mysession")
	if err == nil && data["userConnected"] != "" {
		data_info["cookieExist"] = true
		data_info["userConnected"] = data["userConnected"]
	} else {
		CheckError(err)
		data_info["cookieExist"] = false
		data_info["userConnected"] = ""
	}
}

func DeleteCookie() {
	data["userConnected"] = ""
	data["cookieExist"] = false
	data["already_liked"] = false
}
