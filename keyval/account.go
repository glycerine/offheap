package keyval

type Account struct {
	Id               int64  `protobuf:"varint,1,opt,name=id" json:"id" capid:"0"`
	Dty              int64  `protobuf:"varint,2,opt,name=dty" json:"dty" capid:"1"`
	AcctId           string `protobuf:"bytes,3,opt" json:"AcctId" capid:"2"`
	OpenedFromIP     string `protobuf:"bytes,4,opt" json:"OpenedFromIP" capid:"3"`
	Name             string `protobuf:"bytes,5,opt" json:"Name" capid:"4"`
	Email            string `protobuf:"bytes,6,opt" json:"Email" capid:"5"`
	Disabled         int64  `protobuf:"varint,7,opt" json:"Disabled" capid:"6"`
	XXX_unrecognized []byte `json:"-" capid:"skip"`
}
