package tcp

import (
	"context"
	"encoding/binary"
	"encoding/json"
	errorss "errors"
	"io"

	"github.com/gucooing/spoon/errors"
	"github.com/gucooing/spoon/external"
)

const (
	headLenSize = 2
)

type (
	Read func(ctx context.Context, buf io.Reader) (external.Request, error)
)

// head 内置的tcp包头结构
type head struct {
	MsgID   uint32
	BodyLen uint32
}

func (h *head) GetBody() []byte {
	//TODO implement me
	panic("implement me")
}

func (h *head) GetMsgID() uint32 {
	return h.MsgID
}

type Request struct {
	*head
	body []byte
}

func (r *Request) GetBody() []byte {
	return r.body
}

// defaultRead 默认的读取客户端数据包方法
func defaultRead(ctx context.Context, buf io.Reader) (external.Request, error) {
	headLenBytes := make([]byte, headLenSize)
	_, err := io.ReadFull(buf, headLenBytes)
	if err != nil {
		if errorss.Is(err, io.EOF) {
			return nil, err
		}
		return nil, errors.New(errors.UnknownCode, "io.ReadFull", err.Error())
	}
	headLen := binary.BigEndian.Uint32(headLenBytes)

	headBytes := make([]byte, headLen)
	_, err = io.ReadFull(buf, headBytes)
	if err != nil {
		return nil, errors.New(errors.UnknownCode, "io.ReadFull", err.Error())
	}
	req := new(Request)
	if err = json.Unmarshal(headBytes, &req.head); err != nil {
		return nil, errors.New(errors.UnknownCode, "Head Decoder failed", err.Error())
	}

	bodyBytes := make([]byte, req.head.BodyLen)
	_, err = io.ReadFull(buf, bodyBytes)
	if err != nil {
		return nil, errors.New(errors.UnknownCode, "Body Decoder failed", err.Error())
	}
	req.body = bodyBytes
	return req, nil
}
