/*
|--------------------------------------------------------------------------
| sockio 模块
|--------------------------------------------------------------------------
 */
var _uid = $('meta[name="uid"]').attr('content');

var crm_io = {
    data: {},
    init: function(obj) {
        var socket = io(obj.public_socket);
        // uid可以是自己网站的用户id，以便针对uid推送以及统计在线人数

        // socket连接后以uid登录
        socket.on('connect', function(data) {
            socket.emit('login', _uid);
        });
        socket.on('online', function(data){
            $('.online_users').find('b').html(data.onlineCount);
        });
        socket.on('logout', function(data){
            $('.online_users').find('b').html(data.onlineCount);
        });
        // 后端推送来消息时
        socket.on('hz-pusher-channel:App\\Events\\PusherEvent', function(message) {
            if (!message.uid || message.uid == _uid) {
                if (message.uid == null) {
                    return;
                };
				console.log(message)
            }
        });
    }
}
