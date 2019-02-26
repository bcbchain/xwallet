package client

type ABCIQueryOptions struct {
	Height	int64
	Trusted	bool
}

var DefaultABCIQueryOptions = ABCIQueryOptions{Height: 0, Trusted: false}
