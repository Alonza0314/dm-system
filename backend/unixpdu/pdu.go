package unixpdu

import (
	"backend/constant"
	"encoding/json"
)

type UnixPdu struct {
	MsgType        constant.UNIX_MSG_TYPE `json:"msg_type" binding:"required"`
	ResponseStatus int                    `json:"response_status"`
	Content        string                 `json:"content"`
}

func NewUnixPdu(unixType constant.UNIX_MSG_TYPE, responseStatus int, content string) *UnixPdu {
	return &UnixPdu{
		MsgType:        unixType,
		ResponseStatus: responseStatus,
		Content:        content,
	}
}

func (u *UnixPdu) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (u *UnixPdu) Unmarshal(input []byte) error {
	return json.Unmarshal(input, u)
}
