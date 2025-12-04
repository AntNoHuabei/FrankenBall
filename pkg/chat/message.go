package chat

import (
	"github.com/AntNoHuabei/Remo/pkg/persist"
	"github.com/google/uuid"
	"github.com/ostafen/clover"
)

type Message struct {
	Model       string `json:"model"`
	CreatedTime int64  `json:"created_time"`
	Session     string `json:"session"`
	Id          string `json:"id"`
	Content     string `json:"content"`
	Role        string `json:"role"`
	RequestId   string `json:"request_id"`
}

func Messages(session string) ([]*Message, error) {

	docs, err := persist.DB.Query(persist.Message).Where(clover.Field("session").Eq(session)).FindAll()

	if err != nil {
		return nil, err
	}
	var output = make([]*Message, 0)
	for _, doc := range docs {
		var message = Message{}
		if err := doc.Unmarshal(&message); err == nil {
			output = append(output, &message)
		}
	}
	return output, nil
}

func MessageAppend(session string, message *Message) error {

	if message.Id == "" {
		message.Id = uuid.New().String()
	}

	doc := clover.NewDocumentOf(message)
	doc.Set("_id", message.Id)

	persist.DB.InsertOne(persist.Message, doc)
	return nil
}
