package core

import (
    "fmt"
  //  "io/ioutil"
    "log"
 //   "math/rand"
    // "os"
   "strconv"
    "strings"
   // "time"
    "encoding/json"
    "github.com/TokTok/go-toxcore-c"
)

var debug bool = false


//更新自己的连接状态
func CallbackSelfConnectionStatus(t *tox.Tox, status int, userData interface{}) {
    
    log.Println("on self conn status:", status, userData)
    meInfo.ConnectStatus = status
    rjson := map2sk{
        "type":"CallbackSelfConnectionStatus",
        "data" : meInfo,
    }
    msg  ,err := json.Marshal( rjson)
    if err == nil{ 
        wssSendMessage( string(msg) )
    }else{
        log.Println("websocket send CallbackSelfConnectionStatuss error" , err.Error() )
    }
    
}

//接收好友请求
func CallbackFriendRequest(t *tox.Tox, friendId string, message string, userData interface{}){
        log.Println("on friend requests: ", friendId, message)
        err := dbConn.StoreFriendRequest(friendId , message , 1)
        if err != nil{
            log.Printf("db StoreFriendRequest fail friendId: %s message:%s \n" , friendId , message)
            return
        }

        rjson := map2sk{
            "type":"CallbackFriendRequest",
            "data":map2sk{
                "FriendId":friendId,
                "Message" : message,
                "IsIgnored" : 0,
                "Time":getNowTime(),
            },
            
        }
        msg  ,err := json.Marshal( rjson)
        if err == nil{ 
            wssSendMessage( string(msg) )
        }
}




//接收好友消息
func CallbackFriendMessage(t *tox.Tox, friendNumber uint32, message string, userData interface{}) {
    
    //friendId, err := t.FriendGetPublicKey(friendNumber)
    
    err := dbConn.StoreMessage(friendNumber , true , false, message)

    if err != nil{
        log.Printf("db StoreMessage fail friendId: %s message:%s \n" , friendNumber , message)
    }
    /*
    type Message struct {
    Id        int64
    Message    string
    IsIncoming bool
    IsAction   bool
    Time       int64
}
    */
    rjson := map2sk{
        "type":"CallbackFriendMessage",
        "FriendNumber":friendNumber,
        "data": map2sk{
            "Id": 0,
            "Message":message,
            "IsIncoming": 1,
            "IsAction": 0 ,
            "Time" : getNowTime(),
        },
    }

    msg  ,err := json.Marshal( rjson)
    if err== nil{ 
        wssSendMessage( string(msg) )
        log.Println("on friend message:", friendNumber, message)
    }else{

        log.Println("on friend send error message:", msg , err)
    }

}

//好友连接状态变更回调
func CallbackFriendConnectionStatus(t *tox.Tox, friendNumber uint32, status int, userData interface{}) {
    /*
    type Friend struct{
    friendNumber uint32
    friendId string
    Nickname string
    statusMessage string
    isOnline bool
    ToxStatus int
    addtime int64
}
    */
    
    friendInfo := getFriendInfo(friendNumber)
    friendInfo.IsOnline = status
    dictFriends[friendNumber] = friendInfo

    rjson := map2sk{
        "type":"CallbackFriendConnectionStatus",
        "data":friendInfo,
    }
    msg  ,err := json.Marshal( rjson)
    if err == nil{ 
        wssSendMessage( string(msg) )
    }
    log.Println("on friend connection status:", friendNumber, status, friendInfo.FriendId, err)
    
}

//好友在线状态变更回调
func CallbackFriendStatus(t *tox.Tox, friendNumber uint32, status int, userData interface{}) {
    friendInfo := getFriendInfo(friendNumber)
    friendInfo.ToxStatus = status
    dictFriends[friendNumber] = friendInfo

    rjson := map2sk{
        "type":"CallbackFriendStatus",
        "data":friendInfo,
    }
    msg  ,err := json.Marshal( rjson)
    if err == nil{ 
        wssSendMessage( string(msg) )
    }
     log.Println("on friend online status:", friendNumber, status, friendInfo.FriendId, err)
}

func CallbackFriendName(t *tox.Tox, friendNumber uint32, newName string, userData interface{}){

    friendInfo := getFriendInfo(friendNumber)
    friendInfo.Nickname = newName
    dictFriends[friendNumber] = friendInfo
    dbConn.UpdateFriend(friendInfo)
    rjson := map2sk{
        "type":"CallbackFriendName",
        "data":friendInfo,
    }
    msg  ,err := json.Marshal( rjson)
    if err == nil{ 
        wssSendMessage( string(msg) )
    }

    
    toxsave()
    log.Println("on friend change name:", newName, friendInfo.FriendId, err)
}

//好友签名变更回调
func CallbackFriendStatusMessage(t *tox.Tox, friendNumber uint32, statusText string, userData interface{}) {
    friendInfo := getFriendInfo(friendNumber)
    friendInfo.StatusMessage = statusText
    dictFriends[friendNumber] = friendInfo
    dbConn.UpdateFriend(friendInfo)
    rjson := map2sk{
        "type":"CallbackFriendStatusMessage",
        "data":friendInfo,
    }
    msg  ,err := json.Marshal( rjson)
    if err == nil{ 
        wssSendMessage( string(msg) )
    }
    
    toxsave()
    log.Println("on friend statusMessage change:", friendNumber, statusText, friendInfo.FriendId, err)
}






