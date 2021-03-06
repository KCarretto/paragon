// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/kcarretto/paragon/ent/task"
)

// Task is the model entity for the Task schema.
type Task struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// QueueTime holds the value of the "QueueTime" field.
	QueueTime time.Time `json:"QueueTime,omitempty"`
	// LastChangedTime holds the value of the "LastChangedTime" field.
	LastChangedTime time.Time `json:"LastChangedTime,omitempty"`
	// ClaimTime holds the value of the "ClaimTime" field.
	ClaimTime time.Time `json:"ClaimTime,omitempty"`
	// ExecStartTime holds the value of the "ExecStartTime" field.
	ExecStartTime time.Time `json:"ExecStartTime,omitempty"`
	// ExecStopTime holds the value of the "ExecStopTime" field.
	ExecStopTime time.Time `json:"ExecStopTime,omitempty"`
	// Content holds the value of the "Content" field.
	Content string `json:"Content,omitempty"`
	// Output holds the value of the "Output" field.
	Output string `json:"Output,omitempty"`
	// Error holds the value of the "Error" field.
	Error string `json:"Error,omitempty"`
	// SessionID holds the value of the "SessionID" field.
	SessionID string `json:"SessionID,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TaskQuery when eager-loading is set.
	Edges struct {
		// Tags holds the value of the tags edge.
		Tags []*Tag
		// Job holds the value of the job edge.
		Job *Job
		// Target holds the value of the target edge.
		Target *Target
	} `json:"edges"`
	job_id    *int
	target_id *int
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Task) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // QueueTime
		&sql.NullTime{},   // LastChangedTime
		&sql.NullTime{},   // ClaimTime
		&sql.NullTime{},   // ExecStartTime
		&sql.NullTime{},   // ExecStopTime
		&sql.NullString{}, // Content
		&sql.NullString{}, // Output
		&sql.NullString{}, // Error
		&sql.NullString{}, // SessionID
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Task) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // job_id
		&sql.NullInt64{}, // target_id
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Task fields.
func (t *Task) assignValues(values ...interface{}) error {
	if m, n := len(values), len(task.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	t.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field QueueTime", values[0])
	} else if value.Valid {
		t.QueueTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field LastChangedTime", values[1])
	} else if value.Valid {
		t.LastChangedTime = value.Time
	}
	if value, ok := values[2].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field ClaimTime", values[2])
	} else if value.Valid {
		t.ClaimTime = value.Time
	}
	if value, ok := values[3].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field ExecStartTime", values[3])
	} else if value.Valid {
		t.ExecStartTime = value.Time
	}
	if value, ok := values[4].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field ExecStopTime", values[4])
	} else if value.Valid {
		t.ExecStopTime = value.Time
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field Content", values[5])
	} else if value.Valid {
		t.Content = value.String
	}
	if value, ok := values[6].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field Output", values[6])
	} else if value.Valid {
		t.Output = value.String
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field Error", values[7])
	} else if value.Valid {
		t.Error = value.String
	}
	if value, ok := values[8].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field SessionID", values[8])
	} else if value.Valid {
		t.SessionID = value.String
	}
	values = values[9:]
	if len(values) == len(task.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field job_id", value)
		} else if value.Valid {
			t.job_id = new(int)
			*t.job_id = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field target_id", value)
		} else if value.Valid {
			t.target_id = new(int)
			*t.target_id = int(value.Int64)
		}
	}
	return nil
}

// QueryTags queries the tags edge of the Task.
func (t *Task) QueryTags() *TagQuery {
	return (&TaskClient{t.config}).QueryTags(t)
}

// QueryJob queries the job edge of the Task.
func (t *Task) QueryJob() *JobQuery {
	return (&TaskClient{t.config}).QueryJob(t)
}

// QueryTarget queries the target edge of the Task.
func (t *Task) QueryTarget() *TargetQuery {
	return (&TaskClient{t.config}).QueryTarget(t)
}

// Update returns a builder for updating this Task.
// Note that, you need to call Task.Unwrap() before calling this method, if this Task
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Task) Update() *TaskUpdateOne {
	return (&TaskClient{t.config}).UpdateOne(t)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (t *Task) Unwrap() *Task {
	tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Task is not a transactional entity")
	}
	t.config.driver = tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Task) String() string {
	var builder strings.Builder
	builder.WriteString("Task(")
	builder.WriteString(fmt.Sprintf("id=%v", t.ID))
	builder.WriteString(", QueueTime=")
	builder.WriteString(t.QueueTime.Format(time.ANSIC))
	builder.WriteString(", LastChangedTime=")
	builder.WriteString(t.LastChangedTime.Format(time.ANSIC))
	builder.WriteString(", ClaimTime=")
	builder.WriteString(t.ClaimTime.Format(time.ANSIC))
	builder.WriteString(", ExecStartTime=")
	builder.WriteString(t.ExecStartTime.Format(time.ANSIC))
	builder.WriteString(", ExecStopTime=")
	builder.WriteString(t.ExecStopTime.Format(time.ANSIC))
	builder.WriteString(", Content=")
	builder.WriteString(t.Content)
	builder.WriteString(", Output=")
	builder.WriteString(t.Output)
	builder.WriteString(", Error=")
	builder.WriteString(t.Error)
	builder.WriteString(", SessionID=")
	builder.WriteString(t.SessionID)
	builder.WriteByte(')')
	return builder.String()
}

// Tasks is a parsable slice of Task.
type Tasks []*Task

func (t Tasks) config(cfg config) {
	for _i := range t {
		t[_i].config = cfg
	}
}
