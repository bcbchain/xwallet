package pex

import (
	"time"

	"github.com/tendermint/tendermint/p2p"
)

type knownAddress struct {
	Addr		*p2p.NetAddress	`json:"addr"`
	Src		*p2p.NetAddress	`json:"src"`
	Attempts	int32		`json:"attempts"`
	LastAttempt	time.Time	`json:"last_attempt"`
	LastSuccess	time.Time	`json:"last_success"`
	BucketType	byte		`json:"bucket_type"`
	Buckets		[]int		`json:"buckets"`
}

func newKnownAddress(addr *p2p.NetAddress, src *p2p.NetAddress) *knownAddress {
	return &knownAddress{
		Addr:		addr,
		Src:		src,
		Attempts:	0,
		LastAttempt:	time.Now(),
		BucketType:	bucketTypeNew,
		Buckets:	nil,
	}
}

func (ka *knownAddress) ID() p2p.ID {
	return ka.Addr.ID
}

func (ka *knownAddress) copy() *knownAddress {
	return &knownAddress{
		Addr:		ka.Addr,
		Src:		ka.Src,
		Attempts:	ka.Attempts,
		LastAttempt:	ka.LastAttempt,
		LastSuccess:	ka.LastSuccess,
		BucketType:	ka.BucketType,
		Buckets:	ka.Buckets,
	}
}

func (ka *knownAddress) isOld() bool {
	return ka.BucketType == bucketTypeOld
}

func (ka *knownAddress) isNew() bool {
	return ka.BucketType == bucketTypeNew
}

func (ka *knownAddress) markAttempt() {
	now := time.Now()
	ka.LastAttempt = now
	ka.Attempts++
}

func (ka *knownAddress) markGood() {
	now := time.Now()
	ka.LastAttempt = now
	ka.Attempts = 0
	ka.LastSuccess = now
}

func (ka *knownAddress) addBucketRef(bucketIdx int) int {
	for _, bucket := range ka.Buckets {
		if bucket == bucketIdx {

			return -1
		}
	}
	ka.Buckets = append(ka.Buckets, bucketIdx)
	return len(ka.Buckets)
}

func (ka *knownAddress) removeBucketRef(bucketIdx int) int {
	buckets := []int{}
	for _, bucket := range ka.Buckets {
		if bucket != bucketIdx {
			buckets = append(buckets, bucket)
		}
	}
	if len(buckets) != len(ka.Buckets)-1 {

		return -1
	}
	ka.Buckets = buckets
	return len(ka.Buckets)
}

func (ka *knownAddress) isBad() bool {

	if ka.BucketType == bucketTypeOld {
		return false
	}

	if ka.LastAttempt.Before(time.Now().Add(-1 * time.Minute)) {
		return false
	}

	if ka.LastAttempt.After(time.Now().Add(-1 * numMissingDays * time.Hour * 24)) {
		return true
	}

	if ka.LastSuccess.IsZero() && ka.Attempts >= numRetries {
		return true
	}

	if ka.LastSuccess.Before(time.Now().Add(-1*minBadDays*time.Hour*24)) &&
		ka.Attempts >= maxFailures {
		return true
	}

	return false
}
