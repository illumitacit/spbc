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

// PBJSON represents a protojson message of a specific type, controlled by the generic type parameter.
//
// Note that this handles NULL in postgres separately from an empty protobuf object. A NULL column value will be
// represented with a false value for the Valid field.
type PBJSON[T proto.Message] struct {
	Object T
	Valid  bool
}

// Assert that expected interfaces are implemented
var _ (json.Marshaler) = (*PBJSON[*emptypb.Empty])(nil)
var _ (json.Unmarshaler) = (*PBJSON[*emptypb.Empty])(nil)
var _ (encoding.TextMarshaler) = (*PBJSON[*emptypb.Empty])(nil)
var _ (encoding.TextUnmarshaler) = (*PBJSON[*emptypb.Empty])(nil)
var _ (driver.Valuer) = (*PBJSON[*emptypb.Empty])(nil)
var _ (sql.Scanner) = (*PBJSON[*emptypb.Empty])(nil)

// NewPBJSON creates a new PBJSON object.
func NewPBJSON[T proto.Message](o T, valid bool) PBJSON[T] {
	return PBJSON[T]{
		Object: o,
		Valid:  valid,
	}
}

// PBJSONFrom creates a new PBJSON from a valid protobuf object.
//
// Note that proto.Message is inherently a pointer to a protobuf object and thus can be nil, but golang's type system
// isn't smart enough to allow nil checks against the generic type since it can't guarantee it's a pointer. As such,
// we have to resort to a lousy reflect check.
func PBJSONFrom[T proto.Message](o T) PBJSON[T] {
	return NewPBJSON(o, !reflect.ValueOf(o).IsNil())
}

// IsValid returns true if this carries an explicit value.
func (o PBJSON[T]) IsValid() bool {
	return o.Valid
}

// UnmarshalJSON implements json.Unmarshaler.
func (o *PBJSON[T]) UnmarshalJSON(data []byte) error {
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
func (o *PBJSON[T]) UnmarshalText(text []byte) error {
	return o.UnmarshalJSON(text)
}

// MarshalJSON implements json.Marshaler.
func (o PBJSON[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return null.NullBytes, nil
	}

	return protojson.Marshal(o.Object)
}

// MarshalText implements encoding.TextMarshaler.
func (o PBJSON[T]) MarshalText() ([]byte, error) {
	return o.MarshalJSON()
}

// Scan implements the sql.Scanner interface.
func (o *PBJSON[T]) Scan(value interface{}) error {
	if value == nil {
		o.Valid = false
		return nil
	}

	bytes, isBytes := value.([]byte)
	if !isBytes {
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type PBJSON", value)
	}

	return o.UnmarshalJSON(bytes)
}

// Value implements the driver.Valuer interface.
func (o PBJSON[T]) Value() (driver.Value, error) {
	return o.MarshalJSON()
}
