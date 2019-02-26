package keys

import (
	"github.com/tendermint/go-crypto"
)

type Keybase interface {
	Sign(name, passphrase string, msg []byte) (crypto.Signature, crypto.PubKey, error)

	Create(name, passphrase string, algo CryptoAlgo) (info Info, seed string, err error)

	Recover(name, passphrase, seedphrase string) (info Info, erro error)
	List() ([]Info, error)
	Get(name string) (Info, error)
	Update(name, oldpass, newpass string) error
	Delete(name, passphrase string) error

	Import(name string, armor string) (err error)
	Export(name string) (armor string, err error)
}

type Info struct {
	Name		string		`json:"name"`
	PubKey		crypto.PubKey	`json:"pubkey"`
	PrivKeyArmor	string		`json:"privkey.armor"`
}

func newInfo(name string, pub crypto.PubKey, privArmor string) Info {
	return Info{
		Name:		name,
		PubKey:		pub,
		PrivKeyArmor:	privArmor,
	}
}

func (i Info) Address() string {
	return i.PubKey.Address()
}

func (i Info) bytes() []byte {
	bz, err := cdc.MarshalBinaryBare(i)
	if err != nil {
		panic(err)
	}
	return bz
}

func readInfo(bz []byte) (info Info, err error) {
	err = cdc.UnmarshalBinaryBare(bz, &info)
	return
}
