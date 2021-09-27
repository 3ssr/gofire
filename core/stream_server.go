package core

import (
	"context"
	"gofire/iface"
	"log"
)

type ServerStream struct {
	conn       iface.IConn
	server     *FireServer
	ctx        context.Context
	cancel     context.CancelFunc
	msgChannel chan iface.IMsg
}

func NewServerStream(conn iface.IConn, server *FireServer) iface.IStream {
	ctx, cancel := context.WithCancel(context.Background())
	s := &ServerStream{
		conn:       conn,
		server:     server,
		ctx:        ctx,
		cancel:     cancel,
		msgChannel: make(chan iface.IMsg, 2),
	}
	return s
}

func (s *ServerStream) Flow() {
	go s.ReadLoop()
	go s.WriteLoop()
}

func (s *ServerStream) Close() {
	s.conn.Close()
	close(s.msgChannel)
}

func (s *ServerStream) ReadLoop() {
	// defer s.Close()
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			data, err := s.server.pcodec.Decode(s.conn)
			if err != nil {
				log.Println("pcodec decode error", err)
				return
			}

			msg, err := s.server.mcodec.Decode(data)
			if err != nil {
				log.Println("mcodec decode error", err)
				return
			}

			req := iface.Request{
				Stream: s,
				Msg:    msg,
			}

			log.Printf("server receive msg: %+v\n", msg)
			handler := s.server.GetHandler(msg.GetAction())
			if handler == nil {
				log.Println("not support action for action:", msg.GetAction())
				return
			}

			go handler.Do(req)
		}
	}
}

func (s *ServerStream) WriteLoop() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.msgChannel:
			log.Println("WriteLoop 1")
			msgData, err := s.server.mcodec.Encode(msg)
			if err != nil {
				log.Println("mcodec encode msg error", err)
				continue
			}

			log.Println("WriteLoop 2")

			err = s.server.pcodec.Encode(msgData, s.conn)
			if err != nil {
				log.Println("pcodec encode msg error", err)
				continue
			}
			log.Println("WriteLoop 3")
		}
	}
}

func (s *ServerStream) Write(msg iface.IMsg) {
	s.msgChannel <- msg
	log.Println("server write msg to channel")
}