package main

import "github.com/sonyarouje/simdb/db"

func InitializeDb() (err error, driver *db.Driver) {
	driver, err = db.New("data")
	if err != nil {
		panic(err)
	}
	return err, driver
}

// ---------------------------------------------------------------------------
// Sequence for id field - you call this by NextId
type Counter struct {
	Table  string `json:"table"`
	Number int    `json:"number"`
}

func (c Counter) ID() (jsonField string, value interface{}) {
	value = c.Table
	jsonField = "table"
	return
}

// Sample Code:
// fmt.Println(NextId("product"))
func NextId(table string) (id int) {
	//AsEntity takes a pointer to Counter variable (not an array pointer)
	var counter Counter

	err := Database.Open(Counter{}).Where("table", "=", table).First().AsEntity(&counter)

	if err != nil {
		err = Database.Insert(Counter{Table: table, Number: 1})
		if err != nil {
			panic(err)
		}
		return 1
	}
	counter.Number++
	err = Database.Update(counter)
	if err != nil {
		panic(err)
	}
	return counter.Number

}

var Database *db.Driver

func init() {
	var err error
	err, Database = InitializeDb()

	if err != nil {
		panic(err)
	}
}
