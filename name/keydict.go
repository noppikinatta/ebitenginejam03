package name

var txtimgKeyMapping map[string]string = map[string]string{}

func ImgKey(txtKey string) string {
	k, ok := txtimgKeyMapping[txtKey]
	if ok {
		return k
	}

	// fallback image may be drawn
	return "fallback"
}
