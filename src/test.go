package main


import (

    "fmt"
    "encoding/json"
   // .  "./web"
)

type TINTFOO func(int , int )(int )


var foo = func() (  TINTFOO  ){
    return func(a int , b int ) (int) {
        return a +b
    }

}()


var tox_save_file string = fmt.Sprintf("%s/%s" ,"data" , "webtox.data" )


func main(){
go func( a int ){
        fmt.Printf(" go thead call args %d \n " , a)

    }(200)

    type dmap map[string]interface{} 

    var tmap dmap = dmap{
        "mm": dmap{"cc":"dd"} ,
        "cc":dmap{"ff":1},
    }

    msg  , _ := json.Marshal( tmap )
    fmt.Printf(" json string %s \n" , string(msg)  )

    fmt.Printf(" call test functor %d \n " , foo (1,3) )
    fmt.Printf( "tox_save_file :%s \n" , tox_save_file )

    

}







