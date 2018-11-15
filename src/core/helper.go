package core

import (
    "crypto/rand"
    "encoding/base64"
	"log"
	"time"
)

func RandomString(len int) (string, error) {
    bs := make([]byte, len)
    _, err := rand.Reader.Read(bs)
    if err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(bs), nil
}


func getFriendInfo(friendNumber uint32) Friend{

    friendId, err := t.FriendGetPublicKey(friendNumber)
    if err != nil{
        log.Printf("Not Found friend number :%d %s" , friendId , err.Error())
        return Friend{}
    }

    firendInfo , ok := dictFriends[friendNumber]
    if !ok {
        Nickname , err:= t.FriendGetName(friendNumber)
        statusMessage ,err := t.FriendGetStatusMessage(friendNumber)
        isonline , err := t.FriendGetConnectionStatus(friendNumber)
        toxstatus , err := t.FriendGetStatus(friendNumber)
        err = dbConn.StoreFriend(friendNumber , friendId , Nickname, statusMessage)
        if err != nil {
            log.Printf("CallbackFriendConnectionStatus store friend info fail! friendId:%s " , friendId)
        }
        firendInfo = Friend{friendNumber , friendId ,Nickname , statusMessage , isonline , toxstatus, "" , 0 , 0 }
    }
    return firendInfo
}


func toxsave(){
    t.WriteSavedata(tox_save_file)
}

func getNowTime() int64{
    return time.Now().Unix()*1000
}
