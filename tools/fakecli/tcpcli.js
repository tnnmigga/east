import { decode, encode, initMsgBuilder, print, runCli } from "./common.js";
import { Socket } from 'net';

var host = "127.0.0.1"
var port = 9527

const pkgSizeByteLen = 4
var buf = Buffer.alloc(0)

var socket = new Socket()

initMsgBuilder('pb')

const ping = Buffer.alloc(pkgSizeByteLen)
ping.writeUint32LE(0)

function connect(host, port) {
    socket = new Socket()
    socket.connect(port, host, () => {
        print("connect success")
        send()
        setInterval(() => {
            socket.write(ping)
        }, 3000)
    })
    socket.on("data", (data) => {
        print("on data")
        buf = Buffer.concat([buf, data])
        if (buf.length < pkgSizeByteLen) {
            return
        }
        let pkgSize = buf.readUInt32LE()
        let sliceLen = pkgSize + pkgSizeByteLen
        if (buf.length < sliceLen) {
            return
        }
        let pkg = buf.subarray(pkgSizeByteLen, sliceLen)
        buf = buf.subarray(sliceLen)
        try {
            let [msgName, msg] = decode(pkg)
            print("recv server msg:", msgName, msg)
        } catch {
            print("decode error", data, data.length)
        }
    })
}

function send(msgName = 'SayHelloReq', msgBody = { text: "hello, server!" }) {
    let b = encode(msgName, msgBody)
    let a = Buffer.alloc(pkgSizeByteLen)
    a.writeUInt32LE(b.length)
    socket.write(Buffer.concat([a, b]))
}

connect(host, port)
runCli({ send, connect })

