package unix_test

import (
	"backend/constant"
	"backend/internal/unix"
	"net/http"
	"reflect"
	"testing"
)

var testUnixPduCases = []struct {
	name           string
	msgType        constant.UNIX_MSG_TYPE
	responseStatus int
	content        string
}{
	{
		name:    "Reset account",
		msgType: constant.UNIX_TYPE_RESET_ACCOUNT,
		content: "",
	},
	{
		name:           "Reset account response",
		msgType:        constant.UNIX_TYPE_RESET_ACCOUNT,
		responseStatus: http.StatusOK,
		content:        "",
	},
}

func TestUnixPdu(t *testing.T) {
	for _, tc := range testUnixPduCases {
		t.Run(tc.name, func(t *testing.T) {
			send := unix.NewUnixPdu(tc.msgType, tc.content)

			sendBytes, err := send.Marshal()
			if err != nil {
				t.Fatal(err.Error())
			}

			var recv unix.UnixPdu
			if err := recv.Unmarshal(sendBytes); err != nil {
				t.Fatal(err.Error())
			}

			if !reflect.DeepEqual(*send, recv) {
				t.Fatalf("send and recv is not the same, send: %v, recv: %v", send, recv)
			}
		})
	}
}
