// Code generated by counterfeiter. DO NOT EDIT.
package gojitsufakes

import (
	"net/http"
	"sync"

	"github.com/ddelizia/gojitsu"
)

type FakeResponseBuilder struct {
	HandleStub        func() http.HandlerFunc
	handleMutex       sync.RWMutex
	handleArgsForCall []struct {
	}
	handleReturns struct {
		result1 http.HandlerFunc
	}
	handleReturnsOnCall map[int]struct {
		result1 http.HandlerFunc
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeResponseBuilder) Handle() http.HandlerFunc {
	fake.handleMutex.Lock()
	ret, specificReturn := fake.handleReturnsOnCall[len(fake.handleArgsForCall)]
	fake.handleArgsForCall = append(fake.handleArgsForCall, struct {
	}{})
	stub := fake.HandleStub
	fakeReturns := fake.handleReturns
	fake.recordInvocation("Handle", []interface{}{})
	fake.handleMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeResponseBuilder) HandleCallCount() int {
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	return len(fake.handleArgsForCall)
}

func (fake *FakeResponseBuilder) HandleCalls(stub func() http.HandlerFunc) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = stub
}

func (fake *FakeResponseBuilder) HandleReturns(result1 http.HandlerFunc) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = nil
	fake.handleReturns = struct {
		result1 http.HandlerFunc
	}{result1}
}

func (fake *FakeResponseBuilder) HandleReturnsOnCall(i int, result1 http.HandlerFunc) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = nil
	if fake.handleReturnsOnCall == nil {
		fake.handleReturnsOnCall = make(map[int]struct {
			result1 http.HandlerFunc
		})
	}
	fake.handleReturnsOnCall[i] = struct {
		result1 http.HandlerFunc
	}{result1}
}

func (fake *FakeResponseBuilder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeResponseBuilder) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ gojitsu.ResponseBuilder = new(FakeResponseBuilder)