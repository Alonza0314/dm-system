package unix

import (
	"backend/constant"
	"backend/unixpdu"
	"fmt"
	"net"
	"net/http"
)

func (u *UnixServer) handResetAccount(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			u.UnxLog.Errorf("failed to close unix conn after handler: %v", err)
		}
	}()

	response := unixpdu.NewUnixPdu(constant.UNIX_TYPE_RESET_ACCOUNT, http.StatusOK, "")

	if err := u.DmContext.Db().RemoveAll(constant.COLL_ACCOUNT); err != nil {
		u.UnxLog.Errorf("failed to remove exist account: %v", err)
		response.ResponseStatus, response.Content = http.StatusInternalServerError, fmt.Sprintf("failed to remove exist account: %v", err)
	}

	responseBytes, err := response.Marshal()
	if err != nil {
		u.UnxLog.Errorf("failed to mashal unix PDU: %v", err)
	}

	n, err := conn.Write(responseBytes)
	if err != nil {
		u.UnxLog.Errorf("failed to write back to unix conn: %v", err)
		return
	}
	u.UnxLog.Debugf("Write %d bytes to unix conn", n)
	u.UnxLog.Tracef("Unix write: %v", responseBytes)
}
