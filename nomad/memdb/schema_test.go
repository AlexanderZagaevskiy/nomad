package memdb

import "testing"

func testValidSchema() *DBSchema {
	return &DBSchema{
		Tables: map[string]*TableSchema{
			"main": &TableSchema{
				Name: "main",
				Indexes: map[string]*IndexSchema{
					"id": &IndexSchema{
						Name:    "id",
						Indexer: StringFieldIndex("ID", false),
					},
				},
			},
		},
	}
}

func TestDBSchema_Validate(t *testing.T) {
	s := &DBSchema{}
	err := s.Validate()
	if err == nil {
		t.Fatalf("should not validate, empty")
	}

	s.Tables = map[string]*TableSchema{
		"foo": &TableSchema{Name: "foo"},
	}
	err = s.Validate()
	if err == nil {
		t.Fatalf("should not validate, no indexes")
	}

	valid := testValidSchema()
	err = valid.Validate()
	if err != nil {
		t.Fatalf("should validate: %v", err)
	}
}

func TestTableSchema_Validate(t *testing.T) {
	s := &TableSchema{}
	err := s.Validate()
	if err == nil {
		t.Fatalf("should not validate, empty")
	}

	s.Indexes = map[string]*IndexSchema{
		"foo": &IndexSchema{Name: "foo"},
	}
	err = s.Validate()
	if err == nil {
		t.Fatalf("should not validate, no indexes")
	}

	valid := &TableSchema{
		Name: "main",
		Indexes: map[string]*IndexSchema{
			"id": &IndexSchema{
				Name:    "id",
				Indexer: StringFieldIndex("ID", true),
			},
		},
	}
	err = valid.Validate()
	if err != nil {
		t.Fatalf("should validate: %v", err)
	}
}

func TestIndexSchema_Validate(t *testing.T) {
	s := &IndexSchema{}
	err := s.Validate()
	if err == nil {
		t.Fatalf("should not validate, empty")
	}

	s.Name = "foo"
	err = s.Validate()
	if err == nil {
		t.Fatalf("should not validate, no indexer")
	}

	s.Indexer = StringFieldIndex("Foo", false)
	err = s.Validate()
	if err != nil {
		t.Fatalf("should validate: %v", err)
	}
}
