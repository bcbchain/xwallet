package types

import (
	"encoding/json"
	cfg "github.com/tendermint/tendermint/config"
	"io/ioutil"
	"time"

	"github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
	"bcbchain.io/kms"
)

type GenesisValidator struct {
	RewardAddr	string		`json:"reward_addr"`
	PubKey		crypto.PubKey	`json:"pub_key,omitempty"`
	Power		int64		`json:"power"`
	Name		string		`json:"name"`
}

type GenesisDoc struct {
	GenesisTime	time.Time		`json:"genesis_time"`
	ChainID		string			`json:"chain_id"`
	ConsensusParams	*ConsensusParams	`json:"consensus_params,omitempty"`
	Validators	[]GenesisValidator	`json:"validators"`
	AppHash		cmn.HexBytes		`json:"app_hash"`
	AppStateJSON	json.RawMessage		`json:"app_state,omitempty"`
	AppOptions	json.RawMessage		`json:"app_options,omitempty"`
}

func (genDoc *GenesisDoc) AppState() json.RawMessage {
	if len(genDoc.AppOptions) > 0 {
		return genDoc.AppOptions
	}
	return genDoc.AppStateJSON
}

func (genDoc *GenesisDoc) SaveAs(file string) error {
	genDocBytes, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}
	return cmn.WriteFile(file, genDocBytes, 0644)
}

func (genDoc *GenesisDoc) ValidatorHash() []byte {
	vals := make([]*Validator, len(genDoc.Validators))
	for i, v := range genDoc.Validators {
		if v.Power < 0 {
			v.Power = 0
		}
		vals[i] = NewValidator(v.PubKey, uint64(v.Power), v.RewardAddr, v.Name)
	}
	vset := NewValidatorSet(vals)
	return vset.Hash()
}

func (genDoc *GenesisDoc) ValidateAndComplete() error {

	if genDoc.ChainID == "" {
		return cmn.NewError("Genesis doc must include non-empty chain_id")
	}

	if genDoc.ConsensusParams == nil {
		genDoc.ConsensusParams = DefaultConsensusParams()
	} else {
		if err := genDoc.ConsensusParams.Validate(); err != nil {
			return err
		}
	}

	if len(genDoc.Validators) == 0 {
		return cmn.NewError("The genesis file must have at least one validator")
	}

	for _, v := range genDoc.Validators {
		if v.Power == 0 {
			return cmn.NewError("The genesis file cannot contain validators with no voting power: %v", v)
		}
	}

	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = time.Now()
	}

	return nil
}

func GenesisDocFromJSON(jsonBlob []byte) (*GenesisDoc, error) {
	genDoc := GenesisDoc{}
	err := cdc.UnmarshalJSON(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	return &genDoc, err
}

func GenesisDocFromFile(config *cfg.Config) (*GenesisDoc, error) {
	genesisFile := config.GenesisFile()

	signatureFile := genesisFile[0:len(genesisFile)-5] + "-signature.json"
	_, err := kms.VerifyFileSign(genesisFile, signatureFile)
	if err != nil {
		return nil, cmn.ErrorWrap(err, cmn.Fmt("Genesis file verify failed, %v", err.Error()))
	}

	jsonBlob, err := ioutil.ReadFile(genesisFile)
	if err != nil {
		return nil, cmn.ErrorWrap(err, "Couldn't read GenesisDoc file")
	}
	genDoc, err := GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, cmn.ErrorWrap(err, cmn.Fmt("Error reading GenesisDoc at %v", config.GenesisFile()))
	}
	validators := ValidatorsFromFile(*genDoc, config.ValidatorsFile())
	genDoc.Validators = *validators
	return genDoc, nil
}

func ValidatorsFromFile(genDoc GenesisDoc, validatorsFile string) *[]GenesisValidator {
	jsonBlob, err := ioutil.ReadFile(validatorsFile)
	if err != nil {
		panic("Couldn't read Validators file")
	}
	validators := make([]GenesisValidator, 0)
	err = cdc.UnmarshalJSON(jsonBlob, &validators)
	if err != nil {
		panic(cmn.Fmt("Error reading Validators at %v", validatorsFile))
	}
	genValidators := genDoc.Validators
	flag := false
	for _, v := range genValidators {
		if !inSlice(v, validators) {
			flag = true
		}
	}
	if flag || len(genValidators) != len(validators) {
		panic("genesis.json & validators.json doesn't match!")
	}
	return &validators
}

func inSlice(a GenesisValidator, list []GenesisValidator) bool {
	for _, b := range list {
		if a.RewardAddr == b.RewardAddr && a.Name == b.Name && a.Power == b.Power {
			return true
		}
	}
	return false
}
