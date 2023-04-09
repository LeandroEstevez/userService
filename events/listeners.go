package events

import (
	"context"
	"encoding/json"
	"fmt"
	db "userMicroService/db/sqlc"
	"userMicroService/util"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Listen(message *kafka.Message) {
	key := string(message.Key)

	switch key {
	case "entry_created":
		EntryCreated(message.Value)
	case "entry_updated":
		EntryUpdated(message.Value)
	case "entry_deleted":
		EntryDeleted(message.Value)
	}
}

func EntryCreated(value []byte) {
	var entry db.Entry
	json.Unmarshal(value, &entry)

	err := util.Store.UpdateTotalTx(context.TODO(), entry)
	if err != nil {
		fmt.Println("Failed to update user")
		msg := DefaultMessage{
			Value: "N",
		}
		Produce("user_total_updated", msg)
	} else {
		fmt.Println("Sucees to update user")
		msg := DefaultMessage{
			Value: "Y",
		}
		Produce("user_total_updated", msg)
	}
}

func EntryUpdated(value []byte) {
	var msg db.UpdatedEntryMessage
	json.Unmarshal(value, &msg)

	user, err := util.Store.GetUserForUpdate(context.TODO(), msg.OrignalEntry.Owner)
	if err != nil {
		fmt.Println("There was an error", err)
		return
		// TODO: send back an event to entriesMicroService with the error
	}

	totalExpenseUpdate := user.TotalExpenses + (msg.UpdatedEntry.Amount - msg.OrignalEntry.Amount)

	params := db.UpdateUserParams{
		Username:      msg.OrignalEntry.Owner,
		TotalExpenses: totalExpenseUpdate,
	}
	_, err = util.Store.UpdateUser(context.TODO(), params)
	if err != nil {
		msg := DefaultMessage{
			Value: "N",
		}
		Produce("user_total_updated", msg)
	} else {
		msg := DefaultMessage{
			Value: "Y",
		}
		Produce("user_total_updated", msg)
	}
}

type EntrDeletedMessage struct {
	Owner  string
	Amount int64
}

func EntryDeleted(value []byte) {
	var msg EntrDeletedMessage
	json.Unmarshal(value, &msg)

	user, err := util.Store.GetUserForUpdate(context.TODO(), msg.Owner)
	if err != nil {
		fmt.Println("There was an error", err)
		return
		// TODO: send back an event to entriesMicroService with the error
	}

	totalExpenseUpdate := user.TotalExpenses - msg.Amount

	params := db.UpdateUserParams{
		Username:      msg.Owner,
		TotalExpenses: totalExpenseUpdate,
	}
	_, err = util.Store.UpdateUser(context.TODO(), params)
	if err != nil {
		msg := DefaultMessage{
			Value: "N",
		}
		Produce("user_total_updated", msg)
	} else {
		msg := DefaultMessage{
			Value: "Y",
		}
		Produce("user_total_updated", msg)
	}
}
