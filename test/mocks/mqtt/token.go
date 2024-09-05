package mqtt

import "time"

//go:generate mockgen -source=token.go -destination=../gomock/clients/mqtt/token/token.go

// Token copy interface for mocks.
type Token interface {
	// Wait will wait indefinitely for the Token to complete, ie the Publish
	// to be sent and confirmed receipt from the broker.
	Wait() bool

	// WaitTimeout takes a time.Duration to wait for the flow associated with the
	// Token to complete, returns true if it returned before the timeout or
	// returns false if the timeout occurred. In the case of a timeout the Token
	// does not have an error set in case the caller wishes to wait again.
	WaitTimeout(time.Duration) bool

	// Done returns a channel that is closed when the flow associated
	// with the Token completes. Clients should call Error after the
	// channel is closed to check if the flow completed successfully.
	//
	// Done is provided for use in select statements. Simple use cases may
	// use Wait or WaitTimeout.
	Done() <-chan struct{}

	Error() error
}
