var wsConn = new function(){
    const self = this;
    wss = new WebSocket(wspath)
   // log(wspath)
    self.conn = wss;
    wss.onopen = function(event){
      //  log(`wss open for ${wspath} event ${event}`)

    };
    wss.onmessage = function(event){
        const data = event.data
        log(`wss onmessage ${data}`)
        const data_type = typeof(data)
        switch(data_type){
            case "string":
                const rjson = JSON.parse(data)
                return reactor.OnMessageJson(rjson)
            break;
        }
        

    };
    wss.onerror = function(event){

    };
    wss.onclose = function(event){
        $("#modal-connection-error").modal("show")

    };

}

var view = new function(){
    const self = this;

  
    self.data_self  = {
        ConnectStatus : CONNECTION_NONE,
        Nickname : "",
        Pubkey : "",
        Seckey : "",
        StatusMessage : "",
        ToxStatus : USER_STATUS_NONE,
        Toxid : ""
    }

    self.data_contacts = {
        isFriendRequests : false,
        friends : {},
        friendRequests : {},
        reqfriendToxID : "",
        reqfriendStatusMesage : ""
    }

  
    self.data_chat = {
        isShow:false,
        friendNumber:null,
        SendMessage : ""
    }

    self.send_friend_request = {
        toxid : "",
        message : ""
    }

    //侧边栏的webtox信息
    const sidebar_top = new Vue({
        el: "#profile-card",
        data : self.data_self,
        methods : {

            set_status : function(ToxStatus){
                this.ToxStatus = ToxStatus;
                reactor.ChangeSelfInfo(self.data_self)
            },
            tox_self_change :function(event){
                reactor.ChangeSelfInfo(self.data_self)
            }
        }
    });

    //下边控制栏
    self.bottom_panel = new Vue({
        el:"#button-panel",
        methods:{
            
            group: function(){


            },
            file: function(){


            },
            setting: function(){
                self.data_chat.isShow =false;
                self.data_chat.friendNumber = null;
            },

        }
    })

    //弹出toxid的弹框信息
    self.modal_ToxId = new Vue({
        el:"#modal-toxid",
        data : self.data_self,
    })

    //好友列表信息
    self.sidebar_contacts = new Vue({
        el:"#contact-list-wrapper",
        data: self.data_contacts,
        methods:{
            select_chat : function(friendNumber){
                self.data_chat.friendNumber = friendNumber
                reactor.RequestFriendMessage(friendNumber)
                self.data_chat.isShow = true;
            }
        }
    })

    //好友请求信息列表
    self.modal_friend_request =new Vue({
        el:"#modal-friend-requests",
        data: self.data_contacts,
        methods : {
            accept : function(toxid){
                reactor.ProcessFriendRequest(1 , toxid)
            },
            reject : function(toxid){
                reactor.ProcessFriendRequest(0 , toxid)
            },
            addFriend: function(){
                const friendId = self.data_contacts.reqfriendToxID;
                const statusMessage = self.data_contacts.reqfriendStatusMesage;

                self.data_contacts.reqfriendToxID = ""
                self.data_contacts.reqfriendStatusMesage = ""
                reactor.SendFriendRequest(friendId , statusMessage)
            
            }
        }
    })

    //好友删除节目
    self.modal_friend_del = new Vue({
        el:"#modal-friend-del",
        data:{
            "chat":self.data_chat,
            "contacts":self.data_contacts,
        },
        methods : {
            delFriend: function(){
                const friendNumber = self.data_chat.friendNumber
                self.data_chat.friendNumber = null
                self.data_chat.isShow = false
                reactor.RemoveFriend(friendNumber)
            }
        }

    })

    //聊天窗口
    self.mainview = new Vue({
        el:"#mainview",
        data :{
            "self":self.data_self,
            "chat":self.data_chat,
            "contacts":self.data_contacts,

        },
        methods : {
            send_friend_messge : function(){
                const friendNumber = self.data_chat.friendNumber
                const SendMessage = self.data_chat.SendMessage
                if(!SendMessage){
                    return;
                }
                self.data_chat.SendMessage = ""
                reactor.SendFriendMessage(friendNumber , SendMessage)

            }
        }

    })

    // 监听 Ctrl + Enter发送消息 
    $(window).keydown(function (event) { 
         if (event.ctrlKey && event.keyCode == 13 ) {
            if(!self.data_chat.isShow){
                return;
            }
            const friendNumber = self.data_chat.friendNumber
            const SendMessage = self.data_chat.SendMessage
            if(!SendMessage){
                    return;
                }
            self.data_chat.SendMessage = ""
            reactor.SendFriendMessage(friendNumber , SendMessage)
         }
    });

}


