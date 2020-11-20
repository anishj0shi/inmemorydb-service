package schema

import (
	"fmt"
	"testing"
	"time"
)

func TestInMemDBOperations(t *testing.T) {
	data := []*EventResult{
		{
			ID:         1,
			EventId:    "a",
			E2ELatency: 400,
			EventType:  "order.created",
		},
		{
			ID:         2,
			EventId:    "b",
			E2ELatency: 400,
			EventType:  "order.created",
		},
		{
			ID:         3,
			EventId:    "c",
			E2ELatency: 400,
			EventType:  "order.created",
		},
		{
			ID:         4,
			EventId:    "d",
			E2ELatency: 400,
			EventType:  "order.created",
		},
		{
			ID:         5,
			EventId:    "e",
			E2ELatency: 400,
			EventType:  "order.created",
		}, {
			ID:         5,
			EventId:    "f",
			E2ELatency: 400,
			EventType:  "order.created",
		},
	}
	client := NewDBClient()

	dbClient, err := client.GetDBClient()
	if err != nil {
		t.Fatal(err)
	}

	txn := dbClient.Txn(true)

	for _, d := range data {
		err = txn.Insert(TABLE_NAME, d)
		if err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()

	txn = dbClient.Txn(false)
	it, err := txn.ReverseLowerBound(TABLE_NAME, "id", 3)
	if err != nil {
		t.Fatalf("error %+v", err)
	}
	//it, err := readTxn.ReverseLowerBound(TABLE_NAME, "id", 3)
	result := []*EventResult{}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		e := obj.(*EventResult)
		result = append(result, e)
	}

	if len(result) != 3 {
		t.Fatalf("Unexpected number of items retrieved. Total items: %d", len(result))
	}

}

func TestTime(t *testing.T) {
	receivedTime := time.Unix(0, 1605892478687667000)
	currentTime := time.Now().UTC()

	fmt.Println(currentTime.Sub(receivedTime).Seconds())
}