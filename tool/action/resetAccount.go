package action

import (
	"backend/constant"
	"backend/unixpdu"
	"dmtool/unix"
	"fmt"
	"net/http"
)

func HandleResetAccount(unixClient *unix.UnixClient) {
	request := unixpdu.NewUnixPdu(constant.UNIX_TYPE_RESET_ACCOUNT, -1, "")

	requestBytes, err := request.Marshal()
	if err != nil {
		fmt.Printf("failed to matshal reset account request: %v\n", err)
		return
	}

	if _, err := unixClient.Write(requestBytes); err != nil {
		fmt.Printf("failed to write request to conn: %v\n", err)
		return
	}

	buffer := make([]byte, constant.UNIX_BUFFER_SIZE)
	n, err := unixClient.Read(buffer)
	if err != nil {
		fmt.Printf("failed to read response from conn: %v\n", err)
		return
	}

	var response unixpdu.UnixPdu
	if err := response.Unmarshal(buffer[:n]); err != nil {
		fmt.Printf("failed to unmarshal response: %v\n", err)
		return
	}

	if response.MsgType != constant.UNIX_TYPE_RESET_ACCOUNT {
		fmt.Printf("incorrect unix message type: %d\n", response.MsgType)
		return
	}

	if response.ResponseStatus != http.StatusOK {
		fmt.Printf("failed to reset account: %s\n", response.Content)
	}

	fmt.Println("Reset account success")
}
