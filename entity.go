// AccountManager project main.go
package main

import (
	"time"
)

// Entity relates to an 'Object' or struct
type Entity struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Fields     []Field `json:"fields"`
	EntityType int     `json:"type,omitempty"` // 0..Normal, 1..Lookup
}

// Field is each and every single attribute.
// Object is empty except in case type=slicetype keeps the name of the Object
type Field struct {
	Name      string `json:"name"`
	Type      string `json:"type"` // string, int, bool, lookup
	Object    string `json:"object,omitempty"`
	Maxlength int    `json:"maxlength,omitempty"`
	Size      int    `json:"size,omitempty"`
	Required  bool   `json:"required"`
	Step      int    `json:"step,omitempty"` //for Number fields
	Min       int    `json:"min,omitempty"`  //for Number fields
	Max       int    `json:"max,omitempty"`  //for Number fields
}

func NewEntity() (e *Entity) {
	id := NextId("Entity")
	e = &Entity{Id: id}
	return e
}

func (e *Entity) addField(f Field) {
	e.Fields = append(e.Fields, f)
}

func (e Entity) ID() (jsonField string, value interface{}) {
	value = e.Id
	jsonField = "id"
	return
}

func (e Entity) TimeStamp() string {
	return time.Now().Format(time.UnixDate)
}

// Database access functions
func getAllEntities() (err error, entities []Entity) {
	err = Database.Open(Entity{}).Where("id", ">", 0).Get().AsEntity(&entities)
	if err != nil {
		panic(err)
	}
	return err, entities
}

func getEntityById(id int) (err error, entity Entity) {
	err = Database.Open(Entity{}).Where("id", "=", id).First().AsEntity(&entity)
	if err != nil {
		panic(err)
	}
	return err, entity
}

/* Testcode:
e := NewEntity()
e.Name = "Role"
e.addField(Field{Name: "Id", Type: inttype, Object: ""})
e.addField(Field{Name: "Name", Type: stringtype, Object: ""})
e.addField(Field{Name: "Accounts", Type: slicetype, Object: "Account"})

err := Database.Insert(e)
if err != nil {
	panic(err)
}

_, es := getAllEntities()
*/
