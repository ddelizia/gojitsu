package requestmatcher_test

import (
	"github.com/ddelizia/gojitsu"
	"github.com/ddelizia/gojitsu/requestmatcher"
	"github.com/gorilla/mux"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockEchoHandler struct {
	mock.Mock
}

func (m *mockEchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	return
}

func Test_requestMatcherSimple_Setup(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	matcher := requestmatcher.Simple(
		[]string{http.MethodGet},
		gojitsu.String("/thisIsThePath"),
		[]string{"header", "value"},
	)
	router := mux.NewRouter()

	mockEchoHandler := &mockEchoHandler{}
	mockEchoHandler.On("ServeHTTP", mock.Anything, mock.Anything)
	matcher.Setup(router).HandlerFunc(mockEchoHandler.ServeHTTP)

	req, err := http.NewRequest(http.MethodGet, "/thisIsThePath", nil)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	req.Header.Add("header", "value")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	g.Expect(len(mockEchoHandler.Calls)).To(gomega.Equal(1))
}
