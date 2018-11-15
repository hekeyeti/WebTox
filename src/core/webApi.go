package core



import (
    "fmt"
    "net/http"
    //"regexp"
    //"log"
    //"io/ioutil"
)


/*
下发头像图片
*/
func pushAvatar(w http.ResponseWriter, r *http.Request){    
    w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Fprintf(w , "test %s" , 11)

}

/*
下发接收的文件
*/
func pushFile(w http.ResponseWriter, r *http.Request){


}

/*
导出聊天记录
*/
func exportChatRecord(w http.ResponseWriter, r *http.Request){


}




