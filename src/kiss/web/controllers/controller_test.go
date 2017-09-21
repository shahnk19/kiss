package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

//test any url validation that is not covered by golang lib url.ParseRequestURI
func TestIsValidUrl(t *testing.T) {
	testCases := []struct {
		what string
		want bool
	}{
		{"", false},
		{"https://www.google.com", true},
	}
	for _, ttc := range testCases {
		got := isValidUrl(ttc.what)
		if ttc.want != got {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%s)", ttc.what, ttc.want, got))
		}
	}
}

//test the endponts
type mockModel struct{}

func (m *mockModel) NewEntry(surl string) (int, error) {
	sa := strings.Split(surl, "https://www.test.com/")
	salen := len(sa)
	if salen > 0 {
		i, _ := strconv.Atoi(sa[salen-1])
		return i, nil
	}
	return 0, nil
}

func TestEndpointsEncode(t *testing.T) {
	apiPrefix := "/api/encode?url=https://www.test.com/"
	testCases := []struct {
		what string
		want VM
	}{
		{"1", VM{Value: "1", Status: true, Error: ""}},
		{"1001", VM{Value: "3e9", Status: true, Error: ""}},
	}
	var got VM
	for _, ttc := range testCases {
		rec := runHttpTestRequests(apiPrefix + ttc.what)
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		if err != nil {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%v)", ttc.what, ttc.want, got))
		}
		if ttc.want != got {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%v)", ttc.what, ttc.want, got))
		}
	}
}

func runHttpTestRequests(endpoint string) *httptest.ResponseRecorder {
	ctrl := getTestController()
	router := getTestRouter(ctrl)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", endpoint, nil)
	router.ServeHTTP(rec, req)
	return rec
}

func getTestController() *Ctrl {
	enc := getBaseEncoder()

	return &Ctrl{
		Name:    "mock",
		encoder: enc,
		model:   &mockModel{},
	}
}

func getTestRouter(ctrl *Ctrl) *gin.Engine {
	router := gin.New()
	router.GET("/api/encode", Encode(ctrl))
	return router
}
