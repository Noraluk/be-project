package base

import (
	"be-project/pkg/database"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository[T any] interface {
	Table(name string, args ...interface{}) BaseRepository[T]
	Take(dest interface{}, conds ...interface{}) BaseRepository[T]
	First(dest interface{}, conds ...interface{}) BaseRepository[T]
	Last(dest interface{}, conds ...interface{}) BaseRepository[T]
	Find(dest interface{}, conds ...interface{}) BaseRepository[T]
	Create(t interface{}) BaseRepository[T]
	FirstOrCreate(dest interface{}, conds ...interface{}) BaseRepository[T]
	Select(query interface{}, args ...interface{}) BaseRepository[T]
	Save(t interface{}) BaseRepository[T]
	Update(column string, value interface{}) BaseRepository[T]
	Updates(values interface{}) BaseRepository[T]
	Delete(value interface{}, conds ...interface{}) BaseRepository[T]
	Where(query interface{}, args ...interface{}) BaseRepository[T]
	Joins(query string, args ...interface{}) BaseRepository[T]
	Group(name string) BaseRepository[T]
	Order(value interface{}) BaseRepository[T]
	Count(count *int64) BaseRepository[T]
	Scan(dest interface{}) BaseRepository[T]

	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error

	Omit(column ...string) BaseRepository[T]
	Model(value interface{}) BaseRepository[T]
	Preload(query string, args ...interface{}) BaseRepository[T]

	Session(config *gorm.Session) BaseRepository[T]

	Clauses(conds ...clause.Expression) BaseRepository[T]
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) BaseRepository[T]

	Error() error
	RowsAffected() int64
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any]() BaseRepository[T] {
	return &baseRepository[T]{
		db: database.GetDatabase(),
	}
}

func Wrap[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{
		db: db,
	}
}

func (b baseRepository[T]) Table(name string, args ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Table(name, args...))
}

func (b baseRepository[T]) Take(dest interface{}, conds ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Take(dest, conds...))
}

func (b baseRepository[T]) First(dest interface{}, conds ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.First(dest, conds...))
}

func (b baseRepository[T]) Last(dest interface{}, conds ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Last(dest, conds...))
}

func (b baseRepository[T]) Find(dest interface{}, conds ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Find(dest, conds...))
}

func (b baseRepository[T]) Create(t interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Create(t))
}

func (b baseRepository[T]) FirstOrCreate(dest interface{}, conds ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.FirstOrCreate(dest, conds...))
}

func (b baseRepository[T]) Select(query interface{}, args ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Select(query, args...))
}

func (b baseRepository[T]) Save(t interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Save(t))
}

func (b baseRepository[T]) Update(column string, value interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Update(column, value))
}

func (b baseRepository[T]) Updates(values interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Updates(values))
}

func (b baseRepository[T]) Delete(value interface{}, conds ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Delete(value, conds...))
}

func (b baseRepository[T]) Where(query interface{}, args ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Where(query, args...))
}

func (b baseRepository[T]) Joins(query string, args ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Joins(query, args...))
}

func (b baseRepository[T]) Group(name string) BaseRepository[T] {
	return Wrap[T](b.db.Group(name))
}

func (b baseRepository[T]) Order(value interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Order(value))
}

func (b baseRepository[T]) Count(count *int64) BaseRepository[T] {
	return Wrap[T](b.db.Count(count))
}

func (b baseRepository[T]) Scan(dest interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Scan(dest))
}

func (b baseRepository[T]) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return b.db.Transaction(fc, opts...)
}

func (b baseRepository[T]) Omit(column ...string) BaseRepository[T] {
	return Wrap[T](b.db.Omit(column...))
}

func (b baseRepository[T]) Model(value interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Model(value))
}

func (b baseRepository[T]) Preload(query string, args ...interface{}) BaseRepository[T] {
	return Wrap[T](b.db.Preload(query, args...))
}

func (b baseRepository[T]) Session(config *gorm.Session) BaseRepository[T] {
	return Wrap[T](b.db.Session(config))
}

func (b baseRepository[T]) Clauses(conds ...clause.Expression) BaseRepository[T] {
	return Wrap[T](b.db.Clauses(conds...))
}

func (b baseRepository[T]) Scopes(funcs ...func(*gorm.DB) *gorm.DB) BaseRepository[T] {
	return Wrap[T](b.db.Scopes(funcs...))
}

func (b baseRepository[T]) Error() error {
	return b.db.Error
}

func (b baseRepository[T]) RowsAffected() int64 {
	return b.db.RowsAffected
}
