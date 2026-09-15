package long

var handler = func() {
	var userConfigParam = 1 // want `variable name has 3 words \(user, Config, Param\): use two at most`
	_ = userConfigParam
}

func f() (int, error) {
	return 0, nil
}

func g(s []int, m map[string]int, ch chan int, x any) {
	var parsedJSONDocument = 1 // want `variable name has 3 words \(parsed, JSON, Document\): use two at most`
	activeUserAddr := true     // want `variable name has 3 words \(active, User, Addr\)`
	var a, maxRetryCount = 1, 2 // want `variable name has 3 words \(max, Retry, Count\)`
	var (
		lastSeenValue int // want `\(last, Seen, Value\)`
	)
	readByteCount, err := f() // want `\(read, Byte, Count\)`
	readByteCount, lastErr := f()
	var user_config_param string // want `\(user, config, param\)`
	var base64URLEncoding string // want `\(base64, URL, Encoding\)`
	var userIDsByName map[string]int // want `variable name has 4 words \(user, IDs, By, Name\)`
	var HTTPServerAddr string        // want `\(HTTP, Server, Addr\)`
	for itemIndexValue := range s { // want `\(item, Index, Value\)`
		_ = itemIndexValue
	}
	for key, mapItemValue := range m { // want `\(map, Item, Value\)`
		_, _ = key, mapItemValue
	}
	for loopCounterValue := 0; loopCounterValue < a; loopCounterValue++ { // want `\(loop, Counter, Value\)`
	}
	if lookedUpValue, ok := m["k"]; ok { // want `\(looked, Up, Value\)`
		_ = lookedUpValue
	}
	switch typedInputValue := x.(type) { // want `\(typed, Input, Value\)`
	case int:
		_ = typedInputValue
	case string:
		_ = typedInputValue
	}
	select {
	case receivedChanValue := <-ch: // want `\(received, Chan, Value\)`
		_ = receivedChanValue
	}
	var inner = func(requestBodyText string) {
		var requestBodyLength = len(requestBodyText) // want `\(request, Body, Length\)`
		_ = requestBodyLength
	}
	inner("")
	_, _, _, _, _, _, _, _, _, _, _ = parsedJSONDocument, maxRetryCount, activeUserAddr, lastSeenValue, err, lastErr, readByteCount, user_config_param, base64URLEncoding, userIDsByName, HTTPServerAddr
}
