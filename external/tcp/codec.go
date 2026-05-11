package tcp

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	errorss "errors"
	"github.com/gucooing/spoon/errors"
	"github.com/gucooing/spoon/external"
	"io"
)

const (
	headLenSize = 2
	maxHeadSize = 1 << 10
	maxBodySize = 1 << 24
)

type (
	Read  func(ctx context.Context, buf io.Reader) (external.Request, error)
	Write func(ctx context.Context, buf io.Writer, rsp external.Response) error
)

// head 内置的tcp包头结构
type head struct {
	MsgID   uint32 `json:"msg_id"`
	BodyLen uint32 `json:"body_len"`
	Crc32   uint32 `json:"crc32"`
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
	headLen := binary.BigEndian.Uint16(headLenBytes)
	if headLen > maxHeadSize {
		return nil, errors.New(errors.UnknownCode, "head size", "head size long")
	}
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

func defaultWrite(ctx context.Context, buf io.Writer, rsp external.Response) error {
	headBytes, err := json.Marshal(rsp)
	if err != nil {
		return errors.New(errors.UnknownCode, "Head Encoder failed", err.Error())
	}
	var data bytes.Buffer
	headLen := len(headBytes)
	if headLen > maxHeadSize {
		return errors.New(errors.UnknownCode, "head size", "head size long")
	}
	if err = binary.Write(&data, binary.BigEndian, uint16(headLen)); err != nil {
		return errors.New(errors.UnknownCode, "Head Encoder failed", err.Error())
	}
	if _, err = data.Write(headBytes); err != nil {
		return errors.New(errors.UnknownCode, "Head Write failed", err.Error())
	}
	if _, err = data.Write(rsp.GetBody()); err != nil {
		return errors.New(errors.UnknownCode, "Body Write failed", err.Error())
	}
	n, err := buf.Write(data.Bytes())
	if err != nil {
		return errors.New(errors.UnknownCode, "write rsp", err.Error())
	}
	if n != data.Len() {
		return errors.New(errors.UnknownCode, "send rsp", "send truncated")
	}

	return nil
}
