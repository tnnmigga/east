package idef

type HandlerFn struct {
	Cb  func(msg any)
	RPC func(msg any, resp func(any))
}

type Package struct {
	ServerID uint32
	Body     any
	TTL      int32
}

type BroadcastPackage struct {
	ServerType string
	Body       any
}

type RandomCastPackage struct {
	ServerType string
	Body       any
}

type RPCRequest struct {
	Module   IModule
	ServerID uint32
	Req      any
	Resp     any
	Cb       func(resp any, err error)
	Err      error
}

type RPCPackage struct {
	Req  any
	Resp chan any
	Err  chan error
}
