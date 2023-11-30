package srv

import (
	"miniEdward/bot"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Bot    *bot.Bot
	engine *gin.Engine
}

func (s *Server) Start() {
	s.engine.Run()
}

func NewSrv(b *bot.Bot) *Server {
	r := gin.Default()

	replier := r.Group("replier/")

	replier.PUT("set", b.SetRead)
	replier.PUT("unset", b.UnsetRead)

	return &Server{
		Bot:    b,
		engine: r,
	}
}
