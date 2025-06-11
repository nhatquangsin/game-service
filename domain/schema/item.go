package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("name"),
		field.String("category"),
		field.String("description"),
		field.Int64("created_at").DefaultFunc(func() int64 {
			return time.Now().Unix()
		}).Immutable(),
		field.Int64("updated_at").
			DefaultFunc(func() int64 {
				return time.Now().Unix()
			}).
			UpdateDefault(func() int64 {
				return time.Now().Unix()
			}),
	}
}
