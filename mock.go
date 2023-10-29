package mox

import (
	"fmt"
)

// Mock is ret mock object for testing.
// It is used to mock functions.
type Mock interface {
	// Call calls the function with the given name and arguments.
	// It returns the values that were set with SetReturn.
	// If no values were set, it panics.
	CallMethod(name string, args ...any) []any
	// ClearReturns clears the return values for the given function.
	// If no values were set, it does not perform an operation.
	ClearReturns(name string)
	// SetReturns sets the return values for the given function.
	// If values were already set, they are overwritten.
	SetReturns(name string, values ...any)
	// GetCallParams returns the arguments that were passed to the function.
	// If no values were set, it panics.
	CallParam(name string, call, position int) any
	// GetCallCount returns the number of times the function was called.
	// If no values were set, it panics.
	CallCount(name string) int
}

var NewMock = newMock

type call struct {
	args [][]any
	rets []any
}

type mock struct {
	calls map[string]call
}

func newMock() Mock {
	return &mock{
		calls: make(map[string]call),
	}
}

func (m *mock) CallMethod(name string, args ...any) []any {
	call, ok := m.calls[name]
	if !ok {
		panic(fmt.Sprintf("mock: no such function %s", name))
	}
	call.args = append(call.args, args)
	m.calls[name] = call
	return call.rets
}

func (m *mock) ClearReturns(name string) {
	mo, ok := m.calls[name]
	if !ok {
		return
	}
	mo.rets = make([]any, 0)
	m.calls[name] = mo
}

func (m *mock) SetReturns(name string, values ...any) {
	mo, ok := m.calls[name]
	if !ok {
		mo = call{
			args: make([][]any, 0),
			rets: make([]any, 0),
		}
	}
	mo.rets = values
	m.calls[name] = mo
}

func (m *mock) CallParam(name string, call, position int) any {
	ret, ok := m.calls[name]
	if !ok {
		panic(fmt.Sprintf("mock: no such function %s", name))
	}
	return ret.args[call][position]
}

func (m *mock) CallCount(name string) int {
	ret, ok := m.calls[name]
	if !ok {
		panic(fmt.Sprintf("mock: no such function %s", name))
	}
	return len(ret.args)
}

func Error(ret any) error {
	if ret == nil {
		return nil
	}

	return ret.(error)
}
