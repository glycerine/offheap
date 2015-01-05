package keyval

import (
	"fmt"
	"io"

	capn "github.com/glycerine/go-capnproto"
)

func (s *Account) Save(w io.Writer) {
	seg := capn.NewBuffer(nil)
	AccountGoToCapn(seg, s)
	seg.WriteTo(w)
}

func (s *Account) Load(r io.Reader) {
	capMsg, err := capn.ReadFromStream(r, nil)
	if err != nil {
		panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
	}
	z := ReadRootAccountCapn(capMsg)
	AccountCapnToGo(z, s)
}

func AccountCapnToGo(src AccountCapn, dest *Account) *Account {
	if dest == nil {
		dest = &Account{}
	}
	dest.Id = int64(src.Id())
	dest.Dty = int64(src.Dty())
	dest.AcctId = src.AcctId()
	dest.OpenedFromIP = src.OpenedFromIP()
	dest.Name = src.Name()
	dest.Email = src.Email()
	dest.Disabled = int64(src.Disabled())

	return dest
}

func AccountGoToCapn(seg *capn.Segment, src *Account) AccountCapn {
	dest := AutoNewAccountCapn(seg)
	dest.SetId(src.Id)
	dest.SetDty(src.Dty)
	dest.SetAcctId(src.AcctId)
	dest.SetOpenedFromIP(src.OpenedFromIP)
	dest.SetName(src.Name)
	dest.SetEmail(src.Email)
	dest.SetDisabled(src.Disabled)

	return dest
}

func (s *AcctId) Save(w io.Writer) {
	seg := capn.NewBuffer(nil)
	AcctIdGoToCapn(seg, s)
	seg.WriteTo(w)
}

func (s *AcctId) Load(r io.Reader) {
	capMsg, err := capn.ReadFromStream(r, nil)
	if err != nil {
		panic(fmt.Errorf("capn.ReadFromStream error: %s", err))
	}
	z := ReadRootAcctIdCapn(capMsg)
	AcctIdCapnToGo(z, s)
}

func AcctIdCapnToGo(src AcctIdCapn, dest *AcctId) *AcctId {
	if dest == nil {
		dest = &AcctId{}
	}
	dest.AcctId = src.AcctId()

	return dest
}

func AcctIdGoToCapn(seg *capn.Segment, src *AcctId) AcctIdCapn {
	dest := AutoNewAcctIdCapn(seg)
	dest.SetAcctId(src.AcctId)

	return dest
}
