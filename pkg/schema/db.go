package schema

import (
	"github.com/hashicorp/go-memdb"
)

const (
	TABLE_NAME = "event-results"
)

type DBClient interface {
	GetDBClient() (*memdb.MemDB, error)
}

type db struct{}

func NewDBClient() DBClient {
	return &db{}
}

func (d db) GetDBClient() (*memdb.MemDB, error) {
	db, err := memdb.NewMemDB(getDBSchema())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getDBSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			TABLE_NAME: {
				Name: TABLE_NAME,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:   "id",
						Unique: true,
						Indexer: &memdb.IntFieldIndex{
							Field: "ID",
						},
					},
					"event_id": {
						Name:   "event_id",
						Unique: false,
						Indexer: &memdb.StringFieldIndex{
							Field: "EventId",
						},
					},
					"e2e-latency": {
						Name:   "e2e-latency",
						Unique: false,
						Indexer: &memdb.IntFieldIndex{
							Field: "E2ELatency",
						},
					},
					"event-type": {
						Name:   "event-type",
						Unique: false,
						Indexer: &memdb.StringFieldIndex{
							Field: "EventType",
						},
					},
				},
			},
		},
	}
}
