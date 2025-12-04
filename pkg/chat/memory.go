package chat

import (
	"context"
	"fmt"
	"github.com/AntNoHuabei/Remo/pkg/persist"
	"github.com/cloudwego/eino/compose"
	"github.com/ostafen/clover"
)

func NewStore(session string) (compose.CheckPointStore, error) {
	doc, err := persist.DB.Query(persist.SessionCheckpoint).FindById(session)
	if err != nil {
		doc = clover.NewDocument()
		doc.Set("_id", session)
		_, err = persist.DB.InsertOne(persist.SessionCheckpoint, doc)
		if err != nil {
			return nil, err
		}
	}
	return &sessionStore{
		session: session,
		doc:     doc,
	}, nil
}

type sessionStore struct {
	session string
	doc     *clover.Document
}

func (i *sessionStore) Set(ctx context.Context, key string, value []byte) error {

	if i.doc != nil {
		i.doc.Set(key, value)
	} else {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (i *sessionStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	if i.doc != nil {
		if i.doc.Has(key) {
			v := i.doc.Get(key)
			switch v.(type) {
			case []byte:
				return v.([]byte), true, nil

			default:
				return nil, false, fmt.Errorf("invalid type")
			}
		} else {
			return nil, false, fmt.Errorf("key not found")
		}
	} else {
		return nil, false, fmt.Errorf("session not found")
	}
}
