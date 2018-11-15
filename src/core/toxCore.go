package core

import (
    "fmt"
    "io/ioutil"
    "log"
    "strconv"
    "strings"
    "time"
    "github.com/TokTok/go-toxcore-c"
)

/*
type ToxOptions struct {
    Ipv6_enabled            bool
    Udp_enabled             bool
    Proxy_type              int32
    Proxy_host              string
    Proxy_port              uint16
    Savedata_type           int
    Savedata_data           []byte
    Tcp_port                uint16
    Local_discovery_enabled bool
    Start_port              uint16
    End_port                uint16
    Hole_punching_enabled   bool
    ThreadSafe              bool
    LogCallback             func(_ *Tox, level int, file string, line uint32, fname string, msg string)
}

*/


func ToxCoreStart(cfg JConfig) {
    

    log.SetFlags(log.Flags() | log.Lshortfile)
    bootaddr   :=  strings.Split(cfg.DHT_BOOT_ADDR,":")
    boothost   := bootaddr[0]
    bootport ,err   := strconv.Atoi(bootaddr[1])
    if err != nil {
        log.Println(" bootaddr format error !  eg. 198.98.51.198:33445 " , bootaddr)
        return
    }

 


    opt := tox.NewToxOptions()
    if tox.FileExist(tox_save_file) {
        data, err := ioutil.ReadFile(tox_save_file)
        if err != nil {
            log.Println("Tox savedata read fail! ")
            log.Println( err.Error() )
            return
        } else {
            opt.Savedata_data = data
            opt.Savedata_type = tox.SAVEDATA_TYPE_TOX_SAVE
        }
    }

    opt.Tcp_port = uint16(cfg.TOX_PORT)

    t = tox.NewTox(opt)
    if t == nil {
        log.Println(" Tox Core Init Error!")
        return
    }


    r, err := t.Bootstrap( boothost , uint16(bootport) , cfg.DHT_BOOT_PUBKEY )
    r2, err := t.AddTcpRelay( boothost , uint16(bootport) , cfg.DHT_BOOT_PUBKEY )
    log.Println("bootstrap:", r, err, r2)

    pubkey := t.SelfGetPublicKey()
    seckey := t.SelfGetSecretKey()
    toxid := t.SelfGetAddress()

    //log.Println("keys:", pubkey, seckey, len(pubkey), len(seckey))
    log.Println("toxid:", toxid)

    //读取保存的用户名
    nickName := t.SelfGetName()
    //log.Printf("tox savedata nickname :%s \n " , nickName)
    if "" == nickName{
        t.SelfSetName("webToxUser")
        nickName = t.SelfGetName()
        //log.Printf("tox set default nickname :%s \n " , nickName)
    }

    //读取保存的签名
    statusMessage , err := t.SelfGetStatusMessage()
    if err != nil{
        t.SelfSetStatusMessage("webTox Useing!")
    }
    
    statusMessage, err = t.SelfGetStatusMessage()
    //log.Printf( "Tox StatusText: %s  \n" , statusMessage )
    


  //  sz := t.GetSavedataSize()
   // sd := t.GetSavedata()
   // log.Println("savedata:", sz, t)
    //log.Println("savedata", len(sd), t)

    err = t.WriteSavedata(tox_save_file)
    //log.Println("savedata write:", err)
    

    friends := dbConn.GetFriendList()
    for _,firendinfo := range friends{

        dictFriends[firendinfo.FriendNumber] = firendinfo
    }
    // add friend norequest
    fv := t.SelfGetFriendList()
    for _, fno := range fv {
        fid, err := t.FriendGetPublicKey(fno)
        if err != nil {
            log.Println(err)
        } else {
            t.FriendAddNorequest(fid)
        }
    }
    log.Println("add friends:", len(fv))

    

    // callbacks
    t.CallbackSelfConnectionStatus(CallbackSelfConnectionStatus, nil)
    t.CallbackFriendRequest(CallbackFriendRequest, nil)
    t.CallbackFriendMessage(CallbackFriendMessage, nil)
    t.CallbackFriendConnectionStatus(CallbackFriendConnectionStatus, nil)
    t.CallbackFriendName(CallbackFriendName,nil)
    t.CallbackFriendStatus(CallbackFriendStatus, nil)
    t.CallbackFriendStatusMessage(CallbackFriendStatusMessage, nil)

    //fileControlCallback
    t.CallbackFileRecvControl(CallbackFileRecvControl, nil)
    t.CallbackFileRecv(CallbackFileRecv, nil)
    t.CallbackFileRecvChunk(CallbackFileRecvChunk, nil)
    t.CallbackFileChunkRequest(CallbackFileChunkRequest, nil)

    


   //写入全局变量
  // HandlerTox.toxavHandler = av
   meInfo.Pubkey = pubkey
   meInfo.Toxid = toxid
   meInfo.Seckey = seckey
   meInfo.Nickname = nickName
   meInfo.StatusMessage = statusMessage
   meInfo.ToxStatus = t.SelfGetStatus()
   

    //runing tox loop
    go func(){
        loopc := 0
        itval := 0
        for  {
            iv := t.IterationInterval()
            if iv != itval {
                if debug {
                    if itval-iv > 20 || iv-itval > 20 {
                        log.Println("tox itval changed:", itval, iv)
                    }
                }
                itval = iv
            }

            t.Iterate()
            status := t.SelfGetConnectionStatus()
            if loopc%5500 == 0 {
                if status == 0 {
                    if debug {
                        fmt.Print(".")
                    }
                } else {
                    if debug {
                        fmt.Print(status, ",")
                    }
                }
            }
            loopc += 1
            time.Sleep(1000 * 50 * time.Microsecond)
        }
        t.Kill()
    }()
    
    

}





































