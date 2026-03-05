package credits

import _ "embed"

var (
	//go:embed fonts/license.md
	fonts []byte

	//go:embed go/the-go-gopher.md
	goGopher []byte

	//go:embed go/the-go-programming-language.md
	goProgrammingLanguage []byte

	//go:embed generated/licenses.md
	licenses []byte
)

// Fonts returns the embedded fonts license text
func Fonts() []byte {
	return fonts
}

// GoGopher returns the embedded Go Gopher license text
func GoGopher() []byte {
	return goGopher
}

// GoProgrammingLanguage returns the embedded Go Programming Language license text
func GoProgrammingLanguage() []byte {
	return goProgrammingLanguage
}

// Licenses returns the embedded third-party licenses text
func Licenses() []byte {
	return licenses
}
