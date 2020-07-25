package schema

import (
	"testing"
	"what-unexpected-summer/gbey/gbey/orm/dialect"
)

// schema_test.go
type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)//找到表对应的结构体，然后用结构体做表名建新表
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}