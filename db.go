package fixtures

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewWithGormDB(gormDB *gorm.DB) *DB {
	return &DB{DB: gormDB}
}

func NewDB(dialector gorm.Dialector, opts ...gorm.Option) func(t *testing.T) *DB {
	return func(t *testing.T) *DB {
		gormDB, err := gorm.Open(dialector, opts...)
		if err != nil {
			t.Fatalf("failed to gorm.Open: %v", err)
		}
		return NewWithGormDB(gormDB)
	}
}

func (f *DB) Create(t *testing.T, values ...interface{}) {
	for _, i := range values {
		if r := f.DB.Create(i); r.Error != nil {
			t.Fatalf("failed to create: %#v because of %+v", i, r.Error)
		}
	}
}

func (f *DB) Update(t *testing.T, values ...interface{}) {
	for _, i := range values {
		if r := f.DB.Updates(i); r.Error != nil {
			t.Fatalf("failed to update: %#v because of %+v", i, r.Error)
		}
	}
}

func (f *DB) Delete(t *testing.T, values ...interface{}) {
	for _, i := range values {
		if r := f.DB.Delete(i); r.Error != nil {
			t.Fatalf("failed to delete: %#v because of %+v", i, r.Error)
		}
	}
}

func (f *DB) DeleteFromTable(t *testing.T, models ...interface{}) {
	for _, m := range models {
		stmt := &gorm.Statement{DB: f.DB}
		if err := stmt.Parse(m); err != nil {
			t.Fatalf("failed to parse model: %v", err)
		}
		if r := f.DB.Exec(fmt.Sprintf("DELETE FROM %s", stmt.Schema.Table)); r.Error != nil {
			t.Fatalf("failed to delete table: %v", r.Error)
		}
	}
}
