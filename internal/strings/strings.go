package strings

import (
	"strings"

	internalErrors "github.com/theTardigrade/runner/internal/errors"
)

type Builder strings.Builder

func (b *Builder) WriteString(s string) {
	_, err := (*strings.Builder)(b).WriteString(s)
	internalErrors.Check(err)
}

func (b *Builder) WriteByte(t byte) {
	err := (*strings.Builder)(b).WriteByte(t)
	internalErrors.Check(err)
}

func (b *Builder) String() string {
	return (*strings.Builder)(b).String()
}
