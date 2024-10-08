// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"journeyhub/ent/device"
	"journeyhub/ent/schema/pulid"
	"journeyhub/ent/user"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID pulid.ID `json:"id,omitempty"`
	// FirstName holds the value of the "first_name" field.
	FirstName string `json:"first_name,omitempty"`
	// LastName holds the value of the "last_name" field.
	LastName string `json:"last_name,omitempty"`
	// Nickname holds the value of the "nickname" field.
	Nickname string `json:"nickname,omitempty"`
	// Email holds the value of the "email" field.
	Email string `json:"email,omitempty"`
	// ContactPin holds the value of the "contact_pin" field.
	ContactPin string `json:"contact_pin,omitempty"`
	// Password holds the value of the "password" field.
	Password string `json:"-"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserQuery when eager-loading is set.
	Edges        UserEdges `json:"edges"`
	selectValues sql.SelectValues
}

// UserEdges holds the relations/edges for other nodes in the graph.
type UserEdges struct {
	// Device holds the value of the device edge.
	Device *Device `json:"device,omitempty"`
	// Notifications holds the value of the notifications edge.
	Notifications []*Notification `json:"notifications,omitempty"`
	// Contacts holds the value of the contacts edge.
	Contacts []*User `json:"contacts,omitempty"`
	// Rooms holds the value of the rooms edge.
	Rooms []*Room `json:"rooms,omitempty"`
	// Messages holds the value of the messages edge.
	Messages []*Message `json:"messages,omitempty"`
	// UserContacts holds the value of the user_contacts edge.
	UserContacts []*UserContact `json:"user_contacts,omitempty"`
	// Memberships holds the value of the memberships edge.
	Memberships []*RoomMember `json:"memberships,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [7]bool
	// totalCount holds the count of the edges above.
	totalCount [7]map[string]int

	namedNotifications map[string][]*Notification
	namedContacts      map[string][]*User
	namedRooms         map[string][]*Room
	namedMessages      map[string][]*Message
	namedUserContacts  map[string][]*UserContact
	namedMemberships   map[string][]*RoomMember
}

// DeviceOrErr returns the Device value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserEdges) DeviceOrErr() (*Device, error) {
	if e.Device != nil {
		return e.Device, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: device.Label}
	}
	return nil, &NotLoadedError{edge: "device"}
}

// NotificationsOrErr returns the Notifications value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) NotificationsOrErr() ([]*Notification, error) {
	if e.loadedTypes[1] {
		return e.Notifications, nil
	}
	return nil, &NotLoadedError{edge: "notifications"}
}

// ContactsOrErr returns the Contacts value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) ContactsOrErr() ([]*User, error) {
	if e.loadedTypes[2] {
		return e.Contacts, nil
	}
	return nil, &NotLoadedError{edge: "contacts"}
}

// RoomsOrErr returns the Rooms value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) RoomsOrErr() ([]*Room, error) {
	if e.loadedTypes[3] {
		return e.Rooms, nil
	}
	return nil, &NotLoadedError{edge: "rooms"}
}

// MessagesOrErr returns the Messages value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) MessagesOrErr() ([]*Message, error) {
	if e.loadedTypes[4] {
		return e.Messages, nil
	}
	return nil, &NotLoadedError{edge: "messages"}
}

// UserContactsOrErr returns the UserContacts value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) UserContactsOrErr() ([]*UserContact, error) {
	if e.loadedTypes[5] {
		return e.UserContacts, nil
	}
	return nil, &NotLoadedError{edge: "user_contacts"}
}

// MembershipsOrErr returns the Memberships value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) MembershipsOrErr() ([]*RoomMember, error) {
	if e.loadedTypes[6] {
		return e.Memberships, nil
	}
	return nil, &NotLoadedError{edge: "memberships"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case user.FieldID:
			values[i] = new(pulid.ID)
		case user.FieldFirstName, user.FieldLastName, user.FieldNickname, user.FieldEmail, user.FieldContactPin, user.FieldPassword:
			values[i] = new(sql.NullString)
		case user.FieldCreatedAt, user.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case user.FieldID:
			if value, ok := values[i].(*pulid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				u.ID = *value
			}
		case user.FieldFirstName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field first_name", values[i])
			} else if value.Valid {
				u.FirstName = value.String
			}
		case user.FieldLastName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field last_name", values[i])
			} else if value.Valid {
				u.LastName = value.String
			}
		case user.FieldNickname:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field nickname", values[i])
			} else if value.Valid {
				u.Nickname = value.String
			}
		case user.FieldEmail:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field email", values[i])
			} else if value.Valid {
				u.Email = value.String
			}
		case user.FieldContactPin:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field contact_pin", values[i])
			} else if value.Valid {
				u.ContactPin = value.String
			}
		case user.FieldPassword:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field password", values[i])
			} else if value.Valid {
				u.Password = value.String
			}
		case user.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				u.CreatedAt = value.Time
			}
		case user.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				u.UpdatedAt = value.Time
			}
		default:
			u.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the User.
// This includes values selected through modifiers, order, etc.
func (u *User) Value(name string) (ent.Value, error) {
	return u.selectValues.Get(name)
}

// QueryDevice queries the "device" edge of the User entity.
func (u *User) QueryDevice() *DeviceQuery {
	return NewUserClient(u.config).QueryDevice(u)
}

