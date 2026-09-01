package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func ProfileImageURL(key string) string {
	if key == "" {
		return ""
	}
	return fmt.Sprintf("http://s3-api.garage.local/%s", key)
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("profile_link").Default(""), // will put garage heree!!1 
		field.String("username").NotEmpty().Unique(),
		field.String("email").Unique(),
		field.String("password_hash").NotEmpty().Sensitive(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
