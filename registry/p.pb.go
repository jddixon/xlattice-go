// Code generated by protoc-gen-go.
// source: p.proto
// DO NOT EDIT!

package registry

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type XLRegMsg_Cmd int32

const (
	XLRegMsg_Hello      XLRegMsg_Cmd = 1
	XLRegMsg_HelloReply XLRegMsg_Cmd = 2
	XLRegMsg_Join       XLRegMsg_Cmd = 3
	XLRegMsg_IHave      XLRegMsg_Cmd = 4
	XLRegMsg_Get        XLRegMsg_Cmd = 6
	XLRegMsg_Bye        XLRegMsg_Cmd = 2
	XLRegMsg_Ack        XLRegMsg_Cmd = 4
	XLRegMsg_Error      XLRegMsg_Cmd = 5
)

var XLRegMsg_Cmd_name = map[int32]string{
	1: "Hello",
	2: "HelloReply",
	3: "Join",
	4: "IHave",
	6: "Get",
	// Duplicate value: 2: "Bye",
	// Duplicate value: 4: "Ack",
	5: "Error",
}
var XLRegMsg_Cmd_value = map[string]int32{
	"Hello":      1,
	"HelloReply": 2,
	"Join":       3,
	"IHave":      4,
	"Get":        6,
	"Bye":        2,
	"Ack":        4,
	"Error":      5,
}

func (x XLRegMsg_Cmd) Enum() *XLRegMsg_Cmd {
	p := new(XLRegMsg_Cmd)
	*p = x
	return p
}
func (x XLRegMsg_Cmd) String() string {
	return proto.EnumName(XLRegMsg_Cmd_name, int32(x))
}
func (x XLRegMsg_Cmd) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *XLRegMsg_Cmd) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(XLRegMsg_Cmd_value, data, "XLRegMsg_Cmd")
	if err != nil {
		return err
	}
	*x = XLRegMsg_Cmd(value)
	return nil
}

type XLRegMsg struct {
	Op               *XLRegMsg_Cmd `protobuf:"varint,1,opt,enum=registry.XLRegMsg_Cmd" json:"Op,omitempty"`
	MsgN             *uint64       `protobuf:"varint,2,opt" json:"MsgN,omitempty"`
	KeyIV            []byte        `protobuf:"bytes,3,opt" json:"KeyIV,omitempty"`
	Salt             []byte        `protobuf:"bytes,4,opt" json:"Salt,omitempty"`
	ID               []byte        `protobuf:"bytes,5,opt" json:"ID,omitempty"`
	CommsKey         []byte        `protobuf:"bytes,6,opt" json:"CommsKey,omitempty"`
	SigKey           []byte        `protobuf:"bytes,7,opt" json:"SigKey,omitempty"`
	MyEnd            *string       `protobuf:"bytes,8,opt" json:"MyEnd,omitempty"`
	ClusterID        []byte        `protobuf:"bytes,9,opt" json:"ClusterID,omitempty"`
	Payload          []byte        `protobuf:"bytes,10,opt" json:"Payload,omitempty"`
	ErrDesc          *string       `protobuf:"bytes,11,opt" json:"ErrDesc,omitempty"`
	YourMsgN         *uint64       `protobuf:"varint,12,opt" json:"YourMsgN,omitempty"`
	YourID           []byte        `protobuf:"bytes,13,opt" json:"YourID,omitempty"`
	Sig              []byte        `protobuf:"bytes,14,opt" json:"Sig,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *XLRegMsg) Reset()         { *m = XLRegMsg{} }
func (m *XLRegMsg) String() string { return proto.CompactTextString(m) }
func (*XLRegMsg) ProtoMessage()    {}

func (m *XLRegMsg) GetOp() XLRegMsg_Cmd {
	if m != nil && m.Op != nil {
		return *m.Op
	}
	return 0
}

func (m *XLRegMsg) GetMsgN() uint64 {
	if m != nil && m.MsgN != nil {
		return *m.MsgN
	}
	return 0
}

func (m *XLRegMsg) GetKeyIV() []byte {
	if m != nil {
		return m.KeyIV
	}
	return nil
}

func (m *XLRegMsg) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *XLRegMsg) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *XLRegMsg) GetCommsKey() []byte {
	if m != nil {
		return m.CommsKey
	}
	return nil
}

func (m *XLRegMsg) GetSigKey() []byte {
	if m != nil {
		return m.SigKey
	}
	return nil
}

func (m *XLRegMsg) GetMyEnd() string {
	if m != nil && m.MyEnd != nil {
		return *m.MyEnd
	}
	return ""
}

func (m *XLRegMsg) GetClusterID() []byte {
	if m != nil {
		return m.ClusterID
	}
	return nil
}

func (m *XLRegMsg) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *XLRegMsg) GetErrDesc() string {
	if m != nil && m.ErrDesc != nil {
		return *m.ErrDesc
	}
	return ""
}

func (m *XLRegMsg) GetYourMsgN() uint64 {
	if m != nil && m.YourMsgN != nil {
		return *m.YourMsgN
	}
	return 0
}

func (m *XLRegMsg) GetYourID() []byte {
	if m != nil {
		return m.YourID
	}
	return nil
}

func (m *XLRegMsg) GetSig() []byte {
	if m != nil {
		return m.Sig
	}
	return nil
}

func init() {
	proto.RegisterEnum("registry.XLRegMsg_Cmd", XLRegMsg_Cmd_name, XLRegMsg_Cmd_value)
}
