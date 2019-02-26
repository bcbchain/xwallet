package types

import (
	"bytes"
	"encoding/json"
)

type Validators []Validator

func (v Validators) Len() int {
	return len(v)
}

func (v Validators) Less(i, j int) bool {
	return bytes.Compare(v[i].PubKey, v[j].PubKey) <= 0
}

func (v Validators) Swap(i, j int) {
	v1 := v[i]
	v[i] = v[j]
	v[j] = v1
}

func ValidatorsString(vs Validators) string {
	s := make([]validatorPretty, len(vs))
	for i, v := range vs {
		s[i] = validatorPretty(v)
	}
	b, err := json.Marshal(s)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

type validatorPretty struct {
	PubKey		[]byte	`json:"pub_key"`
	Power		uint64	`json:"power"`
	RewardAddr	string	`json:"reward_addr"`
	Name		string	`json:"name"`
}