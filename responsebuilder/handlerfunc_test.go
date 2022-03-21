package responsebuilder_test

import (
	"github.com/ddelizia/gojitsu/responsebuilder"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHandlerFunc struct {
	mock.Mock
}

func (m *mockHandlerFunc) mockHandlerFunc(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	return
}

func Test_handlerFunc_Handle(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	// Given
	mockHandler := &mockHandlerFunc{}
	mockHandler.On("mockHandlerFunc", mock.Anything, mock.Anything)
	responseBuilder := responsebuilder.HandlerFunc(mockHandler.mockHandlerFunc)

	req, err := http.NewRequest(http.MethodGet, "/hello", nil)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	rr := httptest.NewRecorder()
	responseBuilder.Handle().ServeHTTP(rr, req)

	g.Expect(len(mockHandler.Calls)).To(gomega.Equal(1))
}
