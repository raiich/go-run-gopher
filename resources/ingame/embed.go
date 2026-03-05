package ingame

import _ "embed"

//go:embed go-run1.png
var CharaImageBytes1 []byte

var (
	//go:embed go-run2.png
	CharaImageBytes2 []byte

	//go:embed go-run3.png
	CharaImageBytes3 []byte
)
