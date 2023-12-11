package cgrpc

import (
	"github.com/slclub/easy/vendors/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

/**

suitable for game server
*/

type GameGrpcConn struct {
	*grpc.ClientConn
	Info ServInfo
}

type ServInfo struct {
	AddrToRpc    string
	AddrToClient string
	Name         string
	data         any
}

func (self *ServInfo) Set(data any) {
	self.data = data
}

func (self *ServInfo) Value() any {
	return self.data
}

// game client connn

func newGameRpcConn(assignment option.Assignment) *GameGrpcConn {
	ggc := &GameGrpcConn{}
	if assignment == nil {
		return ggc
	}
	assignment.Target(&ggc.Info)
	assignment.Apply()
	return ggc
}

func (self *GameGrpcConn) Dial(target string, opts ...grpc.DialOption) error {
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return err
	}
	self.ClientConn = conn
	self.Info.AddrToRpc = target
	return nil
}

func (self *GameGrpcConn) Close() {
	self.ClientConn.Close()
}

func (self *GameGrpcConn) GetState() connectivity.State {
	return self.ClientConn.GetState()
}

func NewServInfo(infos []string) *ServInfo {
	s := &ServInfo{}
	switch len(infos) {
	case 0:
		return nil
	case 1:
		s.AddrToRpc = infos[0]
	case 2:
		s.AddrToRpc = infos[0]
		s.AddrToClient = infos[1]
	case 3:
		s.AddrToRpc = infos[0]
		s.AddrToClient = infos[1]
		s.Name = infos[2]
	case 4:
		s.AddrToRpc = infos[0]
		s.AddrToClient = infos[1]
		s.Name = infos[2]
	}
	return s
}
