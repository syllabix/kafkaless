package reverser

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type Service interface {
	Reverse(context.Context, string) (string, error)
}

type reverser struct {
	weaver.Implements[Service]
}

// Reverse the provided string and return it. This function will never fail
// despite returning an error
func (r *reverser) Reverse(ctx context.Context, str string) (string, error) {
	runes := []rune(str)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}

	return string(runes), nil
}
