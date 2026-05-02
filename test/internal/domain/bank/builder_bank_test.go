package bank_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/bank"
	"gitflic.ru/lms/internal/domain/bank/title"
	"github.com/stretchr/testify/assert"
)

type bankBuilder struct {
	title string
}

func newBankBuilder() *bankBuilder {
	return &bankBuilder{
		title: "",
	}
}

func (b *bankBuilder) withTitle(s string) *bankBuilder {
	b.title = s
	return b
}

func (builder *bankBuilder) build(t *testing.T, wantErr error) *bank.Bank {
	t.Helper()

	titleVO, err := title.New(builder.title)
	if err != nil {
		assert.ErrorIs(t, err, wantErr)
		return nil
	}

	b, err := bank.New(titleVO)
	assert.ErrorIs(t, err, wantErr)

	return b
}
