package main

import (
	"fmt"
	"bcbchain.io/tx"
)

func main() {

	UnwrapperCLT()
}

func UnwrapperCLT() {

	tx.InitUnWrapper("local")

	s := "local<tx>.v1.soAgomCKRcxHgXgrHRXJvm9VKfbR7gLz3XHmvQymgatZcZg4V1SUzuQJU3wYyyosaxU7CYz7PUT6qPng1dLnGKqes2EyVWBJv4eXn5neLri4zueAJdJdHnsX5dCFwgNruAz4U6XbtiQz6uuMURWaTs5QFf5jG9zzA6UfYLLpP67uemkPGec6Wj7tkzP4jKRsS9SKLD1bafqTCPWy3e6W6AtzKHbunUBTd9uWA6PMa3vBj4dYjFpoEdvf7SDMpTQ7P9MhB9nnMG5SEmY3893V1WMV2HwLUx6LHPqjRCeCQDiPyaMNBadbdKULjS7nSZNvrLg9HRrJYybRfnwu3ReHd1b9vrHzozA.<1>.YTgiA1gdDGi2L8hw4m73eYYfpkLgrWEtRzU5fUbW2ztgSccoydemsQMNPAVsYXKTxSJpagUSSE4qqPPcoyUfH1GqFyN23eQbFzXQpyVCB5zgXUD14y6uUYJ29wxr9eZgmnDHKBdGBG5rx4Gs8BmQd"

	res := tx.UnpackAndParseTx(s)
	fmt.Println(res)
}
