// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"journeyhub/ent/file"
	"journeyhub/ent/message"
	"journeyhub/ent/messagevoice"
	"journeyhub/ent/room"
	"journeyhub/ent/schema/pulid"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// MessageVoice is the model entity for the MessageVoice schema.
type MessageVoice struct {
	config `json:"-"`
	// ID of the ent.
	ID pulid.ID `json:"id,omitempty"`
	// Length holds the value of the "length" field.
	Length uint64 `json:"length,omitempty"`
	// AttachedAt holds the value of the "attached_at" field.
	AttachedAt time.Time `json:"attached_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MessageVoiceQuery when eager-loading is set.
	Edges               MessageVoiceEdges `json:"edges"`
	file_message_voice  *pulid.ID
	message_voice       *pulid.ID
	room_message_voices *pulid.ID
	selectValues        sql.SelectValues
}

// MessageVoiceEdges holds the relations/edges for other nodes in the graph.
type MessageVoiceEdges struct {
	// Room holds the value of the room edge.
	Room *Room `json:"room,omitempty"`
	// Message holds the value of the message edge.
	Message *Message `json:"message,omitempty"`
	// File holds the value of the file edge.
	File *File `json:"file,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
	// totalCount holds the count of the edges above.
	totalCount [3]map[string]int
}

// RoomOrErr returns the Room value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MessageVoiceEdges) RoomOrErr() (*Room, error) {
	if e.Room != nil {
		return e.Room, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: room.Label}
	}
	return nil, &NotLoadedError{edge: "room"}
}

// MessageOrErr returns the Message value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MessageVoiceEdges) MessageOrErr() (*Message, error) {
	if e.Message != nil {
		return e.Message, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: message.Label}
	}
	return nil, &NotLoadedError{edge: "message"}
}

// FileOrErr returns the File value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MessageVoiceEdges) FileOrErr() (*File, error) {
	if e.File != nil {
		return e.File, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: file.Label}
	}
	return nil, &NotLoadedError{edge: "file"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*MessageVoice) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case messagevoice.FieldID:
			values[i] = new(pulid.ID)
		case messagevoice.FieldLength:
			values[i] = new(sql.NullInt64)
		case messagevoice.FieldAttachedAt:
			values[i] = new(sql.NullTime)
		case messagevoice.ForeignKeys[0]: // file_message_voice
			values[i] = &sql.NullScanner{S: new(pulid.ID)}
		case messagevoice.ForeignKeys[1]: // message_voice
			values[i] = &sql.NullScanner{S: new(pulid.ID)}
		case messagevoice.ForeignKeys[2]: // room_message_voices
			values[i] = &sql.NullScanner{S: new(pulid.ID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the MessageVoice fields.
func (mv *MessageVoice) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case messagevoice.FieldID:
			if value, ok := values[i].(*pulid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				mv.ID = *value
			}
		case messagevoice.FieldLength:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field length", values[i])
			} else if value.Valid {
				mv.Length = uint64(value.Int64)
			}
		case messagevoice.FieldAttachedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field attached_at", values[i])
			} else if value.Valid {
				mv.AttachedAt = value.Time
			}
		case messagevoice.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field file_message_voice", values[i])
			} else if value.Valid {
				mv.file_message_voice = new(pulid.ID)
				*mv.file_message_voice = *value.S.(*pulid.ID)
			}
		case messagevoice.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field message_voice", values[i])
			} else if value.Valid {
				mv.message_voice = new(pulid.ID)
				*mv.message_voice = *value.S.(*pulid.ID)
			}
		case messagevoice.ForeignKeys[2]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field room_message_voices", values[i])
			} else if value.Valid {
				mv.room_message_voices = new(pulid.ID)
				*mv.room_message_voices = *value.S.(*pulid.ID)
			}
		default:
			mv.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the MessageVoice.
// This includes values selected through modifiers, order, etc.
func (mv *MessageVoice) Value(name string) (ent.Value, error) {
	return mv.selectValues.Get(name)
}

// QueryRoom queries the "room" edge of the MessageVoice entity.
func (mv *MessageVoice) QueryRoom() *RoomQuery {
	return NewMessageVoiceClient(mv.config).QueryRoom(mv)
}

// QueryMessage queries the "message" edge of the MessageVoice entity.
func (mv *MessageVoice) QueryMessage() *MessageQuery {
	return NewMessageVoiceClient(mv.config).QueryMessage(mv)
}

// QueryFile queries the "file" edge of the MessageVoice entity.
func (mv *MessageVoice) QueryFile() *FileQuery {
	return NewMessageVoiceClient(mv.config).QueryFile(mv)
}

// Update returns a builder for updating this MessageVoice.
// Note that you need to call MessageVoice.Unwrap() before calling this method if this MessageVoice
// was returned from a transaction, and the transaction was committed or rolled back.
func (mv *MessageVoice) Update() *MessageVoiceUpdateOne {
	return NewMessageVoiceClient(mv.config).UpdateOne(mv)
}

// Unwrap unwraps the MessageVoice entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (mv *MessageVoice) Unwrap() *MessageVoice {
	_tx, ok := mv.config.driver.(*txDriver)
	if !ok {
		panic("ent: MessageVoice is not a transactional entity")
	}
	mv.config.driver = _tx.drv
	return mv
}

// String implements the fmt.Stringer.
func (mv *MessageVoice) String() string {
	var builder strings.Builder
	builder.WriteString("MessageVoice(")
	builder.WriteString(fmt.Sprintf("id=%v, ", mv.ID))
	builder.WriteString("length=")
	builder.WriteString(fmt.Sprintf("%v", mv.Length))
	builder.WriteString(", ")
	builder.WriteString("attached_at=")
	builder.WriteString(mv.AttachedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// MessageVoices is a parsable slice of MessageVoice.
type MessageVoices []*MessageVoice
