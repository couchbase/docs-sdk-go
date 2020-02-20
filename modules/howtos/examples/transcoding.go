package main

import (
	"errors"

	"github.com/couchbase/gocb/v2"
	ffjson "github.com/pquerna/ffjson/ffjson"
	"github.com/tinylib/msgp/msgp"
)

func main() {
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			"Administrator",
			"password",
		},
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}

	bucket := cluster.Bucket("bucket-name")

	collection := bucket.DefaultCollection()

	rawString(collection)
	rawBinary(collection)
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func rawJSONMarshal(collection *gocb.Collection) {
	// #tag::rawjsonmarshal[]
	user := User{Name: "John Smith", Age: 27}
	transcoder := gocb.NewRawJSONTranscoder()

	b, err := ffjson.Marshal(user)
	if err != nil {
		panic(err)
	}

	_, err = collection.Upsert("john-smith", b, &gocb.UpsertOptions{
		Transcoder: transcoder,
	})
	if err != nil {
		panic(err)
	}
	// #end::rawjsonmarshal[]
}

func rawJSONUnmarshal(collection *gocb.Collection) {
	// #tag::rawjsonunmarshal[]
	transcoder := gocb.NewRawJSONTranscoder()

	getRes, err := collection.Get("john-smith", &gocb.GetOptions{
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

	var user User
	err = ffjson.Unmarshal(returned, &user)
	if err != nil {
		panic(err)
	}
	// #end::rawjsonunmarshal[]
}

func rawString(collection *gocb.Collection) {
	// #tag::rawstring[]
	input := "hello world"
	transcoder := gocb.NewRawStringTranscoder()

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

	var returned string
	err = getRes.Content(&returned)
	if err != nil {
		panic(err)
	}
	// #end::rawstring[]
}

func rawBinary(collection *gocb.Collection) {
	// #tag::rawbinary[]
	input := []byte("hello world")
	transcoder := gocb.NewRawBinaryTranscoder()

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
	// #end::rawbinary[]
}

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
