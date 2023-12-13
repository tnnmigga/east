import protobuf from 'protobufjs';
import repl from 'repl'
import { readdirSync } from 'fs'
import { log } from 'console';
import { Socket } from 'net';

const msg_builders = {}
const msgid_to_name = {}
var socket = new Socket()

function connect() {
    socket = new Socket()
    socket.connect('9527', '127.0.0.1', function () {
        print("connect success")
        send("SayHelloReq", { text: "hello, server!" })
    })
    socket.on("data", function (data) {
        try {
            let [msg_name, msg] = decode(data)
            print("recv server msg:", msg_name, msg)
        } catch {
            print("decode error", data, data.length)
        }
    })
}

function init_msg_builder(path) {
    let files = readdirSync(path)
    for (let name of files) {
        if (!name.endsWith(".proto")) {
            continue
        }
        const msg_builder = new protobuf.Root()
        // console.log(path + '/' + name)
        msg_builder.loadSync(path + '/' + name, { keepCase: true })
        for (let msg_name in msg_builder.nested.pb.nested) {
            msg_builders[msg_name] = msg_builder.lookup(msg_name)
            msgid_to_name[nametoid(msg_name)] = msg_name
        }
    }
    // log(msgid_to_name)
}

init_msg_builder("pb")

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

function send(msg_name = 'SayHelloReq', msg_body = { text: "hello, server!" }) {
    let b = encode(msg_name, msg_body)
    socket.write(b)
}

async function recv() {
}
/**
 * encode
 * @param {string} msg_name 
 * @param {object} msg 
 * @returns 
 */
function encode(msg_name, msg) {
    let proto_msg = msg_builders[msg_name].create(msg)
    proto_msg = msg_builders[msg_name].encode(proto_msg).finish()
    let buf = Buffer.alloc(8)
    buf.writeUint32LE(proto_msg.length + 4)
    buf.writeUint32LE(nametoid(msg_name), 4)
    return Buffer.concat([buf, Buffer.from(proto_msg)])
}

/**
 * decode
 * @param {Buffer} buf 
 * @returns 
 */
function decode(buf) {
    let msgid = buf.readUInt32LE()
    let msg_name = msgid_to_name[msgid]
    const proto_msg = msg_builders[msg_name].decode(buf.slice(4))
    return [msg_name, msg_builders[msg_name].toObject(proto_msg)]
}

function nametoid(msg_name) {
    let s = 31
    let v = 0
    for (let c of msg_name) {
        v = uint32(v * s) + c.charCodeAt()
    }
    return uint32(v)
}

// var msg_id = nametoid("SayHelloReq")

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