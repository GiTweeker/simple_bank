// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package token

import (
	"sync"
	"time"
)

// Ensure, that MakerMock does implement Maker.
// If this is not the case, regenerate this file with moq.
var _ Maker = &MakerMock{}

// MakerMock is a mock implementation of Maker.
//
//	func TestSomethingThatUsesMaker(t *testing.T) {
//
//		// make and configure a mocked Maker
//		mockedMaker := &MakerMock{
//			CreateTokenFunc: func(username string, duration time.Duration) (string, *Payload, error) {
//				panic("mock out the CreateToken method")
//			},
//			VerifyTokenFunc: func(token string) (*Payload, error) {
//				panic("mock out the VerifyToken method")
//			},
//		}
//
//		// use mockedMaker in code that requires Maker
//		// and then make assertions.
//
//	}
type MakerMock struct {
	// CreateTokenFunc mocks the CreateToken method.
	CreateTokenFunc func(username string, duration time.Duration) (string, *Payload, error)

	// VerifyTokenFunc mocks the VerifyToken method.
	VerifyTokenFunc func(token string) (*Payload, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateToken holds details about calls to the CreateToken method.
		CreateToken []struct {
			// Username is the username argument value.
			Username string
			// Duration is the duration argument value.
			Duration time.Duration
		}
		// VerifyToken holds details about calls to the VerifyToken method.
		VerifyToken []struct {
			// Token is the token argument value.
			Token string
		}
	}
	lockCreateToken sync.RWMutex
	lockVerifyToken sync.RWMutex
}

// CreateToken calls CreateTokenFunc.
func (mock *MakerMock) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	if mock.CreateTokenFunc == nil {
		panic("MakerMock.CreateTokenFunc: method is nil but Maker.CreateToken was just called")
	}
	callInfo := struct {
		Username string
		Duration time.Duration
	}{
		Username: username,
		Duration: duration,
	}
	mock.lockCreateToken.Lock()
	mock.calls.CreateToken = append(mock.calls.CreateToken, callInfo)
	mock.lockCreateToken.Unlock()
	return mock.CreateTokenFunc(username, duration)
}

// CreateTokenCalls gets all the calls that were made to CreateToken.
// Check the length with:
//
//	len(mockedMaker.CreateTokenCalls())
func (mock *MakerMock) CreateTokenCalls() []struct {
	Username string
	Duration time.Duration
} {
	var calls []struct {
		Username string
		Duration time.Duration
	}
	mock.lockCreateToken.RLock()
	calls = mock.calls.CreateToken
	mock.lockCreateToken.RUnlock()
	return calls
}

// VerifyToken calls VerifyTokenFunc.
func (mock *MakerMock) VerifyToken(token string) (*Payload, error) {
	if mock.VerifyTokenFunc == nil {
		panic("MakerMock.VerifyTokenFunc: method is nil but Maker.VerifyToken was just called")
	}
	callInfo := struct {
		Token string
	}{
		Token: token,
	}
	mock.lockVerifyToken.Lock()
	mock.calls.VerifyToken = append(mock.calls.VerifyToken, callInfo)
	mock.lockVerifyToken.Unlock()
	return mock.VerifyTokenFunc(token)
}

// VerifyTokenCalls gets all the calls that were made to VerifyToken.
// Check the length with:
//
//	len(mockedMaker.VerifyTokenCalls())
func (mock *MakerMock) VerifyTokenCalls() []struct {
	Token string
} {
	var calls []struct {
		Token string
	}
	mock.lockVerifyToken.RLock()
	calls = mock.calls.VerifyToken
	mock.lockVerifyToken.RUnlock()
	return calls
}