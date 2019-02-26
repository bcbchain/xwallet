package db

type DB interface {
	Get([]byte) []byte

	Has(key []byte) bool

	Set([]byte, []byte)
	SetSync([]byte, []byte)

	Delete([]byte)
	DeleteSync([]byte)

	Iterator(start, end []byte) Iterator

	ReverseIterator(start, end []byte) Iterator

	Close()

	NewBatch() Batch

	Print()

	Stats() map[string]string
}

type Batch interface {
	SetDeleter
	Write()
	WriteSync()
}

type SetDeleter interface {
	Set(key, value []byte)
	Delete(key []byte)
}

type Iterator interface {
	Domain() (start []byte, end []byte)

	Valid() bool

	Next()

	Key() (key []byte)

	Value() (value []byte)

	Close()
}

func bz(s string) []byte {
	return []byte(s)
}

func nonNilBytes(bz []byte) []byte {
	if bz == nil {
		return []byte{}
	}
	return bz
}
