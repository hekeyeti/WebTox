package core

import (
    "github.com/TokTok/go-toxcore-c"
    "fmt"

)

//程序配置全局变量
const (
    CONFIG_PATH string  = "./conf/config.json"
    CERT_PUB    string  = "./cert/https.cert.pem"
    CERT_KEY    string  = "./cert/https.key.pem"
    HTML_DIR    string  = "./html"
    SAVE_DIR    string  = "./data"
)



const (
    ACCEPT_CONTACT = 1
    REJECT_CONTACT = 2
)

type map2sk map[string]interface{} 
type map2ik map[int]interface{}



type ToxInfo struct{
    Pubkey      string
    Seckey      string
    Toxid       string
    Nickname    string
    StatusMessage string
    ConnectStatus int
    ToxStatus int
}

type Message struct {
    Id        int64
   // FriendNumber uint32
    Message    string
    IsIncoming bool
    IsAction   bool
    Time       int64
}

type FriendRequest struct {
    FriendId string
    Message   string
    IsIgnored int
    Time int64
}

type Friend struct{
    FriendNumber uint32
    FriendId string
    Nickname string
    StatusMessage string
    IsOnline int
    ToxStatus int
    RemarkName string
    lastReadTime int64
    time int64
}

var tox_save_file string = fmt.Sprintf("%s/%s" , SAVE_DIR , "webtox.data" )
var tox_db_file string = fmt.Sprintf("%s/%s" , SAVE_DIR , "webtox.sqlite3" )
var tox_recv_path string = fmt.Sprintf("%s/%s" , SAVE_DIR , "recvfile")


//toxcore对象
var t *tox.Tox

//当前tox信息对象 
var meInfo *ToxInfo = new(ToxInfo)

//sqlite连接
var dbConn *StorageConn

//好友信息字典
var dictFriends map[uint32]Friend = make(map[uint32]Friend)

var http_auth_user *string 
var http_auth_pass *string











