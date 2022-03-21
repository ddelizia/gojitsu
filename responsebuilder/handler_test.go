package responsebuilder_test

import (
	"github.com/ddelizia/gojitsu/responsebuilder"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	return
}

func Test_handler_Handle(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	// Given
	mockHandler := &mockHandler{}
	responseBuilder := responsebuilder.Handler(mockHandler)
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything)

	req, err := http.NewRequest(http.MethodGet, "/hello", nil)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	rr := httptest.NewRecorder()
	responseBuilder.Handle().ServeHTTP(rr, req)

	g.Expect(len(mockHandler.Calls)).To(gomega.Equal(1))
}
