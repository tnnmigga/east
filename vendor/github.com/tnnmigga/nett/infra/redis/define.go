package redis

import "time"

// redis执行单条命令
// 通过msgbus.RPC调用
// 返回为具体的结果只可能为string或int64
type Exec struct {
	Cmd     []any
	Key     string        // 默认为Cmd[1]
	Timeout time.Duration // 默认为3s
}

// redis执行多条命令批处理
// 通过msgbus.RPC调用
// 返回为具体的结果为[]any
type ExecMulti struct {
	Cmds    [][]any
	Key     string        // 默认为Cmds[0][1]
	Timeout time.Duration // 默认为3s
}
