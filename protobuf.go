package spbc

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/volatiletech/null/v8"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PBO represents a protobuf message of a specific type, controlled by the generic type parameter.
//
// Note that this handles NULL in postgres separately from an empty protobuf object. A NULL column value will be
// represented with a false value for the Valid field.
type PBO[T proto.Message] struct {
	Object T
	Valid  bool
}

// Assert that expected interfaces are implemented
var _ (json.Marshaler) = (*PBO[*emptypb.Empty])(nil)
var _ (json.Unmarshaler) = (*PBO[*emptypb.Empty])(nil)
var _ (encoding.TextMarshaler) = (*PBO[*emptypb.Empty])(nil)
var _ (encoding.TextUnmarshaler) = (*PBO[*emptypb.Empty])(nil)
var _ (driver.Valuer) = (*PBO[*emptypb.Empty])(nil)
var _ (sql.Scanner) = (*PBO[*emptypb.Empty])(nil)

// NewPBO creates a new PBO object.
func NewPBO[T proto.Message](o T, valid bool) PBO[T] {
	return PBO[T]{
		Object: o,
		Valid:  valid,
	}
}

// PBOFrom creates a new PBO from a valid protobuf object.
//
// Note that proto.Message is inherently a pointer to a protobuf object and thus can be nil, but golang's type system
// isn't smart enough to allow nil checks against the generic type since it can't guarantee it's a pointer. As such,
// we have to resort to a lousy reflect check.
func PBOFrom[T proto.Message](o T) PBO[T] {
	return NewPBO(o, !reflect.ValueOf(o).IsNil())
}

// IsValid returns true if this carries an explicit value.
func (o PBO[T]) IsValid() bool {
	return o.Valid
}

// UnmarshalJSON implements json.Unmarshaler.
func (o *PBO[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, null.NullBytes) {
		o.Valid = false
		return nil
	}

	// If the object is nil, generate a new empty object to receive the unmarshalled data.
	if reflect.ValueOf(o.Object).IsNil() {
		o.Object = o.Object.ProtoReflect().New().Interface().(T)
	}

	if err := protojson.Unmarshal(data, o.Object); err != nil {
		return err
	}
	o.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (o *PBO[T]) UnmarshalText(text []byte) error {
	if text == nil {
		o.Valid = false
		return nil
	}

	// If the object is nil, generate a new empty object to receive the unmarshalled data.
	if reflect.ValueOf(o.Object).IsNil() {
		o.Object = o.Object.ProtoReflect().New().Interface().(T)
	}

	if err := proto.Unmarshal(text, o.Object); err != nil {
		return err
	}
	o.Valid = false
	return nil
}

// MarshalJSON implements json.Marshaler.
func (o PBO[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return null.NullBytes, nil
	}

	return protojson.Marshal(o.Object)
}

// MarshalText implements encoding.TextMarshaler.
func (o PBO[T]) MarshalText() ([]byte, error) {
	if !o.Valid {
		return []byte{}, nil
	}
	return proto.Marshal(o.Object)
}

// Scan implements the sql.Scanner interface.
func (o *PBO[T]) Scan(value interface{}) error {
	if value == nil {
		o.Valid = false
		return nil
	}

	bytes, isBytes := value.([]byte)
	if !isBytes {
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type PBO", value)
	}

	// If the object is nil, generate a new empty object to receive the unmarshalled data.
	if reflect.ValueOf(o.Object).IsNil() {
		o.Object = o.Object.ProtoReflect().New().Interface().(T)
	}

	o.Valid = true
	return proto.Unmarshal(bytes, o.Object)
}

// Value implements the driver.Valuer interface.
func (o PBO[T]) Value() (driver.Value, error) {
	if !o.Valid {
		return nil, nil
	}

	return proto.Marshal(o.Object)
}
