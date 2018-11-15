package core

import (
    "fmt"
    "net/http"
    //"html/template"
    authz "../authz"
    "regexp"
   // "strings"
    "log"
   // "io/ioutil"
)



func HttpListen(cfg JConfig){
    
    var pat_ip *regexp.Regexp= regexp.MustCompile(`^[^:]+`)
    go func(){
        //所有请求都跳转到SSL中
        JumpForHttps := func (w http.ResponseWriter, r *http.Request) {
            log.Printf("http connect addr:%s\n",r.RemoteAddr)
            
            host :=pat_ip.FindString(r.Host)
            ssl_location := fmt.Sprintf("https://%s%s/" , host , cfg.HTTPS_LISTEN_ADDR)
           // w.Header().Set("Location", ssl_location)
            //w.WriteHeader(301)   
            log.Println(ssl_location)
             http.Redirect(w, r, ssl_location, http.StatusMovedPermanently)
        }
        http.HandleFunc("/", JumpForHttps) //设置访问的路由
        log.Printf("start http server addr:%s!\n" ,cfg.HTTP_LISTEN_ADDR )
        err := http.ListenAndServe(cfg.HTTP_LISTEN_ADDR , nil) //设置监听的端口
        if err != nil {
            log.Fatal("ListenAndServe: ", err)
        }
    }()

}

// 路由定义
type routeInfo struct {
    pattern string                                       // 正则表达式
    call       func(w http.ResponseWriter, r *http.Request) //Controller函数
}

// 路由定义
var routePath = []routeInfo{
    routeInfo{"^/api/avatar$",  pushAvatar },
    routeInfo{"^/api/ws$",  handleWS.ServeHTTP },
    routeInfo{".*" ,  http.FileServer(http.Dir(HTML_DIR)).ServeHTTP },
}

func Page404(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(404) 
    fmt.Fprint(w, "404 Page Not Found!")
}

func Loop2WebServer(cfg JConfig){
    var err error
    cert_pub := CERT_PUB
    cert_key := CERT_KEY

    /*user := cfg.HTTPS_AUTH_USER
    passwd := cfg.HTTPS_AUTH_PASS*/


    listenAddr := cfg.HTTPS_LISTEN_ADDR

    

    salt  , err := RandomString(32)
    if err != nil{
        log.Printf("RandomString error :%s \n" , err.Error())
        salt = "AAA"
    }
    //log.Printf(" web auth user:%s passwd:%s salt:%s ",user, passwd , salt)
    
    authOptions := authz.NewAuthOptions(*http_auth_user, *http_auth_pass , salt)
    mux := http.NewServeMux()
    //路由代理
    Route :=func (w http.ResponseWriter,  r *http.Request){
        Method := r.Method
        RequestURI   := r.RequestURI
        log.Printf("http access method:%s url:%s" , Method , RequestURI)
        for _, p := range routePath {
            // 这里循环匹配Path，先添加的先匹配
            reg, err := regexp.Compile(p.pattern)
            if err != nil {
                continue
            }
            if reg.MatchString(r.URL.Path) {
                p.call(w, r)
                return
            }
        }
        Page404(w , r)
}

    // paths that require authentication
    //mux.Handle("/events", authz.BasicAuthHandler(http.HandlerFunc(web.SayhelloName), authOptions))
  //  mux.Handle("/api/", authz.BasicAuthHandler(web.handleWS, authOptions))
    mux.Handle("/", authz.BasicAuthHandler( http.HandlerFunc(Route)  , authOptions  ) )
    //mux.Handle("/", authz.BasicAuthHandler(http.FileServer(http.Dir("html/")), authOptions))
    // paths that *do not* require authentication
    //如果添加不经过http auth认证的hander
   // mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("html/img"))))

    /*authz.CreateCertificateIfNotExist(cert_pub, cert_key, "localhost", 3072)
    log.Printf("Web Tox Start Https Sever Addr:%s\n" , listenAddr)*/
    
    //启动http服务器
    HttpListen(cfg)

    //启动HTTPS服务器
    log.Printf("start https server addr:%s!\n",listenAddr)
    err = authz.ListenAndUpgradeTLS(listenAddr, cert_pub, cert_key, mux)
    if err != nil{
        log.Printf("Web Tox Https Server Listen Addr: %s Fail! error :%s \n" , listenAddr , err.Error())
    }
}
