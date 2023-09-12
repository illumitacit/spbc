package spbc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestPBJSONFromValidObj(t *testing.T) {
	o := PBJSONFrom(&_emptyObj)
	assert.True(t, proto.Equal(o.Object, &_emptyObj))
}

func TestPBJSONFromNilObj(t *testing.T) {
	nilObj := PBJSONFrom[*emptypb.Empty](nil)
	assert.False(t, nilObj.Valid, "expected nil PBJSON to be invalid")
}

func TestUnmarshalPBJSONValid(t *testing.T) {
	var o PBJSON[*emptypb.Empty]
	err := json.Unmarshal([]byte("{}"), &o)
	require.NoError(t, err)
	assert.True(t, proto.Equal(o.Object, &_emptyObj))
}

func TestUnmarshalPBJSONInvalid(t *testing.T) {
	var o PBJSON[*emptypb.Empty]
	err := json.Unmarshal([]byte("56"), &o)
	require.Error(t, err)
}

func TestTextUnmarshalPBJSONValid(t *testing.T) {
	var o PBJSON[*emptypb.Empty]
	err := o.UnmarshalText([]byte("{}"))
	require.NoError(t, err)
	assert.True(t, proto.Equal(o.Object, &_emptyObj))
}

func TestTextUnmarshalPBJSONInvalid(t *testing.T) {
	var o PBJSON[*emptypb.Empty]
	err := o.UnmarshalText([]byte("42"))
	require.Error(t, err)
}

func TestMarshalPBJSONValid(t *testing.T) {
	o := PBJSONFrom(&_emptyObj)
	data, err := json.Marshal(o)
	require.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestMarshalPBJSONInvalid(t *testing.T) {
	// invalid values should be encoded as null
	null := NewPBJSON[*emptypb.Empty](nil, false)
	data, err := json.Marshal(null)
	require.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

func TestMarshalPBJSONText(t *testing.T) {
	o := PBJSONFrom(&_emptyObj)
	data, err := o.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestPBJSONScanValid(t *testing.T) {
	var o PBJSON[*emptypb.Empty]
	err := o.Scan([]byte("{}"))
	require.NoError(t, err)
	assert.True(t, proto.Equal(o.Object, &_emptyObj))
}

func TestPBJSONScanNull(t *testing.T) {
	var null PBJSON[*emptypb.Empty]
	err := null.Scan(nil)
	require.NoError(t, err)
	assert.False(t, null.Valid)
}
