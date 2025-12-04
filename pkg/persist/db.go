package persist

import "github.com/ostafen/clover"

const Conversation = "conversation"
const Message = "message"
const SessionCheckpoint = "session_checkpoint"

var DB *clover.DB

func InitDB() error {

	db, err := clover.Open("clover.db")
	if err != nil {
		return err
	}

	DB = db
	if h, _ := db.HasCollection(Conversation); !h {
		err = db.CreateCollection(Conversation)
		if err != nil {
			return err
		}
	}
	if h, _ := db.HasCollection(Message); !h {
		err = db.CreateCollection(Message)
		if err != nil {
			return err
		}
	}
	if h, _ := db.HasCollection(SessionCheckpoint); !h {
		err = db.CreateCollection(SessionCheckpoint)
		if err != nil {
			return err
		}
	}
	return nil
}
