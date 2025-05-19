package chat_v1

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "microservices_course/week8/chat/pkg/chat_v1"

)

func (i *Implementation) CreateChat(ctx context.Context,_ *emptypb.Empty)(*desc.CreateChatResponse,error){
	chatId,err:=uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	i.channels[chatId.String()] = make(chan *desc.Message)

	return &desc.CreateChatResponse{ChatId: chatId.String()},nil
}
