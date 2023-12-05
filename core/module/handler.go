package module

import (
	"east/core/idef"
	"reflect"
)

var (
	rpcPackage = reflect.TypeOf((*idef.RPCPackage)(nil))
	rpcRequest = reflect.TypeOf((*idef.RPCRequest)(nil))
)
