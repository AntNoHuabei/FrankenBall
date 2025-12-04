package chat

import (
	"github.com/AntNoHuabei/Remo/pkg/persist"
	"github.com/google/uuid"
	"github.com/ostafen/clover"
)

type Session struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func CreateSession() *Session {

	id := uuid.New().String()

	var session = &Session{
		Id:    id,
		Title: "New Session",
	}

	doc := clover.NewDocumentOf(session)
	doc.Set("_id", id)
	persist.DB.InsertOne(persist.Conversation, doc)

	return session

}

func DeleteSession(id string) error {

	return persist.DB.Query(persist.Conversation).DeleteById(id)
}

func SessionList(offset, limit int) ([]Session, error) {
	docs, err := persist.DB.Query(persist.Conversation).Skip(offset).Limit(limit).FindAll()
	if err != nil {
		return nil, err
	}

	var output = make([]Session, 0)
	for _, doc := range docs {
		var session = Session{}
		if err := doc.Unmarshal(&session); err == nil {

			output = append(output, session)
		}
	}
	return output, nil
}
