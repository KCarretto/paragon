package schema

import (
	"math"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// File holds file content and metadata.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	content := field.Bytes("Content")
	content.Descriptor().Size = math.MaxUint32

	return []ent.Field{
		field.String("Name").
			NotEmpty().
			MinLen(1).
			Unique().
			Comment("The name of the file, used to reference it for downloads"),
		field.Time("CreationTime").
			Default(func() time.Time {
				return time.Now()
			}).
			Comment("The timestamp for when the File was created"),
		field.Time("LastModifiedTime").
			Comment("The timestamp for when the File was last modified"),
		field.Int("Size").
			Default(0).
			Min(0).
			Comment("The size of the file in bytes"),
		content.
			Comment("The content of the file"),
		field.String("Hash").
			MaxLen(100).
			Comment("A SHA3 digest of the content field"),
		field.String("ContentType").
			Comment("The content type of content"),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("links", Link.Type).
			Comment("A File can have many links assigned to it"),
	}
}
