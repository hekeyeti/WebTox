package main

import (
    "./core"
  //  util "./util"
   // "net/http" 
   //  "fmt"
     "log"
    // "os"

   // "net/http"
   // "strings"
)




func main() {

    //Parse config.json
    var cfg core.JConfig
    err := core.InitConfig(&cfg)

    if err != nil {
        log.Printf("load config error:%s \n !" , err.Error() )
        return
    }
    
    //开启TOX服务
    log.Println("start tox service!")
    core.ToxCoreStart(cfg)

    log.Println("start web server!")
    //开启Web服务
    core.Loop2WebServer(cfg)
}