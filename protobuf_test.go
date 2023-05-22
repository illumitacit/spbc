package spbc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	_emptyObj = emptypb.Empty{}
)

func TestPBOFromValidObj(t *testing.T) {
	o := PBOFrom(&_emptyObj)
	assert.True(t, proto.Equal(o.Object, &_emptyObj))
}

func TestPBOFromNilObj(t *testing.T) {
	nilObj := PBOFrom[*emptypb.Empty](nil)
	assert.False(t, nilObj.Valid, "expected nil PBO to be invalid")
}

func TestUnmarshalPBOValid(t *testing.T) {
	var o PBO[*emptypb.Empty]
	err := json.Unmarshal([]byte("{}"), &o)
	require.NoError(t, err)
	assert.True(t, proto.Equal(o.Object, &_emptyObj))
}

func TestUnmarshalPBOInvalid(t *testing.T) {
	var o PBO[*emptypb.Empty]
	err := json.Unmarshal([]byte("56"), &o)
	require.Error(t, err)
}

func TestTextUnmarshalPBOValid(t *testing.T) {
	var o PBO[*emptypb.Empty]
	err := o.UnmarshalText([]byte(""))
	require.NoError(t, err)
	assert.True(t, proto.Equal(o.Object, &_emptyObj))
}

func TestTextUnmarshalPBOInvalid(t *testing.T) {
	var o PBO[*emptypb.Empty]
	err := o.UnmarshalText([]byte("42"))
	require.Error(t, err)
}

func TestMarshalPBOValid(t *testing.T) {
	o := PBOFrom(&_emptyObj)
	data, err := json.Marshal(o)
	require.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestMarshalPBOInvalid(t *testing.T) {
	// invalid values should be encoded as null
	null := NewPBO[*emptypb.Empty](nil, false)
	data, err := json.Marshal(null)
	require.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

func TestMarshalPBOText(t *testing.T) {
	o := PBOFrom(&_emptyObj)
	data, err := o.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func TestPBOScanValid(t *testing.T) {
	var o PBO[*emptypb.Empty]
	err := o.Scan([]byte(""))
	require.NoError(t, err)
	assert.True(t, proto.Equal(o.Object, &_emptyObj))
}

func TestPBOScanNull(t *testing.T) {
	var null PBO[*emptypb.Empty]
	err := null.Scan(nil)
	require.NoError(t, err)
	assert.False(t, null.Valid)
}
