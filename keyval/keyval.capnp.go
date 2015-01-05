package keyval

// AUTO GENERATED - DO NOT EDIT

import (
	"unsafe"

	C "github.com/glycerine/go-capnproto"
)

type AccountCapn C.Struct

func NewAccountCapn(s *C.Segment) AccountCapn      { return AccountCapn(s.NewStruct(24, 4)) }
func NewRootAccountCapn(s *C.Segment) AccountCapn  { return AccountCapn(s.NewRootStruct(24, 4)) }
func AutoNewAccountCapn(s *C.Segment) AccountCapn  { return AccountCapn(s.NewStructAR(24, 4)) }
func ReadRootAccountCapn(s *C.Segment) AccountCapn { return AccountCapn(s.Root(0).ToStruct()) }
func (s AccountCapn) Id() int64                    { return int64(C.Struct(s).Get64(0)) }
func (s AccountCapn) SetId(v int64)                { C.Struct(s).Set64(0, uint64(v)) }
func (s AccountCapn) Dty() int64                   { return int64(C.Struct(s).Get64(8)) }
func (s AccountCapn) SetDty(v int64)               { C.Struct(s).Set64(8, uint64(v)) }
func (s AccountCapn) AcctId() string               { return C.Struct(s).GetObject(0).ToText() }
func (s AccountCapn) SetAcctId(v string)           { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s AccountCapn) OpenedFromIP() string         { return C.Struct(s).GetObject(1).ToText() }
func (s AccountCapn) SetOpenedFromIP(v string)     { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s AccountCapn) Name() string                 { return C.Struct(s).GetObject(2).ToText() }
func (s AccountCapn) SetName(v string)             { C.Struct(s).SetObject(2, s.Segment.NewText(v)) }
func (s AccountCapn) Email() string                { return C.Struct(s).GetObject(3).ToText() }
func (s AccountCapn) SetEmail(v string)            { C.Struct(s).SetObject(3, s.Segment.NewText(v)) }
func (s AccountCapn) Disabled() int64              { return int64(C.Struct(s).Get64(16)) }
func (s AccountCapn) SetDisabled(v int64)          { C.Struct(s).Set64(16, uint64(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s AccountCapn) MarshalJSON() (bs []byte, err error) { return }

type AccountCapn_List C.PointerList

func NewAccountCapnList(s *C.Segment, sz int) AccountCapn_List {
	return AccountCapn_List(s.NewCompositeList(24, 4, sz))
}
func (s AccountCapn_List) Len() int             { return C.PointerList(s).Len() }
func (s AccountCapn_List) At(i int) AccountCapn { return AccountCapn(C.PointerList(s).At(i).ToStruct()) }
func (s AccountCapn_List) ToArray() []AccountCapn {
	return *(*[]AccountCapn)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
func (s AccountCapn_List) Set(i int, item AccountCapn) { C.PointerList(s).Set(i, C.Object(item)) }

type AcctIdCapn C.Struct

func NewAcctIdCapn(s *C.Segment) AcctIdCapn      { return AcctIdCapn(s.NewStruct(0, 1)) }
func NewRootAcctIdCapn(s *C.Segment) AcctIdCapn  { return AcctIdCapn(s.NewRootStruct(0, 1)) }
func AutoNewAcctIdCapn(s *C.Segment) AcctIdCapn  { return AcctIdCapn(s.NewStructAR(0, 1)) }
func ReadRootAcctIdCapn(s *C.Segment) AcctIdCapn { return AcctIdCapn(s.Root(0).ToStruct()) }
func (s AcctIdCapn) AcctId() string              { return C.Struct(s).GetObject(0).ToText() }
func (s AcctIdCapn) SetAcctId(v string)          { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s AcctIdCapn) MarshalJSON() (bs []byte, err error) { return }

type AcctIdCapn_List C.PointerList

func NewAcctIdCapnList(s *C.Segment, sz int) AcctIdCapn_List {
	return AcctIdCapn_List(s.NewCompositeList(0, 1, sz))
}
func (s AcctIdCapn_List) Len() int            { return C.PointerList(s).Len() }
func (s AcctIdCapn_List) At(i int) AcctIdCapn { return AcctIdCapn(C.PointerList(s).At(i).ToStruct()) }
func (s AcctIdCapn_List) ToArray() []AcctIdCapn {
	return *(*[]AcctIdCapn)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
func (s AcctIdCapn_List) Set(i int, item AcctIdCapn) { C.PointerList(s).Set(i, C.Object(item)) }
