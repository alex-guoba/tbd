package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/alex-guoba/tbd/internal/entity"
	"github.com/alex-guoba/tbd/internal/provider"
	"github.com/alex-guoba/tbd/pkg/logger"
)

type SyncMemos struct {
	ctx   context.Context
	token string
	host  string
	wg    sync.WaitGroup

	// used to referelation for multi-round dialog
	mr *MemoRelaton
}

const RelationComment = "COMMENT"

type MemoRelaton struct {
	MemoID        int    `json:"memoID,omitempty"`
	RelatedMemoId int    `json:"relatedMemoId,omitempty"`
	Type          string `json:"type"`
}
type MemoReq struct {
	Relation        []MemoRelaton `json:"relationList,omitempty"`
	Content         string        `json:"content"`
	PresencePenalty float32       `json:"presence_penalty,omitempty"`
	Visibility      string        `json:"visibility,omitempty"`
	// CreatedTs       int           `json:"createdTs,omitempty"`
}

type MemoRsp struct {
	MemoId     int    `json:"id,omitempty"`
	Status     string `json:"rowStatus"`
	CreatorId  int    `json:"creatorId,omitempty"`
	Visibility string `json:"visibility,omitempty"`
	Content    string `json:"content"`
	ParentID   int    `json:"parentID,omitempty"`
}

func NewMemos(ctx context.Context, token string, host string) *SyncMemos {
	return &SyncMemos{
		ctx:   ctx,
		token: token,
		host:  strings.TrimSuffix(host, "/"),
	}
}

func (m *SyncMemos) formatContent(msg *provider.TopicChatMsg) string {
	content := ""
	if len(msg.Req.Messages) > 0 {
		last := msg.Req.Messages[len(msg.Req.Messages)-1]
		// content += "#TBT#" +  + ": " + last.Content + "\n"
		content = fmt.Sprintf("#TBT# [%s]: %s\n", last.Role, last.Content)
	}
	content = content + fmt.Sprintf("[%s]: ", entity.MessageRoleAssistant) + msg.Rsp.Result
	return content
}

func (m *SyncMemos) postCreate(msg *provider.TopicChatMsg) error {
	url := fmt.Sprintf("%s/api/v1/memo", m.host)

	content := m.formatContent(msg)
	if content == "" {
		return fmt.Errorf("empty msg")
	}

	mq := m.newMemosReq(content)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	data, _ := json.Marshal(mq)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+m.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("memos post eror: %s", err.Error())
		return err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	logger.Debugf("memos reply: %s", body)

	var cRsp MemoRsp
	err = json.Unmarshal(body, &cRsp)
	if err != nil {
		logger.Errorf("memos post eror: %s", err.Error())
		return err
	}

	m.storeRelation(&cRsp)

	return nil
}

func (n *SyncMemos) newMemosReq(content string) *MemoReq {
	req := &MemoReq{
		Content: content,
	}
	if n.mr != nil {
		//set comment relation if not root
		req.Relation = append(req.Relation, *n.mr)
	}
	return req
}

// store root memos ID for multi-round dialog
func (n *SyncMemos) storeRelation(rsp *MemoRsp) {
	if n.mr == nil {
		n.mr = &MemoRelaton{
			RelatedMemoId: rsp.MemoId,
			Type:          RelationComment,
		}
	}
}

func (n *SyncMemos) startConsume(ch <-chan interface{}) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				// closed
				n.wg.Done()
				logger.Debug("memos sync stoped")
				return
			}
			chatMsg, ok := msg.(*provider.TopicChatMsg)
			if !ok {
				logger.Error("msg type error.")
			}
			if err := n.postCreate(chatMsg); err != nil {
				logger.Error("msg sync error.", err)
			}
			// logger.Debug("memos rsp:", chatMsg.Rsp)
			// logger.Debug("memos req:", chatMsg.Req)

		case <-n.ctx.Done():
			n.wg.Done()
			logger.Debug("memos sync stoped")
			return
		}
	}
}

func (n *SyncMemos) Consume(ch <-chan interface{}) {
	n.wg.Add(1)
	go n.startConsume(ch)
}

func (n *SyncMemos) Close() {
	n.wg.Wait()
}