// QueryNotifications queries the "notifications" edge of the User entity.
func (u *User) QueryNotifications() *NotificationQuery {
	return NewUserClient(u.config).QueryNotifications(u)
}

// QueryContacts queries the "contacts" edge of the User entity.
func (u *User) QueryContacts() *UserQuery {
	return NewUserClient(u.config).QueryContacts(u)
}

// QueryRooms queries the "rooms" edge of the User entity.
func (u *User) QueryRooms() *RoomQuery {
	return NewUserClient(u.config).QueryRooms(u)
}

// QueryMessages queries the "messages" edge of the User entity.
func (u *User) QueryMessages() *MessageQuery {
	return NewUserClient(u.config).QueryMessages(u)
}

// QueryUserContacts queries the "user_contacts" edge of the User entity.
func (u *User) QueryUserContacts() *UserContactQuery {
	return NewUserClient(u.config).QueryUserContacts(u)
}

// QueryMemberships queries the "memberships" edge of the User entity.
func (u *User) QueryMemberships() *RoomMemberQuery {
	return NewUserClient(u.config).QueryMemberships(u)
}

// Update returns a builder for updating this User.
// Note that you need to call User.Unwrap() before calling this method if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return NewUserClient(u.config).UpdateOne(u)
}

// Unwrap unwraps the User entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	_tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = _tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v, ", u.ID))
	builder.WriteString("first_name=")
	builder.WriteString(u.FirstName)
	builder.WriteString(", ")
	builder.WriteString("last_name=")
	builder.WriteString(u.LastName)
	builder.WriteString(", ")
	builder.WriteString("nickname=")
	builder.WriteString(u.Nickname)
	builder.WriteString(", ")
	builder.WriteString("email=")
	builder.WriteString(u.Email)
	builder.WriteString(", ")
	builder.WriteString("contact_pin=")
	builder.WriteString(u.ContactPin)
	builder.WriteString(", ")
	builder.WriteString("password=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(u.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(u.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// NamedNotifications returns the Notifications named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedNotifications(name string) ([]*Notification, error) {
	if u.Edges.namedNotifications == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedNotifications[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedNotifications(name string, edges ...*Notification) {
	if u.Edges.namedNotifications == nil {
		u.Edges.namedNotifications = make(map[string][]*Notification)
	}
	if len(edges) == 0 {
		u.Edges.namedNotifications[name] = []*Notification{}
	} else {
		u.Edges.namedNotifications[name] = append(u.Edges.namedNotifications[name], edges...)
	}
}

// NamedContacts returns the Contacts named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedContacts(name string) ([]*User, error) {
	if u.Edges.namedContacts == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedContacts[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedContacts(name string, edges ...*User) {
	if u.Edges.namedContacts == nil {
		u.Edges.namedContacts = make(map[string][]*User)
	}
	if len(edges) == 0 {
		u.Edges.namedContacts[name] = []*User{}
	} else {
		u.Edges.namedContacts[name] = append(u.Edges.namedContacts[name], edges...)
	}
}

// NamedRooms returns the Rooms named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedRooms(name string) ([]*Room, error) {
	if u.Edges.namedRooms == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedRooms[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedRooms(name string, edges ...*Room) {
	if u.Edges.namedRooms == nil {
		u.Edges.namedRooms = make(map[string][]*Room)
	}
	if len(edges) == 0 {
		u.Edges.namedRooms[name] = []*Room{}
	} else {
		u.Edges.namedRooms[name] = append(u.Edges.namedRooms[name], edges...)
	}
}

// NamedMessages returns the Messages named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedMessages(name string) ([]*Message, error) {
	if u.Edges.namedMessages == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedMessages[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedMessages(name string, edges ...*Message) {
	if u.Edges.namedMessages == nil {
		u.Edges.namedMessages = make(map[string][]*Message)
	}
	if len(edges) == 0 {
		u.Edges.namedMessages[name] = []*Message{}
	} else {
		u.Edges.namedMessages[name] = append(u.Edges.namedMessages[name], edges...)
	}
}

// NamedUserContacts returns the UserContacts named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedUserContacts(name string) ([]*UserContact, error) {
	if u.Edges.namedUserContacts == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedUserContacts[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedUserContacts(name string, edges ...*UserContact) {
	if u.Edges.namedUserContacts == nil {
		u.Edges.namedUserContacts = make(map[string][]*UserContact)
	}
	if len(edges) == 0 {
		u.Edges.namedUserContacts[name] = []*UserContact{}
	} else {
		u.Edges.namedUserContacts[name] = append(u.Edges.namedUserContacts[name], edges...)
	}
}

// NamedMemberships returns the Memberships named value or an error if the edge was not
// loaded in eager-loading with this name.
func (u *User) NamedMemberships(name string) ([]*RoomMember, error) {
	if u.Edges.namedMemberships == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := u.Edges.namedMemberships[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (u *User) appendNamedMemberships(name string, edges ...*RoomMember) {
	if u.Edges.namedMemberships == nil {
		u.Edges.namedMemberships = make(map[string][]*RoomMember)
	}
	if len(edges) == 0 {
		u.Edges.namedMemberships[name] = []*RoomMember{}
	} else {
		u.Edges.namedMemberships[name] = append(u.Edges.namedMemberships[name], edges...)
	}
}

// Users is a parsable slice of User.
type Users []*User
