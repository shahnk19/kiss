package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"kiss/web/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	apiPrefix = "/api/encode?url=https://www.test.com/"
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

//dummy models
type mockModel struct {
	lastId    int
	currentId int
}

func (m *mockModel) SaveTiny(surl, id string) (string, error) {
	if m.currentId > 0 && m.currentId > m.lastId {
		m.lastId++
		return "", errors.New("Id collision mock")
	}
	return id, nil
}

func (m *mockModel) GetLastId() int {
	m.lastId++
	return m.lastId
}

func (m *mockModel) GetDestination(id int) (string, error) {
	return "", nil
}

func TestEndpointsEncode_Basic(t *testing.T) {
	want := VM{Value: "1", Status: true, Error: ""}
	got := runHttpTestRequests(t, apiPrefix, &mockModel{lastId: 0})
	if want != got {
		t.Error(fmt.Sprintf("(!!!)Fail (want:%v,got:%v)", want, got))
	}
}

func TestEndpointsEncode_LastIdAlreadyExists(t *testing.T) { //will force retry 1 times
	want := VM{Value: "2", Status: true, Error: ""}
	got := runHttpTestRequests(t, apiPrefix, &mockModel{lastId: 1, currentId: 2})
	if want != got {
		t.Error(fmt.Sprintf("(!!!)Fail (want:%v,got:%v)", want, got))
	}
}

func TestEndpointsEncode_LastIdAlreadyExistsLastRetry(t *testing.T) { //force to go to the last retry
	want := VM{Value: "2", Status: true, Error: ""}
	got := runHttpTestRequests(t, apiPrefix, &mockModel{lastId: 1, currentId: 6})
	if want != got {
		t.Error(fmt.Sprintf("(!!!)Fail (want:%v,got:%v)", want, got))
	}
}

func TestEndpointsEncode_LastIdAlreadyExistsExhaustRetry(t *testing.T) { //exhaust all retry
	want := VM{Value: "", Status: false, Error: ERR_NEW_ENTRY_FAIL}
	got := runHttpTestRequests(t, apiPrefix, &mockModel{lastId: 1, currentId: 7})
	if want != got {
		t.Error(fmt.Sprintf("(!!!)Fail (want:%v,got:%v)", want, got))
	}
}

func runHttpTestRequests(t *testing.T, endpoint string, model models.IModel) VM {
	ctrl := getTestController(model)
	router := getTestRouter(ctrl)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", endpoint, nil)
	router.ServeHTTP(rec, req)
	var got VM
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	if err != nil {
		t.Error(fmt.Sprintf("(!!!)Fail unmarshall response :%v", err))
	}
	return got
}

func getTestController(model models.IModel) *Ctrl {
	enc := getBaseEncoder()

	return &Ctrl{
		Name:    "mock",
		encoder: enc,
		model:   model, //&mockModel{lastId: 0},
	}
}

func getTestRouter(ctrl *Ctrl) *gin.Engine {
	router := gin.New()
	router.GET("/api/encode", Encode(ctrl))
	return router
}
