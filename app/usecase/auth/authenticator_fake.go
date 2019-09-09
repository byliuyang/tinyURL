package auth

import (
	"time"

	"github.com/byliuyang/app/mdtest"
)

func NewAuthenticatorFake(current time.Time, validPeriod time.Duration) Authenticator {
	tokenizer := mdtest.CryptoTokenizerFake{}
	timer := mdtest.NewTimerFake(current)
	return NewAuthenticator(tokenizer, timer, validPeriod)
}
