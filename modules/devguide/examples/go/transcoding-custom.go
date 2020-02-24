package main

import (
	"errors"

	"github.com/couchbase/gocb/v2"
	"github.com/tinylib/msgp/msgp"
)

// This is an example of a custom transcoder that uses Message Pack as the data format.

// #tag::msgpack[]
const customFlags = (0x02 << 24) | ('M' << 16) | ('P' << 8) | ('K' << 0)

type MsgPackTranscoder struct {
}

func (t *MsgPackTranscoder) Encode(value interface{}) ([]byte, uint32, error) {
	msgPckVal, ok := value.(msgp.Marshaler)
	if !ok {
		return nil, 0, errors.New("MsgPackTranscoder only supports types that satisfy msgp.Marshaler")
	}

	data, err := msgPckVal.MarshalMsg(nil)
	if err != nil {
		return nil, 0, err
	}

	return data, customFlags, nil
}

func (t *MsgPackTranscoder) Decode(bytes []byte, flags uint32, out interface{}) error {
	if flags != customFlags {
		return errors.New("unexpected expectedFlags value")
	}

	msgPckVal, ok := out.(msgp.Unmarshaler)
	if !ok {
		return errors.New("MsgPackTranscoder only supports types that satisfy msgp.Unmarshaler")
	}

	var err error
	out, err = msgPckVal.UnmarshalMsg(bytes)
	if err != nil {
		return err
	}

	return nil
}

// #end::msgpack[]

// input here must satisfy the msgp package interfaces
func msgpackTranscode(collection *gocb.Collection, input interface{}) {
	// #tag::msgpack-transcode[]
	transcoder := &MsgPackTranscoder{}

	_, err := collection.Upsert("key", input, &gocb.UpsertOptions{
		Transcoder: transcoder,
	})
	if err != nil {
		panic(err)
	}

	getRes, err := collection.Get("key", &gocb.GetOptions{
		Transcoder: transcoder,
	})
	if err != nil {
		panic(err)
	}

	var returned []byte
	err = getRes.Content(&returned)
	if err != nil {
		panic(err)
	}
	// #end::msgpack-transcode[]
}
