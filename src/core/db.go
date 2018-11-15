
package core

import (
    "database/sql"
    "errors"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "sync"
    "time"

)

var (
    KeyNotFound = errors.New("Key does not exist")
)


type StorageConn struct {
    db  *sql.DB
    mtx sync.Mutex
}



// Open creates a connection to the database
// always close the connection with `defer storageConn.Close()`
func DBOpen(filename string) (*StorageConn, error) {
    db, err := sql.Open("sqlite3", filename)
    if err != nil {
        log.Fatal(err)
        return &StorageConn{}, err
    }

    // create database tables
    sqlStmt := `
    CREATE TABLE IF NOT EXISTS contacts (
        friendNumber INTEGER ,
        friendId TEXT PRIMARY KEY,
        Nickname  TEXT,
        statusMessage TEXT,
        remarkName  TEXT,
        lastReadTime time,
        time INTEGER
    );
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY autoincrement,
        message TEXT NOT NULL,
        time INTEGER,
        fileName TEXT,
        friendNumber INTEGER,
        isIncoming INTEGER,
        isAction INTEGER,
        isFile INTEGER       
    );
    CREATE TABLE IF NOT EXISTS contact_requests (
        friendId TEXT NOT NULL PRIMARY KEY,
        message TEXT NOT NULL,
        isIgnored INTEGER,
        isIncoming INTEGER,
        time INTEGER
    );
    CREATE TABLE IF NOT EXISTS contact_group(
        groupNumber INTEGER,
        groupName TEXT,
        PRIMARY KEY(groupNumber)
    );
    CREATE TABLE IF NOT EXISTS contact_group_rel(
        id INTEGER PRIMARY KEY,
        groupNumber INTEGER,
        publicKey TEXT
    );
    CREATE TABLE IF NOT EXISTS keyValueStorage (
        key TEXT PRIMARY KEY,
        value TEXT
    );`

    _, err = db.Exec(sqlStmt)
    if err != nil {
        log.Panicf("%q: %s\n", err, sqlStmt)
        return &StorageConn{}, err
    }

    s := &StorageConn{db: db}
    return s, nil
}

// Close safely closes the connection to the database
func (s *StorageConn) Close() {
    s.db.Close()
}


func (s *StorageConn) StoreKeyValue(key string, value string) error {
    s.mtx.Lock()
    defer s.mtx.Unlock()

    _, err := s.db.Exec(`INSERT OR REPLACE INTO keyValueStorage(key, value) VALUES(?, ?)`, key, value)
    if err != nil {
        log.Print("[db StoreKeyValue] INSERT statement failed")
        return err
    }

    return nil
}


func (s *StorageConn) GetKeyValue(key string) (string, error) {
    rows, err := s.db.Query("SELECT value FROM keyValueStorage WHERE key = ?", key)
    if err != nil {
        log.Print("[db GetKeyValue] SELECT statement failed")
        return "", err
    }
    defer rows.Close()

    if rows.Next() {
        var value string
        rows.Scan(&value)
        return value, nil
    }

    return "", KeyNotFound
}

func (s *StorageConn) DeleteFriend(friendNumber uint32) error{
   
    _,err := s.db.Exec(`Delete from contacts where friendNumber = ? `, 
        friendNumber  )
    if err != nil {
        log.Print("[db DeleteFriend] Delete statement failed")
        return err
    }
    _,err = s.db.Exec(`Delete from messages where friendNumber = ?  `, 
        friendNumber  )
    if err != nil {
        log.Print("[db DeleteFriend] Delete statement failed")
        return err
    }
    return nil
}

func (s *StorageConn) UpdateFriend(friend Friend) error{

    _, err := s.db.Exec(`INSERT OR REPLACE INTO 
        contacts(friendNumber, friendId, Nickname , statusMessage ) 
        VALUES(?, ?, ? ,? )`, 
        friend.FriendNumber , friend.FriendId , friend.Nickname , friend.StatusMessage  )
    if err != nil {
        log.Print("[db UpdataFriend] REPLACE statement failed")
        return err
    }
    return nil
}


func (s *StorageConn) GetFriendList() []Friend{
    var Friends []Friend = []Friend{}

    rows , err := s.db.Query("SELECT * from contacts")
    if err != nil {
        log.Print("[db GetFriendList] SELECT statement failed")
        return Friends
    }
    defer rows.Close()


    
    for rows.Next(){
        var friendNumber uint32
        var friendId string
        var Nickname string
        var statusMessage string
        var addtime int64
        var remarkName string
        var lastReadTime int64
        rows.Scan(&friendNumber , &friendId , &Nickname , &statusMessage ,&remarkName , &lastReadTime , &addtime)
        Friends = append( Friends , Friend{ friendNumber , friendId , Nickname , statusMessage , 0 , 0 ,remarkName , lastReadTime, addtime } )
    }    
    return Friends 

}


