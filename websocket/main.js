var Chat = (function() {
    function Chat(url, app_id, user_id, message_token, options) {
        this.app_id = app_id;
        this.user_id = user_id;
        this.message_token = message_token;
        this.conn = new WebSocket(url);
        this.conn.onclose = function() {
            console.log('conn closed');
            return
        }
        this.conn.onmessage = function() {
            
        }
        if (options) {
            this.onTextMessage = options.onTextMessage;
            this.onImageMessage = options.onImageMessage;
            this.onOnlineMessage = options.onOnlineMessage;
            this.onOfflineMessage = options.onOfflineMessage;
        }
    }
    Chat.prototype.sendTextMessage = function(channel_id, content) {
        console.log('sendTextMessage', channel_id, content);
    }
    Chat.prototype.sendImageMessage = function(channel_id, image_url) {
        console.log('sendImageMessage', channel_id, image_url)
    }
    return Chat;
}());
var chat = new Chat('app1', 'u1', 'mt1', {
    onTextMessage: function(a) {
        console.log(a);
    }
});
chat.sendTextMessage('ch1', 'hello');
