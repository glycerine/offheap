package offheap

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import "github.com/tinylib/msgp/msgp"

// DecodeMsg implements msgp.Decodable
func (z *Cell) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zbai uint32
	zbai, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zbai > 0 {
		zbai--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "UnHashedKey":
			z.UnHashedKey, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "ByteKey":
			err = dc.ReadExactBytes(z.ByteKey[:])
			if err != nil {
				return
			}
		case "Value":
			err = dc.ReadExactBytes(z.Value[:])
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Cell) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "UnHashedKey"
	err = en.Append(0x83, 0xab, 0x55, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x65, 0x64, 0x4b, 0x65, 0x79)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.UnHashedKey)
	if err != nil {
		return
	}
	// write "ByteKey"
	err = en.Append(0xa7, 0x42, 0x79, 0x74, 0x65, 0x4b, 0x65, 0x79)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.ByteKey[:])
	if err != nil {
		return
	}
	// write "Value"
	err = en.Append(0xa5, 0x56, 0x61, 0x6c, 0x75, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Value[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Cell) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "UnHashedKey"
	o = append(o, 0x83, 0xab, 0x55, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x65, 0x64, 0x4b, 0x65, 0x79)
	o = msgp.AppendUint64(o, z.UnHashedKey)
	// string "ByteKey"
	o = append(o, 0xa7, 0x42, 0x79, 0x74, 0x65, 0x4b, 0x65, 0x79)
	o = msgp.AppendBytes(o, z.ByteKey[:])
	// string "Value"
	o = append(o, 0xa5, 0x56, 0x61, 0x6c, 0x75, 0x65)
	o = msgp.AppendBytes(o, z.Value[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Cell) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcmr uint32
	zcmr, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcmr > 0 {
		zcmr--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "UnHashedKey":
			z.UnHashedKey, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "ByteKey":
			bts, err = msgp.ReadExactBytes(bts, z.ByteKey[:])
			if err != nil {
				return
			}
		case "Value":
			bts, err = msgp.ReadExactBytes(bts, z.Value[:])
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Cell) Msgsize() (s int) {
	s = 1 + 12 + msgp.Uint64Size + 8 + msgp.ArrayHeaderSize + (64 * (msgp.ByteSize)) + 6 + msgp.ArrayHeaderSize + (56 * (msgp.ByteSize))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *HashTable) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zajw uint32
	zajw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zajw > 0 {
		zajw--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "MagicNumber":
			z.MagicNumber, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "CellSz":
			z.CellSz, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "ArraySize":
			z.ArraySize, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Population":
			z.Population, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "ZeroUsed":
			z.ZeroUsed, err = dc.ReadBool()
			if err != nil {
				return
			}
		case "ZeroCell":
			err = z.ZeroCell.DecodeMsg(dc)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *HashTable) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 6
	// write "MagicNumber"
	err = en.Append(0x86, 0xab, 0x4d, 0x61, 0x67, 0x69, 0x63, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.MagicNumber)
	if err != nil {
		return
	}
	// write "CellSz"
	err = en.Append(0xa6, 0x43, 0x65, 0x6c, 0x6c, 0x53, 0x7a)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.CellSz)
	if err != nil {
		return
	}
	// write "ArraySize"
	err = en.Append(0xa9, 0x41, 0x72, 0x72, 0x61, 0x79, 0x53, 0x69, 0x7a, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.ArraySize)
	if err != nil {
		return
	}
	// write "Population"
	err = en.Append(0xaa, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.Population)
	if err != nil {
		return
	}
	// write "ZeroUsed"
	err = en.Append(0xa8, 0x5a, 0x65, 0x72, 0x6f, 0x55, 0x73, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.ZeroUsed)
	if err != nil {
		return
	}
	// write "ZeroCell"
	err = en.Append(0xa8, 0x5a, 0x65, 0x72, 0x6f, 0x43, 0x65, 0x6c, 0x6c)
	if err != nil {
		return err
	}
	err = z.ZeroCell.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *HashTable) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "MagicNumber"
	o = append(o, 0x86, 0xab, 0x4d, 0x61, 0x67, 0x69, 0x63, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	o = msgp.AppendInt(o, z.MagicNumber)
	// string "CellSz"
	o = append(o, 0xa6, 0x43, 0x65, 0x6c, 0x6c, 0x53, 0x7a)
	o = msgp.AppendUint64(o, z.CellSz)
	// string "ArraySize"
	o = append(o, 0xa9, 0x41, 0x72, 0x72, 0x61, 0x79, 0x53, 0x69, 0x7a, 0x65)
	o = msgp.AppendUint64(o, z.ArraySize)
	// string "Population"
	o = append(o, 0xaa, 0x50, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Population)
	// string "ZeroUsed"
	o = append(o, 0xa8, 0x5a, 0x65, 0x72, 0x6f, 0x55, 0x73, 0x65, 0x64)
	o = msgp.AppendBool(o, z.ZeroUsed)
	// string "ZeroCell"
	o = append(o, 0xa8, 0x5a, 0x65, 0x72, 0x6f, 0x43, 0x65, 0x6c, 0x6c)
	o, err = z.ZeroCell.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HashTable) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zwht uint32
	zwht, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zwht > 0 {
		zwht--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "MagicNumber":
			z.MagicNumber, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "CellSz":
			z.CellSz, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "ArraySize":
			z.ArraySize, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Population":
			z.Population, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "ZeroUsed":
			z.ZeroUsed, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		case "ZeroCell":
			bts, err = z.ZeroCell.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *HashTable) Msgsize() (s int) {
	s = 1 + 12 + msgp.IntSize + 7 + msgp.Uint64Size + 10 + msgp.Uint64Size + 11 + msgp.Uint64Size + 9 + msgp.BoolSize + 9 + z.ZeroCell.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Iterator) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxhx uint32
	zxhx, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxhx > 0 {
		zxhx--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Tab":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Tab = nil
			} else {
				if z.Tab == nil {
					z.Tab = new(HashTable)
				}
				err = z.Tab.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Pos":
			z.Pos, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "Cur":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.Cur = nil
			} else {
				if z.Cur == nil {
					z.Cur = new(Cell)
				}
				var zlqf uint32
				zlqf, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zlqf > 0 {
					zlqf--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "UnHashedKey":
						z.Cur.UnHashedKey, err = dc.ReadUint64()
						if err != nil {
							return
						}
					case "ByteKey":
						err = dc.ReadExactBytes(z.Cur.ByteKey[:])
						if err != nil {
							return
						}
					case "Value":
						err = dc.ReadExactBytes(z.Cur.Value[:])
						if err != nil {
							return
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Iterator) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Tab"
	err = en.Append(0x83, 0xa3, 0x54, 0x61, 0x62)
	if err != nil {
		return err
	}
	if z.Tab == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.Tab.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Pos"
	err = en.Append(0xa3, 0x50, 0x6f, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteInt64(z.Pos)
	if err != nil {
		return
	}
	// write "Cur"
	err = en.Append(0xa3, 0x43, 0x75, 0x72)
	if err != nil {
		return err
	}
	if z.Cur == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		// map header, size 3
		// write "UnHashedKey"
		err = en.Append(0x83, 0xab, 0x55, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x65, 0x64, 0x4b, 0x65, 0x79)
		if err != nil {
			return err
		}
		err = en.WriteUint64(z.Cur.UnHashedKey)
		if err != nil {
			return
		}
		// write "ByteKey"
		err = en.Append(0xa7, 0x42, 0x79, 0x74, 0x65, 0x4b, 0x65, 0x79)
		if err != nil {
			return err
		}
		err = en.WriteBytes(z.Cur.ByteKey[:])
		if err != nil {
			return
		}
		// write "Value"
		err = en.Append(0xa5, 0x56, 0x61, 0x6c, 0x75, 0x65)
		if err != nil {
			return err
		}
		err = en.WriteBytes(z.Cur.Value[:])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Iterator) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Tab"
	o = append(o, 0x83, 0xa3, 0x54, 0x61, 0x62)
	if z.Tab == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Tab.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Pos"
	o = append(o, 0xa3, 0x50, 0x6f, 0x73)
	o = msgp.AppendInt64(o, z.Pos)
	// string "Cur"
	o = append(o, 0xa3, 0x43, 0x75, 0x72)
	if z.Cur == nil {
		o = msgp.AppendNil(o)
	} else {
		// map header, size 3
		// string "UnHashedKey"
		o = append(o, 0x83, 0xab, 0x55, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x65, 0x64, 0x4b, 0x65, 0x79)
		o = msgp.AppendUint64(o, z.Cur.UnHashedKey)
		// string "ByteKey"
		o = append(o, 0xa7, 0x42, 0x79, 0x74, 0x65, 0x4b, 0x65, 0x79)
		o = msgp.AppendBytes(o, z.Cur.ByteKey[:])
		// string "Value"
		o = append(o, 0xa5, 0x56, 0x61, 0x6c, 0x75, 0x65)
		o = msgp.AppendBytes(o, z.Cur.Value[:])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Iterator) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zdaf uint32
	zdaf, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zdaf > 0 {
		zdaf--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Tab":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Tab = nil
			} else {
				if z.Tab == nil {
					z.Tab = new(HashTable)
				}
				bts, err = z.Tab.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Pos":
			z.Pos, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "Cur":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Cur = nil
			} else {
				if z.Cur == nil {
					z.Cur = new(Cell)
				}
				var zpks uint32
				zpks, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for zpks > 0 {
					zpks--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "UnHashedKey":
						z.Cur.UnHashedKey, bts, err = msgp.ReadUint64Bytes(bts)
						if err != nil {
							return
						}
					case "ByteKey":
						bts, err = msgp.ReadExactBytes(bts, z.Cur.ByteKey[:])
						if err != nil {
							return
						}
					case "Value":
						bts, err = msgp.ReadExactBytes(bts, z.Cur.Value[:])
						if err != nil {
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Iterator) Msgsize() (s int) {
	s = 1 + 4
	if z.Tab == nil {
		s += msgp.NilSize
	} else {
		s += z.Tab.Msgsize()
	}
	s += 4 + msgp.Int64Size + 4
	if z.Cur == nil {
		s += msgp.NilSize
	} else {
		s += 1 + 12 + msgp.Uint64Size + 8 + msgp.ArrayHeaderSize + (64 * (msgp.ByteSize)) + 6 + msgp.ArrayHeaderSize + (56 * (msgp.ByteSize))
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Key_t) DecodeMsg(dc *msgp.Reader) (err error) {
	err = dc.ReadExactBytes(z[:])
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Key_t) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes(z[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Key_t) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, z[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Key_t) UnmarshalMsg(bts []byte) (o []byte, err error) {
	bts, err = msgp.ReadExactBytes(bts, z[:])
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Key_t) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (64 * (msgp.ByteSize))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Val_t) DecodeMsg(dc *msgp.Reader) (err error) {
	err = dc.ReadExactBytes(z[:])
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Val_t) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes(z[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Val_t) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, z[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Val_t) UnmarshalMsg(bts []byte) (o []byte, err error) {
	bts, err = msgp.ReadExactBytes(bts, z[:])
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Val_t) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (56 * (msgp.ByteSize))
	return
}
