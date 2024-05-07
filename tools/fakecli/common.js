import protobuf from 'protobufjs';
import { readdirSync } from 'fs'
import { log } from 'console';
import repl from 'repl'


const msgBuilders = {}
const msgidToName = {}

export function initMsgBuilder(path) {
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

/**
 * encode
 * @param {string} msgName 
 * @param {object} msg 
 * @returns 
 */
export function encode(msgName, msg) {
    let protoMsg = msgBuilders[msgName].create(msg)
    protoMsg = msgBuilders[msgName].encode(protoMsg).finish()
    let buf = Buffer.alloc(4)
    buf.writeUint32LE(nametoid(msgName))
    return Buffer.concat([buf, Buffer.from(protoMsg)])
}

/**
 * decode
 * @param {Buffer} buf 
 * @returns 
 */
export function decode(buf) {
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


function uint32(x) {
    return x % Math.pow(2, 32);
}

export function print(...any) {
    log(...any)
    process.stdout.write('> ') // 模拟prompt
}

export function runCli(context = {}, name = 'REPL') {
    const r = repl.start({
        // prompt: `${name} > `,
        preview: true,
        terminal: true,
    });
    Object.setPrototypeOf(r.context, context);
    global.console = r.context.console;
}
// function int(x) {
//     x = Number(x);
//     return x < 0 ? Math.ceil(x) : Math.floor(x);
// }

// function mod(a, b) {
//     return a - Math.floor(a / b) * b;
// }


// var m = encode("SayHelloReq", { text: "hello" })
// var n = decode(m)
// log(n)