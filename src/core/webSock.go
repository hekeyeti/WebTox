package core



import (
    "log"
   // "github.com/codedust/go-tox"
    "golang.org/x/net/websocket"
    "encoding/json"
   // "strconv"
   // "time"
)



var wsConns = make( map[ *websocket.Conn ]  bool )

//下发数据
func wssSendMessage(msg string) {
    go func() {
        for conn, _ := range wsConns {
            if err := websocket.Message.Send(conn, msg); err != nil {
                log.Println("[handleWS] Could not send message to ", conn.RemoteAddr, err.Error())
            }
        }
    }()
}


//分发处理函数
var routeWSUP map[string]func(*websocket.Conn , map2sk)=map[string]func(*websocket.Conn , map2sk){
  `wsUpdateSelfInfo`:wsUpdateSelfInfo,
  `wsSendFriendRequest`:wsSendFriendRequest,
  `wsGetFriendRequest`:wsGetFriendRequest,
  `wsAcceptFriendRequest`:wsAcceptFriendRequest,
  `wsSendFriendMsg`:wsSendFriendMsg,
  `wsGetFriendMessage`:wsGetFriendMessage,
  `wsRemoveFriend`:wsRemoveFriend,
}

func handlerWs(conn *websocket.Conn , clientMessage string){
    log.Println("ws upload message  ", clientMessage)
    
    dictWSUP := make(map2sk)
    json.Unmarshal([]byte(clientMessage) , &dictWSUP) 
    callbackNameIn ,ok  := dictWSUP["type"]
    if !ok{
       log.Println("ws upload message format error, data :" , clientMessage)
       return
    }
    callbackName , ok := callbackNameIn.(string)
    callbackfunc , ok := routeWSUP[callbackName]
    if !ok{
      log.Println("ws upload message callbackName undefine! data:" , clientMessage)
      return
    }
   
    callbackfunc(conn , dictWSUP)

}



var handleWS = websocket.Handler(func(conn *websocket.Conn) {
    var err error
    var clientMessage string

    // cleanup on server side
    defer func() {
        //从列表中移除ws连接
        delete(   wsConns ,conn )
        if err = conn.Close(); err != nil {
            log.Println("[handleWS] Websocket could not be closed", err.Error())
        }
        log.Println("[handleWS] Number of clients still connected:", len(wsConns))
    }()

    log.Println("[handleWS] Client connected:", conn.Request().RemoteAddr)

    //如果websocket的连接数大于1,保证当前只有一个连接，则返回错误信息
    if len(wsConns) > 0{
        msg  , _ := json.Marshal( map2sk{
                "type": "init" ,
                "value" : "error" , 
                "message":"To many client for webtox server!",
            }  )
        if  err =  websocket.Message.Send(conn , string(msg) ) ; err != nil{
            log.Println("[handleWS] Could not send init  message to ", conn.RemoteAddr, err.Error())
        }
        log.Println(" To many client for webtox! ")
        return
    }

    //把当前连接写入websocket连接列表
    wsConns[conn] = true
    log.Println("[handleWS] Number of clients connected:", len(wsConns))

    //下发当前信息，初始化Tox页面
    wsAfterConnected(conn)

    for {
       // log.Println("web socket start Receive!")
        if err = websocket.Message.Receive(conn, &clientMessage); err != nil {
            // the connection is closed
            log.Println("[handleWS] Read error. Removing client.", err.Error())
            //tox.SelfSetStatus(gotox.TOX_USERSTATUS_AWAY)
            
            return
        }
        handlerWs(conn , clientMessage)
        //log.Printf("web socket receive msg :%s \n" , clientMessage )
       // broadcastToClients("xxxxxxxx")

    }
})








