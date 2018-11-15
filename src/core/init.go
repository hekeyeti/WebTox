package core

import (
    "encoding/json"
    "io/ioutil"
    "errors"
    "log"
    "flag"
)

type JConfig struct{
    DHT_BOOT_PUBKEY string `json:"DHT_BOOT_PUBKEY"`
    DHT_BOOT_ADDR string `json:"DHT_BOOT_ADDR"`
    TOX_PORT int `json:"TOX_PORT"`

    HTTPS_LISTEN_ADDR string `json:"HTTPS_LISTEN_ADDR"`
    HTTP_LISTEN_ADDR string `json:"HTTP_LISTEN_ADDR"`
}



func InitConfig(cfg *JConfig) (error) {


    http_auth_user = flag.String("u","tox","http auth login user")
    http_auth_pass = flag.String("p","tox","http auth login passwd")
    flag.Parse()


    pdb , err := DBOpen(tox_db_file)
    dbConn = pdb
    //dbConn.StoreKeyValue("tox" , "true")
    if err != nil{
        log.Println("Cannot open sqlite3 file :" ,  tox_db_file)
        return errors.New("dbConn init error!")
    }else{
        log.Println("dbConn init success ! " ,  tox_db_file)

    }

    json_content, err := ioutil.ReadFile(CONFIG_PATH)
    if err != nil {
        return  err
    }
    
    if err := json.Unmarshal(json_content, cfg); err == nil {
        return nil
    } else {
        return err
    }

}




