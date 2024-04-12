import { WebSocket } from 'ws';
import { decode, encode, initMsgBuilder, print, runCli } from "./common.js";

var host = "127.0.0.1"
var port = 9527
var ws = null
var timer = null

function connect(host = host, port = port) {
    ws = new WebSocket(`ws://${host}:${port}/`);
    ws.on('open', function open() {
        print('connect success');
        timer = setInterval(() => {
            ws.send(Buffer.alloc(0))
        }, 3000)
        send()
    });

    ws.on('message', function (data) {
        try {
            let [msgName, msg] = decode(data)
            print("recv server msg:", msgName, msg)
        } catch {
            print("decode error", data, data.length)
        }
    });

    ws.on('close', function () {
        print('websocket close');
        clearInterval(timer)
    });

    ws.on('error', function error(err) {
        print('websocket error ', err.message);
    });
}

function send(msgName = 'SayHelloReq', msgBody = { text: "hello, server!" }) {
    ws.send(encode(msgName, msgBody))
}

initMsgBuilder('pb')

connect(host, port)
runCli({ send, connect })