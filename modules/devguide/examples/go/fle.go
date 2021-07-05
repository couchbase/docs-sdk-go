package main

import (
	"errors"
	"fmt"

	"github.com/couchbase/gocb/v2"
	gocbfieldcrypt "github.com/couchbase/gocbencryption/v2"
)

// tag::annotation[]
type PersonAddress struct {
	HouseName  string `json:"houseName" encrypted:"one"`
	StreetName string `json:"streetName"`
}

type Person struct {
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Password  string          `json:"password" encrypted:"one"`
	Addresses []PersonAddress `json:"address" encrypted:"two"`

	Phone string `json:"phone" encrypted:"two"`
}

// end::annotation[]

func main() {
	authenticator := gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	}

	// tag::keys[]
	keyB := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x3f, 0x3f, 0x3f, 0x3f, 0x3d, 0x3e, 0x3f,
	}
	key1 := gocbfieldcrypt.Key{
		ID:    "mykey",
		Bytes: keyB,
	}
	key2 := gocbfieldcrypt.Key{
		ID:    "myotherkey",
		Bytes: keyB,
	}

	// Create an insecure keyring and add two keys.
	keyring := gocbfieldcrypt.NewInsecureKeyring()
	keyring.Add(key1)
	keyring.Add(key2)
	// end::keys[]

	// tag::provider[]
	// Create a provider.
	// AES-256 authenticated with HMAC SHA-512. Requires a 64-byte key.
	provider := gocbfieldcrypt.NewAeadAes256CbcHmacSha512Provider(keyring)

	// Create the manager and add the providers.
	mgr := gocbfieldcrypt.NewDefaultCryptoManager(nil)

	// We need to create and then register encrypters.
	// The keyID here is used by the encrypter to lookup the key from the store when encrypting a document.
	// The key.ID returned from the store at encryption time is written into the data for the field to be encrypted.
	// The key ID that was written is then used on the decrypt side to find the corresponding key from the store.
	keyOneEncrypter := provider.EncrypterForKey(key1.ID)

	// We register the providers for both encryption and decryption.
	// The alias used here is the value which corresponds to the "encrypted" field annotation.
	err := mgr.RegisterEncrypter("one", keyOneEncrypter)
	if err != nil {
		panic(err)
	}

	err = mgr.RegisterEncrypter("two", provider.EncrypterForKey(key2.ID))
	if err != nil {
		panic(err)
	}

	// We don't need to add a default encryptor but if we do then any fields with an
	// empty encrypted tag will use this encryptor.
	err = mgr.DefaultEncrypter(keyOneEncrypter)
	if err != nil {
		panic(err)
	}

	// We only set one decrypter per algorithm.
	// The crypto manager will work out which decrypter to use based on the `alg` field embedded in the field data.
	// The decrypter will use the key embedded in the field data to determine which key to fetch from the key store for decryption.
	err = mgr.RegisterDecrypter(provider.Decrypter())
	if err != nil {
		panic(err)
	}
	// end::provider[]

	// tag::transcoder[]
	// Create our transcoder, not setting a base transcoder will cause it to fallback to the
	// SDK JSON transcoder.
	transcoder := gocbfieldcrypt.NewTranscoder(nil, mgr)

	// Register the encryption transcoder with the SDK.
	opts := gocb.ClusterOptions{
		Authenticator: authenticator,
		Transcoder:    transcoder,
	}
	cluster, err := gocb.Connect("localhost", opts)
	if err != nil {
		panic(err)
	}
	// end::transcoder[]

	// tag::upsert[]
	bucket := cluster.Bucket("travel-sample")
	collection := bucket.Scope("inventory").Collection("airport")

	person := Person{
		FirstName: "Barry",
		LastName:  "Sheen",
		Password:  "bang!",
		Addresses: []PersonAddress{
			{
				HouseName:  "my house",
				StreetName: "my street",
			},
			{
				HouseName:  "my other house",
				StreetName: "my other street",
			},
		},
		Phone: "123456",
	}

	_, err = collection.Upsert("p1", person, nil)
	if err != nil {
		panic(err)
	}
	// end::upsert[]

	// tag::getmap[]
	res, err := collection.Get("p1", nil)
	if err != nil {
		panic(err)
	}

	var resData map[string]interface{}
	err = res.Content(&resData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resData)
	// end::getmap[]

	// tag::getstr[]
	personRes, err := collection.Get("p1", nil)
	if err != nil {
		panic(err)
	}

	var personResData Person
	err = personRes.Content(&personResData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", personResData)
	// end::getstr[]
}

func oldAnnotation() {
	// tag::oldannotation[]
	mgr := gocbfieldcrypt.NewDefaultCryptoManager(&gocbfieldcrypt.DefaultCryptoManagerOptions{
		EncryptedFieldPrefix: "__crypt",
	})
	// end::oldannotation[]

	keyring := gocbfieldcrypt.NewInsecureKeyring()

	// tag::legacy[]
	// NewLegacyAes256CryptoDecrypter takes a function parameter so that the single decrypter can support multiple
	// keys. The function accepts a public key name and returns the corresponding private key name.
	decrypter := gocbfieldcrypt.NewLegacyAes256CryptoDecrypter(keyring, func(s string) (string, error) {
		if s == "mykey" {
			return "myhmackey", nil
		}

		return "", errors.New("unknown key")
	})
	err := mgr.RegisterDecrypter(decrypter)
	if err != nil {
		panic(err)
	}
	// end::legacy[]
}
