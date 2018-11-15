package core



import (
    "log"
   // "github.com/codedust/go-tox"
    "golang.org/x/net/websocket"
    "encoding/json"
   // "strconv"
   // "time"
)



//刚连接上之后的初始化信息
func wsAfterConnected(conn *websocket.Conn){


   rjson := map2sk{
        "type" : "wsAfterConnected",
        "data": *meInfo,
   }

   msg,err := json.Marshal(rjson)
   if err != nil{
      log.Println(" wsAfterConnected Error "  , err.Error() )
      return 
   }


   if err := websocket.Message.Send(conn , string(msg) ); err != nil{
        log.Println("wsFriendList Send Error " , err.Error())
   }

   rjson = map2sk{
        "type" : "wsFriendList",
        "data": dictFriends,
   }

   msg,err = json.Marshal(rjson)
   if err != nil{
      log.Println(" wsFriendList Error "  , err.Error() )
      return 
   }


   if err = websocket.Message.Send(conn , string(msg) ); err != nil{
        log.Println("wsFriendList Send Error " , err.Error())
   }

   friendRequests := dbConn.GetFriendRequests()
   
   rjson = map2sk{
        "type" : "wsFriendRequestsList",
        "data": friendRequests,
   }

   msg,err = json.Marshal(rjson)
   if err != nil{
      log.Println(" wsFriendRequestsList Error "  , err.Error() )
      return 
   }

   if err = websocket.Message.Send(conn , string(msg) ); err != nil{
        log.Println("wsFriendRequestsList Send Error " , err.Error())
   }

   log.Println(" wsAfterConnected success "  )

}

//设置Tox状态
func wsUpdateSelfInfo(conn *websocket.Conn , recvMap map2sk ){
    for key,_ := range recvMap{
        log.Println( " wsUpdateSelfInfo key ", key )
    }

    data, ok := recvMap["data"]
    if !ok{
        log.Println("upload json data none!")
        return
    }
    rjson , ok  := data.(map[string]interface{})

    Nickname := rjson[`Nickname`]
    ToxStatus := rjson[`ToxStatus`]
    StatusMessage :=rjson[`StatusMessage`]

    t.SelfSetStatus( uint8(int(ToxStatus.(float64))) )
    t.SelfSetName(Nickname.(string))
    t.SelfSetStatusMessage(StatusMessage.(string))

    err := t.WriteSavedata(tox_save_file)
    if err != nil{
        log.Println("tox data save fail!" , err)
        return
    }
    meInfo.Nickname = Nickname.(string)
    meInfo.StatusMessage = StatusMessage.(string)
    meInfo.ToxStatus  = int(ToxStatus.(float64))

    log.Println("wsUpdateSelfInfo success!")
}

//获取好友请求列表
func wsGetFriendRequest(conn *websocket.Conn , recvMap map2sk){

    FriendRequests := dbConn.GetFriendRequests()

    rjson := map2sk{
        "type":"wsGetFriendRequest",
        "data":FriendRequests,
    }

    msg,err := json.Marshal(rjson)
   if err != nil{
      log.Println("wsGetFriendRequest Error "  , err.Error() )
      return 
   }


   if err = websocket.Message.Send(conn , string(msg) ); err != nil{
      log.Println("wsGetFriendRequest Send Error " , err.Error())
   }

}


//获取单个好友的信息列
func wsGetFriendMessage(conn *websocket.Conn , recvMap map2sk ){
    fnm := recvMap["friendNumber"]
    friendNumber := uint32(fnm.(float64))
    MessageList := dbConn.GetMessages(friendNumber  )


    rjson := map2sk{
        "type":"wsGetFriendMessage",
        "friendNumber": friendNumber,
        "data":MessageList,
    }


    msg,err := json.Marshal(rjson)
   if err != nil{
      log.Println("wsGetFriendMessage Error "  , err.Error() )
      return 
   }


   if err = websocket.Message.Send(conn , string(msg) ); err != nil{
      log.Println("wsGetFriendMessage Send Error " , err.Error())
   }


}

//发送消息到联系人
func wsSendFriendMsg(conn *websocket.Conn , recvMap map2sk){
    fnm,_  :=  recvMap["FriendNumber"]
    msg, _ := recvMap["message"]

    friendNumber := uint32( fnm.(float64) )
    message := msg.(string)

    _ , err := t.FriendSendMessage(friendNumber , message)
    if err != nil{
        log.Printf("wsSendFriendMsg fail friendNumber:%d message:%s" , friendNumber , message)
    }

    dbConn.StoreMessage(friendNumber , false , false , message)

}


//发起好友请求
func wsSendFriendRequest(conn *websocket.Conn , recvMap map2sk ){
    id , _ := recvMap["FriendId"]
    msg , _ := recvMap["Message"]

    FriendId := id.(string)
    Message := msg.(string)


    friendNumber ,err := t.FriendAdd(FriendId , Message)

    if err != nil {
        log.Printf("wsSendFriendRequest error %s" , err.Error())
        return
    }

    dbConn.StoreFriend(friendNumber , FriendId , "" , "")
    toxsave()
    log.Printf("wsSendFriendRequest success %d %s  ", FriendId , Message)
}



//是否接受好友请求
func wsAcceptFriendRequest(conn *websocket.Conn , recvMap map2sk){

    Id,_  := recvMap["toxid"]
    Accept ,_ := recvMap["isAccept"]
    
    FriendId,_ := Id.(string)
    isAccept  := int(Accept.(float64) )
    log.Println(FriendId , isAccept)
    if 0 == isAccept{
        err := dbConn.RejectFriendRequset(FriendId)
        if err != nil{
            log.Printf("wsAcceptFriendRequest Reject stroe fail friendId : %s\n" , FriendId )
        }
        return 
    }

    friendNumber ,err := t.FriendAddNorequest(FriendId)

    if err != nil {
        log.Printf("wsAcceptFriendRequest Accept error %s" , err.Error())
        return
    }
    Nickname , err:= t.FriendGetName(friendNumber)
    StatusMessage ,err := t.FriendGetStatusMessage(friendNumber)
    dbConn.StoreFriend(friendNumber , FriendId , Nickname, StatusMessage)
    dbConn.AcceptFriendRequest(FriendId)
    
    log.Printf("wsAcceptFriendRequest success %d %s %s %s ", friendNumber , FriendId , Nickname , StatusMessage)
    toxsave()

}

//删除好友
func wsRemoveFriend(conn *websocket.Conn , recvMap map2sk){

    fn , _ := recvMap["FriendNumber"]
    FriendNumber := uint32( fn.(float64) )
    t.FriendDelete(FriendNumber)
    dbConn.DeleteFriend(FriendNumber)
    toxsave()
}

