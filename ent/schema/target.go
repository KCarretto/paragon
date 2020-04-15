package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Target holds the schema definition for the Target entity.
type Target struct {
	ent.Schema
}

type OSType string

var osTypes []string

func newOSType(value string) OSType {
	osTypes = append(osTypes, value)
	return OSType(value)
}

// IF you update this you MUST regen the ent stuff
var (
	LinuxOS   OSType = newOSType("LINUX")
	WindowsOS OSType = newOSType("WINDOWS")
	BSDOS     OSType = newOSType("BSD")
)

// Fields of the Target.
func (Target) Fields() []ent.Field {
	return []ent.Field{
		field.String("Name").
			Unique().
			Comment("The name of the Target"),
		field.Enum("OS").
			Values(osTypes...).
			Comment("The OS of the target"),
		field.String("PrimaryIP").
			Comment("The IP Address for the primary interface of the Target"),
		field.String("MachineUUID").
			MaxLen(250).
			Unique().
			Optional().
			Comment("The machine UUID of the Target"),
		field.String("PublicIP").
			Optional().
			Comment("The Public IP Address for the Target"),
		field.String("PrimaryMAC").
			Optional().
			Comment("The MAC Address for the primary interface of the Target"),
		field.String("Hostname").
			Optional().
			Comment("The hostname for the Target"),
		field.Time("LastSeen").
			Optional().
			Comment("The time the Target was last seen"),
	}
}

// Edges of the Target.
func (Target) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type).
			Comment("A Target can have many tasks connected to it"),
		edge.To("tags", Tag.Type).
			Comment("A Target can have many Tags"),
		edge.To("credentials", Credential.Type).
			Comment("A Target can have many credentials connected to it"),
	}
}