//接受好友请求，保存db里面，tox端操作另外做
func (s *StorageConn) StoreFriend ( friendNumber uint32 ,friendId string , Nickname string , statusMessage string ) error{
    s.mtx.Lock()
    defer s.mtx.Unlock()

    NowTime :=time.Now().Unix()*1000


    _, err := s.db.Exec(`INSERT OR REPLACE INTO contacts(friendNumber, friendId, Nickname , statusMessage , time) VALUES(?, ?, ? ,? , ?)`, friendNumber ,friendId , Nickname , statusMessage , NowTime )
    if err != nil {
        log.Print("[db AcceptFriendRequest] INSERT statement failed")
        return err
    }

    
    return nil

}


func (s *StorageConn) StoreMessage(friendNumber uint32, isIncoming bool, isAction bool , message string) error {
    s.mtx.Lock()
    defer s.mtx.Unlock()

    _, err := s.db.Exec(`INSERT INTO 
        messages(friendNumber, isIncoming, isAction,isFile, time, message) 
        VALUES(?, ?, ? ,? , ?, ?)`, 
        friendNumber, isIncoming, isAction, 0 , time.Now().Unix()*1000, message)
    if err != nil {
        log.Print("[db StoreMessage] INSERT statement failed")
        return err
    }
    return nil
}


func (s *StorageConn) GetMessages(friendNumber uint32) []Message {

    rows, err := s.db.Query("SELECT id , message , time, isIncoming , isAction , isFile  FROM messages WHERE friendNumber = ? and isFile = 0 ORDER BY id desc", friendNumber)
    if err != nil {
        log.Print("[db GetMessages] SELECT statement failed")
        return nil
    }
    defer rows.Close()

    var messages []Message
    /*
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY autoincrement,
        friendNumber INTEGER,
        isIncoming INTEGER,
        isAction INTEGER,
        isFile INTEGER,
        fileName TEXT,
        message TEXT NOT NULL,
        time INTEGER
    );
    type Message struct {
    Id        int64
    Message    string
    IsIncoming bool
    IsAction   bool
    Time       int64
}
    */
    for rows.Next() {
        var id int64
       // var friendNumber int64
        var isIncoming bool
        var isAction bool
        var isFile bool
     //   var fileName string
        var time int64
        var message string
        rows.Scan(&id, &message , &time ,&isIncoming , &isAction , &isFile )
        messages = append(messages, Message{id , message , isIncoming , isAction  , time})
    }

    if messages == nil {
        return nil
    }

    return messages
}



func (s *StorageConn) StoreFriendRequest(friendId string,  message string ,isIncoming int) error {
    s.mtx.Lock()
    defer s.mtx.Unlock()
    NowTime :=time.Now().Unix()*1000
    _, err := s.db.Exec(`INSERT OR REPLACE INTO 
        contact_requests(friendId, message, isIncoming ,isIgnored , time ) 
        VALUES(?, ?, ?, ?,?)`, 
        friendId, message, isIncoming ,0 , NowTime)
    if err != nil {
        log.Print("[db StoreFriendRequest] INSERT statement failed")
        return err
    }
    return nil
}

func (s *StorageConn) GetFriendRequests() []FriendRequest {

    rows, err := s.db.Query("SELECT friendId, message, isIgnored ,time FROM contact_requests WHERE isIgnored == 0 and isIncoming = 1 ORDER BY isIgnored DESC ", )
    if err != nil {
        log.Print("[db GetFriendRequests] SELECT statement failed")
        return nil
    }
    defer rows.Close()

    var friendRequests []FriendRequest


    for rows.Next() {
        var friendId string
        var message string
        var isIgnored int
        var time int64
        rows.Scan(&friendId, &message, &isIgnored ,&time)
        friendRequests = append(friendRequests, FriendRequest{friendId , message ,  isIgnored , time })
    }

    if friendRequests == nil {
        return nil
    }

    return friendRequests
}

//删除好友请求
func (s *StorageConn) DeleteFriendRequest(friendId string) error {
    s.mtx.Lock()
    defer s.mtx.Unlock()

    _, err := s.db.Exec(`DELETE FROM contact_requests WHERE friendId = ?`, friendId)
    if err != nil {
        log.Print("[db DeleteFriendRequest] DELETE statement failed")
        return err
    }
    return nil
}


//接受好友请求，保存db里面，tox端操作另外做
func (s *StorageConn) AcceptFriendRequest( friendId string ) error{
    s.mtx.Lock()
    defer s.mtx.Unlock()

    _, err := s.db.Exec(`UPDATE  contact_requests set isIgnored = ? WHERE friendId = ?`, ACCEPT_CONTACT, friendId)
    if err != nil {
        log.Print("[db AcceptFriendRequest] UPDATE statement failed")
        return err
    }
    return nil

}

//拒绝好友请求DB入库
func (s *StorageConn) RejectFriendRequset(friendId string ) error {
     s.mtx.Lock()
    defer s.mtx.Unlock()

    _, err := s.db.Exec(`UPDATE  contact_requests SET isIgnored = ? WHERE friendId = ?`,REJECT_CONTACT  ,  friendId)
    if err != nil {
        log.Print("[db RejectFriendRequset] UPDATE statement failed")
        return err
    }
    return nil

}



