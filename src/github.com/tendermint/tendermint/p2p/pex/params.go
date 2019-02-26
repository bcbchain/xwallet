package pex

import "time"

const (
	needAddressThreshold	= 1000

	dumpAddressInterval	= time.Minute * 2

	oldBucketSize	= 64

	oldBucketCount	= 64

	newBucketSize	= 64

	newBucketCount	= 256

	oldBucketsPerGroup	= 4

	newBucketsPerGroup	= 32

	maxNewBucketsPerAddress	= 4

	numMissingDays	= 7

	numRetries	= 3

	maxFailures	= 10

	minBadDays	= 7

	getSelectionPercent	= 23

	minGetSelection	= 32

	maxGetSelection	= 250
)
