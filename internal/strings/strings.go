package strings

import (
	"strings"

	internalErrors "github.com/theTardigrade/runner/internal/errors"
)

type Builder strings.Builder

func (b *Builder) Native() *strings.Builder {
	return (*strings.Builder)(b)
}

func (b *Builder) WriteString(s string) {
	_, err := b.Native().WriteString(s)
	internalErrors.Check(err)
}

func (b *Builder) WriteByte(t byte) {
	err := b.Native().WriteByte(t)
	internalErrors.Check(err)
}

func (b *Builder) String() string {
	return b.Native().String()
}
