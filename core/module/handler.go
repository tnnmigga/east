package module

import (
	"east/core/define"
	"reflect"
)

var (
	rpcPackage = reflect.TypeOf((*define.RPCPackage)(nil))
	rpcRequest = reflect.TypeOf((*define.RPCRequest)(nil))
)
