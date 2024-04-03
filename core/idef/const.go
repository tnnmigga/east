package idef

type ServerState int

const (
	ServerStateInit ServerState = iota // 初始化
	ServerStateRun                     // 运行
	ServerStateStop                    // 停止
	ServerStateExit                    // 进程退出
)

const (
	ConstKeyNone         = "none"
	ConstKeyNonuseStream = "nonuse-stream"
	ConstKeyOneOfMods    = "one-of-mods"
	ConstKeyServerID     = "server-id"
	ConstKeyExpires      = "expires"
)
