package idef

type ServerState int

const (
	ServerStateNone  ServerState = iota // None
	ServerStateInit                     // 初始化
	ServerStateRun                      // 运行
	ServerStateStop                     // 停止
	ServerStateClose                    // 进程退出
	ServerStateMax                      // END
)
