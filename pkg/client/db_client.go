package client

import (
	"encoding/json"
	"github.com/anishj0shi/inmemorydb-service/pkg/schema"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/hashicorp/go-memdb"
	"log"
	"net/http"
	"sync/atomic"
)

type EventResultService interface {
	GetEventResult(w http.ResponseWriter, r *http.Request)
	PostEventResult(w http.ResponseWriter, r *http.Request)
}

func NewEventResultSrvice() EventResultService {
	memDb, err := schema.NewDBClient().GetDBClient()
	var count uint64 = 0
	if err != nil {
		panic(err)
	}
	return &eventResultRetriver{
		db:      memDb,
		counter: &count,
	}
}

type eventResultRetriver struct {
	db      *memdb.MemDB
	counter *uint64
}

func (e eventResultRetriver) GetEventResult(w http.ResponseWriter, r *http.Request) {
	topParameter, err := types.ParseInteger(r.URL.Query()["top"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	skipParameter, err := types.ParseInteger(r.URL.Query()["skip"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Printf("Query Parameter top: %d, skip: %d", topParameter, skipParameter)

	res, err := e.readEventResultFromDB(int(topParameter), int(skipParameter))
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Json marshaling error", http.StatusInternalServerError)
	}
	w.Write(jsonStr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (e eventResultRetriver) PostEventResult(w http.ResponseWriter, r *http.Request) {
	eventResult := &schema.EventResult{}
	atomic.AddUint64(e.counter, 1)

	log.Print("Write Request Received")
	err := json.NewDecoder(r.Body).Decode(eventResult)
	if err != nil {
		http.Error(w, "Bad Request Body", http.StatusBadRequest)
	}
	eventResult.ID = int(*e.counter)

	err = e.writeEventResultToDB(eventResult)
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (e eventResultRetriver) writeEventResultToDB(eventResult *schema.EventResult) error {
	txn := e.db.Txn(true)
	err := txn.Insert(schema.TABLE_NAME, eventResult)
	if err != nil {
		log.Print(err)
		return err
	}
	txn.Commit()
	return nil
}

func (e eventResultRetriver) readEventResultFromDB(top, skip int) ([]*schema.EventResult, error) {
	var result []*schema.EventResult

	txn := e.db.Txn(false)
	it, err := txn.LowerBound(schema.TABLE_NAME, "id", skip)
	if err != nil {
		return nil, err
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		e := obj.(*schema.EventResult)
		if e.ID > top {
			break
		}
		result = append(result, e)
	}
	return result, nil

}
