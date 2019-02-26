package pex

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/p2p"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	bucketTypeNew	= 0x01
	bucketTypeOld	= 0x02
)

type AddrBook interface {
	cmn.Service

	AddOurAddress(*p2p.NetAddress)

	OurAddress(*p2p.NetAddress) bool

	AddAddress(addr *p2p.NetAddress, src *p2p.NetAddress) error
	RemoveAddress(*p2p.NetAddress)

	HasAddress(*p2p.NetAddress) bool

	NeedMoreAddrs() bool

	PickAddress(biasTowardsNewAddrs int) *p2p.NetAddress

	MarkGood(*p2p.NetAddress)
	MarkAttempt(*p2p.NetAddress)
	MarkBad(*p2p.NetAddress)

	IsGood(*p2p.NetAddress) bool

	GetSelection() []*p2p.NetAddress

	GetSelectionWithBias(biasTowardsNewAddrs int) []*p2p.NetAddress

	ListOfKnownAddresses() []*knownAddress

	Save()
}

var _ AddrBook = (*addrBook)(nil)

type addrBook struct {
	cmn.BaseService

	filePath		string
	routabilityStrict	bool
	key			string

	mtx		sync.Mutex
	rand		*cmn.Rand
	ourAddrs	map[string]struct{}
	addrLookup	map[p2p.ID]*knownAddress
	bucketsOld	[]map[string]*knownAddress
	bucketsNew	[]map[string]*knownAddress
	nOld		int
	nNew		int

	wg	sync.WaitGroup
}

func NewAddrBook(filePath string, routabilityStrict bool) *addrBook {
	am := &addrBook{
		rand:			cmn.NewRand(),
		ourAddrs:		make(map[string]struct{}),
		addrLookup:		make(map[p2p.ID]*knownAddress),
		filePath:		filePath,
		routabilityStrict:	routabilityStrict,
	}
	am.init()
	am.BaseService = *cmn.NewBaseService(nil, "AddrBook", am)
	return am
}

func (a *addrBook) init() {
	a.key = crypto.CRandHex(24)

	a.bucketsNew = make([]map[string]*knownAddress, newBucketCount)
	for i := range a.bucketsNew {
		a.bucketsNew[i] = make(map[string]*knownAddress)
	}

	a.bucketsOld = make([]map[string]*knownAddress, oldBucketCount)
	for i := range a.bucketsOld {
		a.bucketsOld[i] = make(map[string]*knownAddress)
	}
}

func (a *addrBook) OnStart() error {
	if err := a.BaseService.OnStart(); err != nil {
		return err
	}
	a.loadFromFile(a.filePath)

	a.wg.Add(1)
	go a.saveRoutine()

	return nil
}

func (a *addrBook) OnStop() {
	a.BaseService.OnStop()
}

func (a *addrBook) Wait() {
	a.wg.Wait()
}

func (a *addrBook) FilePath() string {
	return a.filePath
}

func (a *addrBook) AddOurAddress(addr *p2p.NetAddress) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.Logger.Info("Add our address to book", "addr", addr)
	a.ourAddrs[addr.String()] = struct{}{}
}

func (a *addrBook) OurAddress(addr *p2p.NetAddress) bool {
	a.mtx.Lock()
	_, ok := a.ourAddrs[addr.String()]
	a.mtx.Unlock()
	return ok
}

func (a *addrBook) AddAddress(addr *p2p.NetAddress, src *p2p.NetAddress) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.addAddress(addr, src)
}

func (a *addrBook) RemoveAddress(addr *p2p.NetAddress) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	ka := a.addrLookup[addr.ID]
	if ka == nil {
		return
	}
	a.Logger.Info("Remove address from book", "addr", ka.Addr, "ID", ka.ID)
	a.removeFromAllBuckets(ka)
}

func (a *addrBook) IsGood(addr *p2p.NetAddress) bool {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.addrLookup[addr.ID].isOld()
}

func (a *addrBook) HasAddress(addr *p2p.NetAddress) bool {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	ka := a.addrLookup[addr.ID]
	return ka != nil
}

