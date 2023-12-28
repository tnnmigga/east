package idef

type ServerState int

const (
	ServerStateInit  ServerState = iota // 初始化
	ServerStateRun                      // 运行
	ServerStateStop                     // 停止
	ServerStateClose                    // 进程退出
)
