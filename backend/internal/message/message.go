package message

import (
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"msg_app/internal/util"
	"net"
	"unicode/utf8"
)

type Message struct {
	logger *slog.Logger
}

func NewMsg(logger *slog.Logger) *Message {
	return &Message{
		logger: logger,
	}
}

func (m *Message) HandleMsg(conn *net.Conn) (string, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(*conn, header); err != nil {
		return "", fmt.Errorf("client disconnected or error reading header: %v\n", err)
	}

	msgLen := binary.BigEndian.Uint32(header)

	if msgLen > util.MAX_MSG_BYTE {
		return "", fmt.Errorf("message payload to large, %d\n", msgLen)
	}

	payload := make([]byte, msgLen)
	if _, err := io.ReadFull(*conn, payload); err != nil {
		return "", fmt.Errorf("failed to read payload: %v\n", err)
	}

	msg := string(payload)
	charCount := utf8.RuneCountInString(msg)

	if charCount > util.MAX_MSG_LEN {
		return "", fmt.Errorf("message rejected, exceeded max characters\n")
	}

	return msg, nil
}