func (a *addrBook) NeedMoreAddrs() bool {
	return a.Size() < needAddressThreshold
}

func (a *addrBook) PickAddress(biasTowardsNewAddrs int) *p2p.NetAddress {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.size() == 0 {
		return nil
	}
	if biasTowardsNewAddrs > 100 {
		biasTowardsNewAddrs = 100
	}
	if biasTowardsNewAddrs < 0 {
		biasTowardsNewAddrs = 0
	}

	oldCorrelation := math.Sqrt(float64(a.nOld)) * (100.0 - float64(biasTowardsNewAddrs))
	newCorrelation := math.Sqrt(float64(a.nNew)) * float64(biasTowardsNewAddrs)

	var bucket map[string]*knownAddress
	pickFromOldBucket := (newCorrelation+oldCorrelation)*a.rand.Float64() < oldCorrelation
	if (pickFromOldBucket && a.nOld == 0) ||
		(!pickFromOldBucket && a.nNew == 0) {
		return nil
	}

	for len(bucket) == 0 {
		if pickFromOldBucket {
			bucket = a.bucketsOld[a.rand.Intn(len(a.bucketsOld))]
		} else {
			bucket = a.bucketsNew[a.rand.Intn(len(a.bucketsNew))]
		}
	}

	randIndex := a.rand.Intn(len(bucket))
	for _, ka := range bucket {
		if randIndex == 0 {
			return ka.Addr
		}
		randIndex--
	}
	return nil
}

func (a *addrBook) MarkGood(addr *p2p.NetAddress) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	ka := a.addrLookup[addr.ID]
	if ka == nil {
		return
	}
	ka.markGood()
	if ka.isNew() {
		a.moveToOld(ka)
	}
}

func (a *addrBook) MarkAttempt(addr *p2p.NetAddress) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	ka := a.addrLookup[addr.ID]
	if ka == nil {
		return
	}
	ka.markAttempt()
}

func (a *addrBook) MarkBad(addr *p2p.NetAddress) {
	a.RemoveAddress(addr)
}

func (a *addrBook) GetSelection() []*p2p.NetAddress {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.size() == 0 {
		return nil
	}

	allAddr := make([]*p2p.NetAddress, a.size())
	i := 0
	for _, ka := range a.addrLookup {
		allAddr[i] = ka.Addr
		i++
	}

	numAddresses := cmn.MaxInt(
		cmn.MinInt(minGetSelection, len(allAddr)),
		len(allAddr)*getSelectionPercent/100)
	numAddresses = cmn.MinInt(maxGetSelection, numAddresses)

	for i := 0; i < numAddresses; i++ {

		j := cmn.RandIntn(len(allAddr)-i) + i
		allAddr[i], allAddr[j] = allAddr[j], allAddr[i]
	}

	return allAddr[:numAddresses]
}

