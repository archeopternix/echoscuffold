// AccountManager project main.go
package main

// Entity relates to an 'Object' or struct
type Entity struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

// Field is each and every single attribute.
// Object is empty except in case type=slicetype keeps the name of the Object
type Field struct {
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Object string `json:"object"`
}

const (
	stringtype int = iota + 1
	inttype
	booltype
	slicetype
)

func NewEntity() (e *Entity) {
	_, id := NextId("Entity")
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
