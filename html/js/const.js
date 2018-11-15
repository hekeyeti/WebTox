const log = console.log
const myhost =window.location.host
const wspath = `wss://${myhost}/api/ws`
const getUnixTimeStamp = function(){return Math.floor(new Date().getTime() ) }
const getLocalTime = function (nS) {    return new Date(parseInt(nS) * 1000).toLocaleString().substr(0,17)} 

/*
CONNECTION_NONE :0
CONNECTION_TCP  :1
CONNECTION_UDP  :2
*/

const CONNECTION_NONE = 0
const CONNECTION_TCP = 1
const CONNECTION_UDP = 2

/*
USER_STATUS_NONE :0
USER_STATUS_AWAY :1
USER_STATUS_BUSY :2
*/

const USER_STATUS_BUSY = 2
const USER_STATUS_AWAY = 1
const USER_STATUS_NONE = 0