func (a *addrBook) GetSelectionWithBias(biasTowardsNewAddrs int) []*p2p.NetAddress {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.size() == 0 {
		return nil
	}

	if biasTowardsNewAddrs > 100 {
		biasTowardsNewAddrs = 100
	}
	if biasTowardsNewAddrs < 0 {
		biasTowardsNewAddrs = 0
	}

	numAddresses := cmn.MaxInt(
		cmn.MinInt(minGetSelection, a.size()),
		a.size()*getSelectionPercent/100)
	numAddresses = cmn.MinInt(maxGetSelection, numAddresses)

	selection := make([]*p2p.NetAddress, numAddresses)

	oldBucketToAddrsMap := make(map[int]map[string]struct{})
	var oldIndex int
	newBucketToAddrsMap := make(map[int]map[string]struct{})
	var newIndex int

	selectionIndex := 0
ADDRS_LOOP:
	for selectionIndex < numAddresses {
		pickFromOldBucket := int((float64(selectionIndex)/float64(numAddresses))*100) >= biasTowardsNewAddrs
		pickFromOldBucket = (pickFromOldBucket && a.nOld > 0) || a.nNew == 0
		bucket := make(map[string]*knownAddress)

		for len(bucket) == 0 {
			if pickFromOldBucket {
				oldIndex = a.rand.Intn(len(a.bucketsOld))
				bucket = a.bucketsOld[oldIndex]
			} else {
				newIndex = a.rand.Intn(len(a.bucketsNew))
				bucket = a.bucketsNew[newIndex]
			}
		}

		randIndex := a.rand.Intn(len(bucket))

		var selectedAddr *p2p.NetAddress
		for _, ka := range bucket {
			if randIndex == 0 {
				selectedAddr = ka.Addr
				break
			}
			randIndex--
		}

		if pickFromOldBucket {
			if addrsMap, ok := oldBucketToAddrsMap[oldIndex]; ok {
				if _, ok = addrsMap[selectedAddr.String()]; ok {
					continue ADDRS_LOOP
				}
			} else {
				oldBucketToAddrsMap[oldIndex] = make(map[string]struct{})
			}
			oldBucketToAddrsMap[oldIndex][selectedAddr.String()] = struct{}{}
		} else {
			if addrsMap, ok := newBucketToAddrsMap[newIndex]; ok {
				if _, ok = addrsMap[selectedAddr.String()]; ok {
					continue ADDRS_LOOP
				}
			} else {
				newBucketToAddrsMap[newIndex] = make(map[string]struct{})
			}
			newBucketToAddrsMap[newIndex][selectedAddr.String()] = struct{}{}
		}

		selection[selectionIndex] = selectedAddr
		selectionIndex++
	}

	return selection
}

func (a *addrBook) ListOfKnownAddresses() []*knownAddress {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	addrs := []*knownAddress{}
	for _, addr := range a.addrLookup {
		addrs = append(addrs, addr.copy())
	}
	return addrs
}

func (a *addrBook) Size() int {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.size()
}

func (a *addrBook) size() int {
	return a.nNew + a.nOld
}

func (a *addrBook) Save() {
	a.saveToFile(a.filePath)
}

func (a *addrBook) saveRoutine() {
	defer a.wg.Done()

	saveFileTicker := time.NewTicker(dumpAddressInterval)
out:
	for {
		select {
		case <-saveFileTicker.C:
			a.saveToFile(a.filePath)
		case <-a.Quit():
			break out
		}
	}
	saveFileTicker.Stop()
	a.saveToFile(a.filePath)
	a.Logger.Info("Address handler done")
}

func (a *addrBook) getBucket(bucketType byte, bucketIdx int) map[string]*knownAddress {
	switch bucketType {
	case bucketTypeNew:
		return a.bucketsNew[bucketIdx]
	case bucketTypeOld:
		return a.bucketsOld[bucketIdx]
	default:
		cmn.PanicSanity("Should not happen")
		return nil
	}
}

func (a *addrBook) addToNewBucket(ka *knownAddress, bucketIdx int) bool {

	if ka.isOld() {
		a.Logger.Error(cmn.Fmt("Cannot add address already in old bucket to a new bucket: %v", ka))
		return false
	}

	addrStr := ka.Addr.String()
	bucket := a.getBucket(bucketTypeNew, bucketIdx)

	if _, ok := bucket[addrStr]; ok {
		return true
	}

	if len(bucket) > newBucketSize {
		a.Logger.Info("new bucket is full, expiring new")
		a.expireNew(bucketIdx)
	}

	bucket[addrStr] = ka

	if ka.addBucketRef(bucketIdx) == 1 {
		a.nNew++
	}

	a.addrLookup[ka.ID()] = ka

	return true
}

func (a *addrBook) addToOldBucket(ka *knownAddress, bucketIdx int) bool {

	if ka.isNew() {
		a.Logger.Error(cmn.Fmt("Cannot add new address to old bucket: %v", ka))
		return false
	}
	if len(ka.Buckets) != 0 {
		a.Logger.Error(cmn.Fmt("Cannot add already old address to another old bucket: %v", ka))
		return false
	}

	addrStr := ka.Addr.String()
	bucket := a.getBucket(bucketTypeOld, bucketIdx)

	if _, ok := bucket[addrStr]; ok {
		return true
	}

	if len(bucket) > oldBucketSize {
		return false
	}

	bucket[addrStr] = ka
	if ka.addBucketRef(bucketIdx) == 1 {
		a.nOld++
	}

	a.addrLookup[ka.ID()] = ka

	return true
}

