package sync

import (
	"context"
	"sync"

	"github.com/alex-guoba/tbd/internal/provider"
	"github.com/alex-guoba/tbd/pkg/logger"
)

type SyncNotion struct {
	ctx   context.Context
	token string
	db    string
	wg    sync.WaitGroup
}

func NewNotion(ctx context.Context, token string, db string) *SyncNotion {
	return &SyncNotion{
		ctx:   ctx,
		token: token,
		db:    db,
	}
}

func (n *SyncNotion) startConsume(ch <-chan interface{}) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				// closed
				n.wg.Done()
				logger.Debug("notion sync stoped")
				return
			}
			chatMsg, ok := msg.(*provider.TopicChatMsg)
			if !ok {
				logger.Error("msg type error.")
			}
			logger.Debug("notion rsp:", chatMsg.Rsp)
			logger.Debug("notion req:", chatMsg.Req)

		case <-n.ctx.Done():
			n.wg.Done()
			logger.Debug("notion sync stoped")
			return
		}
	}
}

func (n *SyncNotion) Consume(ch <-chan interface{}) {
	n.wg.Add(1)
	go n.startConsume(ch)
}

func (n *SyncNotion) Close() {
	n.wg.Wait()
}
