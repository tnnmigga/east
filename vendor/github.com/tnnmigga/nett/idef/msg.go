package idef

type Handler struct {
	Cb  func(msg any)
	RPC func(req any, resolve func(any), reject func(error))
}

type CastPackage struct {
	ServerID uint32
	Body     any
}

type StreamCastPackage struct {
	ServerID uint32
	Body     any
	Header   map[string]string
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
	Req  any
	Resp chan any
	Err  chan error
}

type RPCResponse struct {
	Module IModule
	Req    any
	Resp   any
	Err    error
	Cb     func(resp any, err error)
}

type RPCContext struct {
	Caller   IModule
	ServerID uint32
	Req      any
	Resp     any
	Cb       func(resp any, err error)
}