func (a *addrBook) removeFromBucket(ka *knownAddress, bucketType byte, bucketIdx int) {
	if ka.BucketType != bucketType {
		a.Logger.Error(cmn.Fmt("Bucket type mismatch: %v", ka))
		return
	}
	bucket := a.getBucket(bucketType, bucketIdx)
	delete(bucket, ka.Addr.String())
	if ka.removeBucketRef(bucketIdx) == 0 {
		if bucketType == bucketTypeNew {
			a.nNew--
		} else {
			a.nOld--
		}
		delete(a.addrLookup, ka.ID())
	}
}

func (a *addrBook) removeFromAllBuckets(ka *knownAddress) {
	for _, bucketIdx := range ka.Buckets {
		bucket := a.getBucket(ka.BucketType, bucketIdx)
		delete(bucket, ka.Addr.String())
	}
	ka.Buckets = nil
	if ka.BucketType == bucketTypeNew {
		a.nNew--
	} else {
		a.nOld--
	}
	delete(a.addrLookup, ka.ID())
}

func (a *addrBook) pickOldest(bucketType byte, bucketIdx int) *knownAddress {
	bucket := a.getBucket(bucketType, bucketIdx)
	var oldest *knownAddress
	for _, ka := range bucket {
		if oldest == nil || ka.LastAttempt.Before(oldest.LastAttempt) {
			oldest = ka
		}
	}
	return oldest
}

func (a *addrBook) addAddress(addr, src *p2p.NetAddress) error {
	if a.routabilityStrict && !addr.Routable() {
		return fmt.Errorf("Cannot add non-routable address %v", addr)
	}
	if _, ok := a.ourAddrs[addr.String()]; ok {

		return fmt.Errorf("Cannot add ourselves with address %v", addr)
	}

	ka := a.addrLookup[addr.ID]

	if ka != nil {

		if ka.isOld() {
			return nil
		}

		if len(ka.Buckets) == maxNewBucketsPerAddress {
			return nil
		}

		factor := int32(2 * len(ka.Buckets))
		if a.rand.Int31n(factor) != 0 {
			return nil
		}
	} else {
		ka = newKnownAddress(addr, src)
	}

	bucket := a.calcNewBucket(addr, src)
	added := a.addToNewBucket(ka, bucket)
	if !added {
		a.Logger.Info("Can't add new address, addr book is full", "address", addr, "total", a.size())
	}

	a.Logger.Info("Added new address", "address", addr, "total", a.size())
	return nil
}

func (a *addrBook) expireNew(bucketIdx int) {
	for addrStr, ka := range a.bucketsNew[bucketIdx] {

		if ka.isBad() {
			a.Logger.Info(cmn.Fmt("expiring bad address %v", addrStr))
			a.removeFromBucket(ka, bucketTypeNew, bucketIdx)
			return
		}
	}

	oldest := a.pickOldest(bucketTypeNew, bucketIdx)
	a.removeFromBucket(oldest, bucketTypeNew, bucketIdx)
}

