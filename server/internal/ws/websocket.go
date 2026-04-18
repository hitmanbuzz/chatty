package ws

import (
	"log/slog"
	"msg_app/internal/message"
	"msg_app/internal/util"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

type Websocket struct {
	logger *slog.Logger
}

func Init(logger *slog.Logger) *Websocket {
	return &Websocket{
		logger: logger,
	}
}

func (w *Websocket) Messaging(g *gin.Context) {
	conn, err := websocket.Accept(g.Writer, g.Request, nil)
	if err != nil {
		w.logger.Error("failed to accept websocket", "error", err)
		return
	}

	defer conn.CloseNow()

	conn.SetReadLimit(util.MAX_MSG_BYTE)
	ctx := g.Request.Context()

	for {
		msg := message.Init(w.logger)
		content, isExit := msg.ReadMessage(ctx, conn)

		if isExit {
			break
		}

		w.logger.Info("message received", "content", content)
	}

	conn.Close(websocket.StatusNormalClosure, "server closing connection")
}

func (w *Websocket) LoginUser(g *gin.Context) {
	// conn, err := websocket.Accept(g.Writer, g.Request, nil)
}
