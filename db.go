package fixtures

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

// DB wraps a gorm.DB instance.
type DB struct {
	*gorm.DB
}

// NewWithGormDB creates a new DB instance with an existing gorm.DB.
func NewWithGormDB(gormDB *gorm.DB) *DB {
	return &DB{DB: gormDB}
}

// NewDB initializes a new DB instance with the given gorm.Dialector and options.
func NewDB(dialector gorm.Dialector, opts ...gorm.Option) func(t *testing.T) *DB {
	return func(t *testing.T) *DB {
		gormDB, err := gorm.Open(dialector, opts...)
		if err != nil {
			t.Fatalf("failed to gorm.Open: %v", err)
		}
		return NewWithGormDB(gormDB)
	}
}

// Create inserts the given values into the database.
func (f *DB) Create(t *testing.T, values ...interface{}) {
	for _, i := range values {
		if r := f.DB.Create(i); r.Error != nil {
			t.Fatalf("failed to create: %#v because of %+v", i, r.Error)
		}
	}
}

// CreateAndReload creates the given values in the database and then reloads them.
func (f *DB) CreateAndReload(t *testing.T, values ...interface{}) {
	f.Create(t, values...)
	f.Reload(t, values...)
}

// Update modifies the given values in the database.
func (f *DB) Update(t *testing.T, values ...interface{}) {
	for _, i := range values {
		if r := f.DB.Updates(i); r.Error != nil {
			t.Fatalf("failed to update: %#v because of %+v", i, r.Error)
		}
	}
}

func (f *DB) UpdateAndReload(t *testing.T, values ...interface{}) {
	f.Update(t, values...)
	f.Reload(t, values...)
}

// Delete removes the given values from the database.
func (f *DB) Delete(t *testing.T, values ...interface{}) {
	for _, i := range values {
		if r := f.DB.Delete(i); r.Error != nil {
			t.Fatalf("failed to delete: %#v because of %+v", i, r.Error)
		}
	}
}

// DeleteFromTable deletes all records from the tables of the given models.
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

// Reload reloads the given models from the database.
func (f *DB) Reload(t *testing.T, models ...interface{}) {
	for _, i := range models {
		if r := f.DB.First(i); r.Error != nil {
			t.Fatalf("failed to reload: %#v because of %+v", i, r.Error)
		}
	}
}
