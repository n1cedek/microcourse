package chat_v1

import (
	desc "microservices_course/week8/chat/pkg/chat_v1"
	"sync"
)

type Chat struct {
	streams map[string]desc.ChatV1_ConnectChatServer
	m       *sync.Mutex
}

type Implementation struct {
	desc.UnimplementedChatV1Server

	chats   map[string]*Chat
	mxChat sync.RWMutex

	channels map[string]chan *desc.Message
	mxChanel sync.RWMutex
}

func NewImplementation() *Implementation {
	return &Implementation{
		chats:     make(map[string]*Chat),
		channels: make(map[string]chan *desc.Message),
	}
}
