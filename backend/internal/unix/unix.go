package unix

import (
	"backend/constant"
	bctx "backend/internal/context"
	"backend/logger"
	"errors"
	"net"
	"os"
	"time"
)

type UnixServerParams struct {
	Enable bool
	Path   string

	*bctx.DmContext

	*logger.BackendLogger
}

type UnixServer struct {
	enable bool
	path   string

	net.Listener

	*bctx.DmContext

	*logger.BackendLogger
}

func NewUnixServer(params *UnixServerParams) *UnixServer {
	if err := os.Remove(params.Path); err != nil && !errors.Is(err, os.ErrNotExist) {
		params.UnxLog.Errorf("failed to remove old unix socket file at %s: %v", params.Path, err)
		return nil
	}

	listener, err := net.Listen("unix", params.Path)
	if err != nil {
		params.UnxLog.Errorf("failed to init unix listener at %s: %v", params.Path, err)
		return nil
	}
	params.UnxLog.Debugf("unix listener init at %s", params.Path)

	return &UnixServer{
		enable: params.Enable,
		path:   params.Path,

		Listener: listener,

		DmContext: params.DmContext,

		BackendLogger: params.BackendLogger,
	}
}

func (u *UnixServer) Start() {
	u.UnxLog.Infoln("Starting unix server...")

	if !u.enable {
		u.UnxLog.Warnln("Unix server is disabled")
	}

	go u.run()
	time.Sleep(300 * time.Millisecond)

	u.UnxLog.Infoln("Unix server started")
}

func (u *UnixServer) Stop() {
	u.UnxLog.Infof("Stopping unix server...")

	if err := u.Listener.Close(); err != nil {
		u.UnxLog.Errorf("failed to close unix listener: %v", err)
		return
	}

	u.UnxLog.Infof("Unix server stopped successfully")
}

func (u *UnixServer) run() {
	for {
		conn, err := u.Listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			u.UnxLog.Errorf("failed to accept unix connection: %v", err)
			continue
		}

		go u.dispatcher(conn)
	}
}

func (u *UnixServer) dispatcher(conn net.Conn) {
	buffer := make([]byte, constant.UNIX_BUFFER_SIZE)

	n, err := conn.Read(buffer)
	if err != nil {
		u.UnxLog.Errorf("failed to read from unix conn: %v", err)

		if err := conn.Close(); err != nil {
			u.UnxLog.Errorf("failed to close unix conn while read error happened: %v", err)
		}

		return
	}
	u.UnxLog.Debugf("Read %d bytes from unix conn", n)
	u.UnxLog.Tracef("Unix read: %v", buffer[:n])

	var unixPdu UnixPdu
	if err := unixPdu.Unmarshal(buffer[:n]); err != nil {
		u.UnxLog.Errorf("failed to umarshl unix buffer: %v", err)

		if err := conn.Close(); err != nil {
			u.UnxLog.Errorf("failed to close unix conn while unmarshal error happened: %v", err)
		}

		return
	}

	switch unixPdu.MsgType {
	case constant.UNIX_TYPE_RESET_ACCOUNT:
		u.handResetAccount(conn)
	default:
		u.UnxLog.Errorf("unknown unix PDU type %d", unixPdu.MsgType)

		if err := conn.Close(); err != nil {
			u.UnxLog.Errorf("failed to close unix conn with unknown PDU type: %v", err)
		}
	}
}