func (a *addrBook) moveToOld(ka *knownAddress) {

	if ka.isOld() {
		a.Logger.Error(cmn.Fmt("Cannot promote address that is already old %v", ka))
		return
	}
	if len(ka.Buckets) == 0 {
		a.Logger.Error(cmn.Fmt("Cannot promote address that isn't in any new buckets %v", ka))
		return
	}

	freedBucket := ka.Buckets[0]

	a.removeFromAllBuckets(ka)

	ka.BucketType = bucketTypeOld

	oldBucketIdx := a.calcOldBucket(ka.Addr)
	added := a.addToOldBucket(ka, oldBucketIdx)
	if !added {

		oldest := a.pickOldest(bucketTypeOld, oldBucketIdx)
		a.removeFromBucket(oldest, bucketTypeOld, oldBucketIdx)

		newBucketIdx := a.calcNewBucket(oldest.Addr, oldest.Src)
		added := a.addToNewBucket(oldest, newBucketIdx)

		if !added {
			added := a.addToNewBucket(oldest, freedBucket)
			if !added {
				a.Logger.Error(cmn.Fmt("Could not migrate oldest %v to freedBucket %v", oldest, freedBucket))
			}
		}

		added = a.addToOldBucket(ka, oldBucketIdx)
		if !added {
			a.Logger.Error(cmn.Fmt("Could not re-add ka %v to oldBucketIdx %v", ka, oldBucketIdx))
		}
	}
}

func (a *addrBook) calcNewBucket(addr, src *p2p.NetAddress) int {
	data1 := []byte{}
	data1 = append(data1, []byte(a.key)...)
	data1 = append(data1, []byte(a.groupKey(addr))...)
	data1 = append(data1, []byte(a.groupKey(src))...)
	hash1 := doubleSha256(data1)
	hash64 := binary.BigEndian.Uint64(hash1)
	hash64 %= newBucketsPerGroup
	var hashbuf [8]byte
	binary.BigEndian.PutUint64(hashbuf[:], hash64)
	data2 := []byte{}
	data2 = append(data2, []byte(a.key)...)
	data2 = append(data2, a.groupKey(src)...)
	data2 = append(data2, hashbuf[:]...)

	hash2 := doubleSha256(data2)
	return int(binary.BigEndian.Uint64(hash2) % newBucketCount)
}

func (a *addrBook) calcOldBucket(addr *p2p.NetAddress) int {
	data1 := []byte{}
	data1 = append(data1, []byte(a.key)...)
	data1 = append(data1, []byte(addr.String())...)
	hash1 := doubleSha256(data1)
	hash64 := binary.BigEndian.Uint64(hash1)
	hash64 %= oldBucketsPerGroup
	var hashbuf [8]byte
	binary.BigEndian.PutUint64(hashbuf[:], hash64)
	data2 := []byte{}
	data2 = append(data2, []byte(a.key)...)
	data2 = append(data2, a.groupKey(addr)...)
	data2 = append(data2, hashbuf[:]...)

	hash2 := doubleSha256(data2)
	return int(binary.BigEndian.Uint64(hash2) % oldBucketCount)
}

func (a *addrBook) groupKey(na *p2p.NetAddress) string {
	if a.routabilityStrict && na.Local() {
		return "local"
	}
	if a.routabilityStrict && !na.Routable() {
		return "unroutable"
	}

	if ipv4 := na.IP.To4(); ipv4 != nil {
		return (&net.IPNet{IP: na.IP, Mask: net.CIDRMask(16, 32)}).String()
	}
	if na.RFC6145() || na.RFC6052() {

		ip := net.IP(na.IP[12:16])
		return (&net.IPNet{IP: ip, Mask: net.CIDRMask(16, 32)}).String()
	}

	if na.RFC3964() {
		ip := net.IP(na.IP[2:7])
		return (&net.IPNet{IP: ip, Mask: net.CIDRMask(16, 32)}).String()

	}
	if na.RFC4380() {

		ip := net.IP(make([]byte, 4))
		for i, byte := range na.IP[12:16] {
			ip[i] = byte ^ 0xff
		}
		return (&net.IPNet{IP: ip, Mask: net.CIDRMask(16, 32)}).String()
	}

	bits := 32
	heNet := &net.IPNet{IP: net.ParseIP("2001:470::"),
		Mask:	net.CIDRMask(32, 128)}
	if heNet.Contains(na.IP) {
		bits = 36
	}

	return (&net.IPNet{IP: na.IP, Mask: net.CIDRMask(bits, 128)}).String()
}

func doubleSha256(b []byte) []byte {
	hasher := sha256.New()
	hasher.Write(b)
	sum := hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(sum)
	return hasher.Sum(nil)
}