var reactor = new function(){
    const self = this;
    const wss = wsConn.conn

    
    

    //接受自我状态变更
    const OnChangeSelfInfo = function(rjson){
        const vdict = rjson.data
        for(const key in vdict){
            const item = vdict[key]
            view.data_self[key] = item
        }
    }

    const OnFriendsList = function(rjson){
        const fdict = rjson.data
        for(var FirendNumber in fdict){
            view.data_contacts.friends[FirendNumber] = fdict[FirendNumber]
        }
        view.sidebar_contacts.$forceUpdate();
    }

    //好友请求列表
    const OnFriendRequstList = function(rjson){
        const list = rjson.data
        if(!list){
            return
        }
        view.data_contacts.isFriendRequests = true;

        for (var i = list.length - 1; i >= 0; i--) {
            const item = list[i]
            view.data_contacts.friendRequests[item.FriendNumber] = item
        }
        view.sidebar_contacts.$forceUpdate();
    }

    //单条好友请求
    const OnFriendRequsts = function(rjson){
        const data = rjson.data
        view.data_contacts.friendRequests[data.FriendId] = data
        view.data_contacts.isFriendRequests = true;
        view.modal_friend_request.$forceUpdate();
    }

    //好友列表变更
    const OnFriendsUpdata = function(rjson){
        const friend = rjson.data
        if(!friend){
            return
        }
        view.data_contacts.friends[friend.FriendNumber] = friend
        view.sidebar_contacts.$forceUpdate();
    }

    //接受server批量下发的聊天纪录
    const OnFriendMessageList = function(rjson){
        const friendNumber = rjson.friendNumber
        let messagelist = rjson.data
        if (!messagelist){
            messagelist = []
        }
        view.data_contacts.friends[friendNumber].message = messagelist
        view.mainview.$forceUpdate()
    }

    //单条信息回调
    const OnFriendMessage = function(rjson){
        const friendNumber = rjson.FriendNumber
        const message = rjson.data
        if(null == view.data_contacts.friends[friendNumber].message){
            view.data_contacts.friends[friendNumber].message = []
        } 
        view.data_contacts.friends[friendNumber].message.unshift(message)
        view.mainview.$forceUpdate()
    }


    //分发由Server端下发的JSON
    self.OnMessageJson =function(rjson){
        const type = rjson['type']
        const dispatch2type = {
            "wsAfterConnected":OnChangeSelfInfo,
            "CallbackSelfConnectionStatus":OnChangeSelfInfo,
            "wsFriendList" : OnFriendsList,
            "wsFriendRequestsList":OnFriendRequstList,
            "wsGetFriendMessage":OnFriendMessageList,
            "CallbackFriendMessage":OnFriendMessage,

            "CallbackFriendRequest" : OnFriendRequsts,


            "CallbackFriendConnectionStatus" : OnFriendsUpdata,
            "CallbackFriendStatus" : OnFriendsUpdata,
            "CallbackFriendName" : OnFriendsUpdata,
            "CallbackFriendStatusMessage" : OnFriendsUpdata,

        }
        const callbackfunc = dispatch2type[type]
        if(typeof(callbackfunc) == "function"){
            callbackfunc(rjson)
        }
    }


    /*===========================*/
    
    //更改名称，在线状态等信息
    self.ChangeSelfInfo = function(data){

        const up_json = {
            type:"wsUpdateSelfInfo",
            data : data
        }

        const updata = JSON.stringify(up_json)
        wss.send(updata)
    }
   
    self.ProcessFriendRequest = function(isAccept , toxid){
        delete view.data_contacts.friendRequests[toxid]
        view.data_contacts.isFriendRequests = false;
        const up_json = {
            type:"wsAcceptFriendRequest",
            toxid: toxid,
            isAccept : isAccept
        }
        const updata = JSON.stringify(up_json)
        wss.send(updata)
        view.modal_friend_request.$forceUpdate();
    }

    //请求Server下发聊天记录
    self.RequestFriendMessage = function(friendNumber){
        const up_json = {
            "type":"wsGetFriendMessage",
            "friendNumber":friendNumber
        }
        const updata = JSON.stringify(up_json)
        wss.send(updata)
    }

    self.SendFriendMessage = function(FirendNumber , message){
        const up_json = {
            "type":"wsSendFriendMsg",
            "message": message,
            "FriendNumber":FirendNumber
        }
        view.data_contacts.friends[FirendNumber].message.unshift({
            Message : message,
            IsIncoming : 0,
            IsAction : 0,
            Time:getUnixTimeStamp()
        })
        view.mainview.$forceUpdate()
        const updata = JSON.stringify(up_json)
        wss.send(updata)
    }

    self.SendFriendRequest = function(friendId , stausMessage){
        const up_json = {
            "type":"wsSendFriendRequest",
            "Message": stausMessage,
            "FriendId":friendId
        }
        const updata = JSON.stringify(up_json)
        wss.send(updata)
        
        view.modal_friend_request.$forceUpdate();
    }


    self.RemoveFriend = function(friendNumber){

         const up_json = {
            "type":"wsRemoveFriend",
            "FriendNumber":friendNumber
        }
        const updata = JSON.stringify(up_json)
        wss.send(updata)
        delete view.data_contacts.friends[friendNumber]
        view.sidebar_contacts.$forceUpdate();
        
    }


}