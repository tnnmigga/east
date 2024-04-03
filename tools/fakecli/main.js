import protobuf from 'protobufjs';
import repl from 'repl'
import { readdirSync } from 'fs'
import { log } from 'console';
import { Socket } from 'net';

const msgBuilders = {}
const msgidToName = {}
var socket = new Socket()

var ping = Buffer.alloc(4)
ping.writeUint32LE(0)

function connect() {
    socket = new Socket()
    socket.connect('10078', '127.0.0.1', function () {
        print("connect success")
        send("SayHelloReq", { text: "hello, server!" })
        setInterval(function () {
            // print("ping")
            socket.write(ping)
        }, 3000)
    })
    socket.on("data", function (data) {
        try {
            let [msgName, msg] = decode(data)
            print("recv server msg:", msgName, msg)
        } catch {
            print("decode error", data, data.length)
        }
    })
}

function initMsgBuilder(path) {
    let files = readdirSync(path)
    for (let name of files) {
        if (!name.endsWith(".proto")) {
            continue
        }
        const msgBuilder = new protobuf.Root()
        // console.log(path + '/' + name)
        msgBuilder.loadSync(path + '/' + name, { keepCase: true })
        for (let msgName in msgBuilder.nested.pb.nested) {
            msgBuilders[msgName] = msgBuilder.lookup(msgName)
            msgidToName[nametoid(msgName)] = msgName
        }
    }
}

initMsgBuilder("pb")

function runCli(context = {}, name = 'REPL') {
    const r = repl.start({
        // prompt: `${name} > `,
        preview: true,
        terminal: true,
    });
    Object.setPrototypeOf(r.context, context);
    global.console = r.context.console;
}

connect()
runCli({ send, connect })

function send(msgName = 'SayHelloReq', msgBody = { text: "hello, server!" }) {
    let b = encode(msgName, msgBody)
    socket.write(b)
}

async function recv() {
}
/**
 * encode
 * @param {string} msgName 
 * @param {object} msg 
 * @returns 
 */
function encode(msgName, msg) {
    let protoMsg = msgBuilders[msgName].create(msg)
    protoMsg = msgBuilders[msgName].encode(protoMsg).finish()
    let buf = Buffer.alloc(8)
    buf.writeUint32LE(protoMsg.length + 4)
    buf.writeUint32LE(nametoid(msgName), 4)
    return Buffer.concat([buf, Buffer.from(protoMsg)])
}

/**
 * decode
 * @param {Buffer} buf 
 * @returns 
 */
function decode(buf) {
    let msgid = buf.readUInt32LE()
    let msgName = msgidToName[msgid]
    const protoMsg = msgBuilders[msgName].decode(buf.slice(4))
    return [msgName, msgBuilders[msgName].toObject(protoMsg)]
}

function nametoid(msgName) {
    let s = 31
    let v = 0
    for (let c of msgName) {
        v = uint32(v * s) + c.charCodeAt()
    }
    return uint32(v)
}


// var m = encode("SayHelloReq", { text: "hello" })
// var n = decode(m)
// log(n)


function int(x) {
    x = Number(x);
    return x < 0 ? Math.ceil(x) : Math.floor(x);
}

function mod(a, b) {
    return a - Math.floor(a / b) * b;
}
function uint32(x) {
    return mod(int(x), Math.pow(2, 32));
}

function print(...any) {
    log(...any)
    process.stdout.write('> ') // 模拟prompt
}