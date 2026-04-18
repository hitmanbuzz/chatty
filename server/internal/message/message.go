package message

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"msg_app/internal/util"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Message struct {
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}

func new_msg(user_id uint, content string) Message {
	return Message{
		UserID:  user_id,
		Content: content,
	}
}

type MessageHandler struct {
	logger *slog.Logger
	msg    Message
}

func Init(logger *slog.Logger) *MessageHandler {
	return &MessageHandler{
		logger: logger,
	}
}

func (m *MessageHandler) ReadMessage(ctx context.Context, conn *websocket.Conn) (string, bool) {
	var v Message
	err := wsjson.Read(ctx, conn, &v)
	if err != nil {
		msg_type := handleMsgType(err)
		err_msg, level, isExit, err := msg_type_string(msg_type, err)
		if err != nil && level == util.ERROR {
			m.logger.Error(err_msg)
		} else {
			switch level {
			case util.INFO:
				m.logger.Info(err_msg)
			case util.WARN:
				m.logger.Warn(err_msg)
			}
		}

		return "", isExit
	}

	isValid := isValidMsg(&v.Content)
	if !isValid {
		m.logger.Warn("message length exceeded the limit", "length", len(v.Content))
		return "", false
	}

	return v.Content, false
}

func isValidMsg(content *string) bool {
	if len(*content) > util.MAX_MSG_LEN {
		return false
	}
	return true
}

// The below 2 functions could have been merged since there are some duplications but it works so I will ignore it for now

func handleMsgType(e error) util.MessageType {
	status := websocket.CloseStatus(e)
	switch status {
	case websocket.StatusNormalClosure, websocket.StatusGoingAway:
		return util.DISCONNECT
	case websocket.StatusAbnormalClosure:
		return util.ABNORMAL_DISCONNECT
	case websocket.StatusMessageTooBig:
		return util.EXCEED_MAX_MSG_BYTE
	case -1:
		if errors.Is(e, io.EOF) || errors.Is(e, context.Canceled) {
			return util.ABRUPT_DISCONNECT
		}
		return util.FAIL_JSON_PARSE
	}
	return util.NIL
}

func msg_type_string(msg_type util.MessageType, e error) (string, util.LogLevel, bool, error) {
	switch msg_type {
	case util.DISCONNECT:
		return "user disconnected gracefully", util.INFO, true, nil
	case util.ABNORMAL_DISCONNECT:
		return "user disconnected abnormally", util.WARN, true, nil
	case util.ABRUPT_DISCONNECT:
		return "user disconnected abruptly", util.WARN, true, nil
	case util.EXCEED_MAX_MSG_BYTE:
		return "message exceed the bytes limit", util.WARN, false, nil
	case util.FAIL_JSON_PARSE:
		return "failed to parse the content to json", util.ERROR, false, fmt.Errorf("failed to parse the content to json")
	default:
		return fmt.Sprintf("unkown error: %v", e), util.ERROR, true, fmt.Errorf("unkown error: %w", e)
	}
}
