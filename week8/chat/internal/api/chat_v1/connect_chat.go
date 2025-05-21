package chat_v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	desc "microservices_course/week8/chat/pkg/chat_v1"
	"sync"
)

func (i *Implementation) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer) error {
	i.mxChanel.RLock()
	chatChan, ok := i.channels[req.GetChatId()]
	i.mxChanel.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	i.mxChat.Lock()
	if _, okChat := i.chats[req.GetChatId()]; !okChat {
		i.chats[req.GetChatId()] = &Chat{
			streams: make(map[string]desc.ChatV1_ConnectChatServer),
			m:       new(sync.Mutex),
		}
	}
	i.mxChat.Unlock()

	i.chats[req.GetChatId()].m.Lock()
	i.chats[req.GetChatId()].streams[req.GetUsername()] = stream
	i.chats[req.GetChatId()].m.Unlock()

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}
			for _, st := range i.chats[req.GetChatId()].streams {
				if err := st.Send(msg); err != nil {
					return err
				}
			}
		case <-stream.Context().Done():
			i.chats[req.GetChatId()].m.Lock()
			delete(i.chats[req.GetChatId()].streams, req.GetUsername())
			i.chats[req.GetChatId()].m.Unlock()
			return nil

		}
	}
}
