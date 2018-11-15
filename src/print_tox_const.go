package main

import (
    "fmt"
    . "github.com/TokTok/go-toxcore-c"
)



var list_name []string = []string{
	`PUBLIC_KEY_SIZE           `,
	`SECRET_KEY_SIZE           `,
	`ADDRESS_SIZE              `,
	`MAX_NAME_LENGTH           `,
	`MAX_STATUS_MESSAGE_LENGTH `,
	`MAX_FRIEND_REQUEST_LENGTH `,
	`MAX_MESSAGE_LENGTH        `,
	`MAX_CUSTOM_PACKET_SIZE    `,
	`HASH_LENGTH               `,
	`FILE_ID_LENGTH            `,
	`MAX_FILENAME_LENGTH   `,
	`USER_STATUS_NONE `,
	`USER_STATUS_AWAY `,
	`USER_STATUS_BUSY `,
	`CONNECTION_NONE `,
	`CONNECTION_TCP  `,
	`CONNECTION_UDP  `,
	`FILE_CONTROL_RESUME `,
	`FILE_CONTROL_PAUSE  `,
	`FILE_CONTROL_CANCEL `,
	`FILE_KIND_DATA   `,
	`FILE_KIND_AVATAR `,
	`MESSAGE_TYPE_NORMAL `,
	`MESSAGE_TYPE_ACTION `,
	`CONFERENCE_TYPE_TEXT `,
	`CONFERENCE_TYPE_AV   `,
}

var list_number []interface{} = []interface{}{
PUBLIC_KEY_SIZE      ,
SECRET_KEY_SIZE      ,
ADDRESS_SIZE         ,
MAX_NAME_LENGTH      ,
MAX_STATUS_MESSAGE_LENGTH,
MAX_FRIEND_REQUEST_LENGTH,
MAX_MESSAGE_LENGTH   ,
MAX_CUSTOM_PACKET_SIZE,
HASH_LENGTH          ,
FILE_ID_LENGTH       ,
MAX_FILENAME_LENGTH  ,
USER_STATUS_NONE ,
USER_STATUS_AWAY ,
USER_STATUS_BUSY ,
CONNECTION_NONE ,
CONNECTION_TCP  ,
CONNECTION_UDP  ,
FILE_CONTROL_RESUME ,
FILE_CONTROL_PAUSE  ,
FILE_CONTROL_CANCEL ,
FILE_KIND_DATA   ,
FILE_KIND_AVATAR ,
MESSAGE_TYPE_NORMAL ,
MESSAGE_TYPE_ACTION ,
CONFERENCE_TYPE_TEXT ,
CONFERENCE_TYPE_AV   ,
}

func main(){
	for index,val := range(list_name){
		fmt.Printf("%s:%d\n" ,val , list_number[index] )
	}
}