//文件传输控制

type FuncSendChunk  func(friendNumber uint32, fileNumber uint32, position uint64)
func makekey(no uint32, a0 interface{}, a1 interface{}) string {
    return fmt.Sprintf("%d_%v_%v", no, a0, a1)
}

var recvFiles = make(map[uint64]uint32, 0)
var sendFiles = make(map[uint64]uint32, 0)
var sendDatas = make(map[string][]byte, 0)
var chunkReqs = make([]string, 0)
var trySendChunk = func() ( FuncSendChunk ){
    // some vars for file echo
    
    trySendChunk := func(friendNumber uint32, fileNumber uint32, position uint64) {
        sentKeys := make(map[string]bool, 0)
        for _, reqkey := range chunkReqs {
            lst := strings.Split(reqkey, "_")
            pos, err := strconv.ParseUint(lst[2], 10, 64)
            if err != nil {
            }
            if data, ok := sendDatas[reqkey]; ok {
                r, err := t.FileSendChunk(friendNumber, fileNumber, pos, data)
                if err != nil {
                    if err.Error() == "toxcore error: 7" || err.Error() == "toxcore error: 8" {
                    } else {
                        log.Println("file send chunk error:", err, r, reqkey)
                    }
                    break
                } else {
                    delete(sendDatas, reqkey)
                    sentKeys[reqkey] = true
                }
            }
        }
        leftChunkReqs := make([]string, 0)
        for _, reqkey := range chunkReqs {
            if _, ok := sentKeys[reqkey]; !ok {
                leftChunkReqs = append(leftChunkReqs, reqkey)
            }
        }
        chunkReqs = leftChunkReqs
    }
    return trySendChunk 
}()

func CallbackFileRecvControl(t *tox.Tox, friendNumber uint32, fileNumber uint32,control int, userData interface{}) {
        friendId, err := t.FriendGetPublicKey(friendNumber)
        log.Println("on recv file control:", friendNumber, fileNumber, control, friendId, err)

        key := uint64(uint64(friendNumber)<<32 | uint64(fileNumber))
        if control == tox.FILE_CONTROL_RESUME {
            if fno, ok := sendFiles[key]; ok {
                t.FileControl(friendNumber, fno, tox.FILE_CONTROL_RESUME)
            }
        } else if control == tox.FILE_CONTROL_PAUSE {
            if fno, ok := sendFiles[key]; ok {
                t.FileControl(friendNumber, fno, tox.FILE_CONTROL_PAUSE)
            }
        } else if control == tox.FILE_CONTROL_CANCEL {
            if fno, ok := sendFiles[key]; ok {
                t.FileControl(friendNumber, fno, tox.FILE_CONTROL_CANCEL)
            }
        }
    }

func CallbackFileRecv(t *tox.Tox, friendNumber uint32, fileNumber uint32, kind uint32,
        fileSize uint64, fileName string, userData interface{}) {
        friendId, err := t.FriendGetPublicKey(friendNumber)
        log.Println("on recv file:", friendNumber, fileNumber, kind, fileSize, fileName, friendId, err)

        if fileSize > 1024*1024*1024 {
            // good guy
        }

        var reFileName = "Re_" + fileName
        reFileNumber, err := t.FileSend(friendNumber, kind, fileSize, reFileName, reFileName)
        if err != nil {
        }
        recvFiles[uint64(uint64(friendNumber)<<32|uint64(fileNumber))] = reFileNumber
        sendFiles[uint64(uint64(friendNumber)<<32|uint64(reFileNumber))] = fileNumber
    }

func CallbackFileRecvChunk(t *tox.Tox, friendNumber uint32, fileNumber uint32,
        position uint64, data []byte, userData interface{}) {
        friendId, err := t.FriendGetPublicKey(friendNumber)
        if debug {
            // log.Println("on recv chunk:", friendNumber, fileNumber, position, len(data), friendId, err)
        }

        if len(data) == 0 {
            if debug {
                log.Println("recv file finished:", friendNumber, fileNumber, friendId, err)
            }
        } else {
            reFileNumber := recvFiles[uint64(uint64(fileNumber)<<32|uint64(fileNumber))]
            key := makekey(friendNumber, reFileNumber, position)
            sendDatas[key] = data
            trySendChunk(friendNumber, reFileNumber, position)
        }
    }


func CallbackFileChunkRequest(t *tox.Tox, friendNumber uint32, fileNumber uint32, position uint64,
        length int, userData interface{}) {
        friendId, err := t.FriendGetPublicKey(friendNumber)
        if length == 0 {
            if debug {
                log.Println("send file finished:", friendNumber, fileNumber, friendId, err)
            }
            origFileNumber := sendFiles[uint64(uint64(fileNumber)<<32|uint64(fileNumber))]
            delete(sendFiles, uint64(uint64(fileNumber)<<32|uint64(fileNumber)))
            delete(recvFiles, uint64(uint64(fileNumber)<<32|uint64(origFileNumber)))
        } else {
            key := makekey(friendNumber, fileNumber, position)
            chunkReqs = append(chunkReqs, key)
            trySendChunk(friendNumber, fileNumber, position)
        }
    }

