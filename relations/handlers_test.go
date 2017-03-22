package relations

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type test struct {
	name         string
	req          *http.Request
	dummyService Driver
	statusCode   int
	body         string
}

const knownUUID = "f78c1482-a65c-413e-b753-ca3ce3cb84f0"
const successfulResponse = `{"curatedRelatedContent":[{"id":"http://id-f78c1482-a65c-413e-b753-ca3ce3cb84f0", "apiUrl":"http://apiurl-f78c1482-a65c-413e-b753-ca3ce3cb84f0"}],
"contains":[{"id":"http://id-f78c1482-a65c-413e-b753-ca3ce3cb84f0", "apiUrl":"http://apiurl-f78c1482-a65c-413e-b753-ca3ce3cb84f0"}],
"containedIn":[{"id":"http://id-f78c1482-a65c-413e-b753-ca3ce3cb84f0", "apiUrl":"http://apiurl-f78c1482-a65c-413e-b753-ca3ce3cb84f0"}]}`

func TestGetHandler(t *testing.T) {
	tests := []test{
		{"Success", newRequest("GET", fmt.Sprintf("/content/%s/relations", knownUUID), nil), dummyService{contentUUID: knownUUID}, http.StatusOK, successfulResponse},
		{"NotFound", newRequest("GET", fmt.Sprintf("/content/%s/relations", "db90a9db-6cb6-4ba0-8648-c0676087aba2"), nil), dummyService{contentUUID: knownUUID}, http.StatusNotFound, message("No relations found for content with uuid db90a9db-6cb6-4ba0-8648-c0676087aba2")},
		{"InvalidUuid", newRequest("GET", fmt.Sprintf("/content/%s/relations", "99999"), nil), dummyService{contentUUID: knownUUID}, http.StatusBadRequest, message("The given uuid is not valid, err=uuid: UUID string too short: 99999")},
		{"ReadError", newRequest("GET", fmt.Sprintf("/content/%s/relations", knownUUID), nil), dummyService{contentUUID: knownUUID, failRead: true}, http.StatusServiceUnavailable, message("Error retrieving relations for f78c1482-a65c-413e-b753-ca3ce3cb84f0, err=TEST failing to READ")},
	}

	for _, test := range tests {
		hh := HttpHandlers{test.dummyService, ""}
		rec := httptest.NewRecorder()
		r := mux.NewRouter()
		r.HandleFunc("/content/{uuid}/relations", hh.GetRelations).Methods("GET")
		r.ServeHTTP(rec, test.req)
		assert.True(t, test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.JSONEq(t, test.body, rec.Body.String(), fmt.Sprintf("%s: Wrong body", test.name))
	}
}

func newRequest(method, url string, body []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	return req
}

func message(errMsg string) string {
	return fmt.Sprintf("{\"message\": \"%s\"}\n", errMsg)
}

type dummyService struct {
	contentUUID string
	failRead    bool
}

func (dS dummyService) read(contentUUID string) (relations, bool, error) {
	if dS.failRead {
		return relations{}, false, errors.New("TEST failing to READ")
	}
	if contentUUID == dS.contentUUID {
		return relations{
			[]relatedContent{{ID: "http://id-" + contentUUID, APIURL: "http://apiurl-" + contentUUID}},
			[]relatedContent{{ID: "http://id-" + contentUUID, APIURL: "http://apiurl-" + contentUUID}},
			[]relatedContent{{ID: "http://id-" + contentUUID, APIURL: "http://apiurl-" + contentUUID}},
		}, true, nil
	}
	return relations{}, false, nil
}

func (dS dummyService) checkConnectivity() error {
	return nil
}
