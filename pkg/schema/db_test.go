package schema

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
	"testing"
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
		err = txn.Insert(schema_name, d)
		if err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()

	txn = dbClient.Txn(false)
	it, err := txn.ReverseLowerBound(schema_name, "id", 3)
	if err != nil {
		t.Fatalf("error %+v", err)
	}
	//it, err := readTxn.ReverseLowerBound(schema_name, "id", 3)
	result := []*EventResult{}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		e := obj.(*EventResult)
		result = append(result, e)
	}

	if len(result) != 3 {
		t.Fatalf("Unexpected number of items retrieved. Total items: %d", len(result))
	}

}

func TestBla(t *testing.T) {
	type Person struct {
		Email string
		Name  string
		Age   int
	}

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"person": &memdb.TableSchema{
				Name: "person",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"age": &memdb.IndexSchema{
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := db.Txn(true)

	// Insert some people
	people := []*Person{
		&Person{"joe@aol.com", "Joe", 30},
		&Person{"lucy@aol.com", "Lucy", 35},
		&Person{"tariq@aol.com", "Tariq", 21},
		&Person{"dorothy@aol.com", "Dorothy", 53},
	}
	for _, p := range people {
		if err := txn.Insert("person", p); err != nil {
			panic(err)
		}
	}

	// Commit the transaction
	txn.Commit()

	// Create read-only transaction
	txn = db.Txn(false)
	defer txn.Abort()

	// Lookup by email
	raw, err := txn.First("person", "id", "joe@aol.com")
	if err != nil {
		panic(err)
	}

	// Say hi!
	fmt.Printf("Hello %s!\n", raw.(*Person).Name)

	// List all the people
	it, err := txn.Get("person", "id")
	if err != nil {
		panic(err)
	}

	fmt.Println("All the people:")
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Person)
		fmt.Printf("  %s\n", p.Name)
	}

	// Range scan over people with ages between 25 and 35 inclusive
	it, err = txn.LowerBound("person", "age", 25)
	if err != nil {
		panic(err)
	}

	fmt.Println("People aged 25 - 35:")
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Person)
		if p.Age > 35 {
			break
		}
		fmt.Printf("  %s is aged %d\n", p.Name, p.Age)
	}
}
