package relations

import (
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
)

var givenNeoRelatedContent []neoRelatedContent
var expectedRelatedContent []relatedContent

func TestTransformToRelatedContent(t *testing.T) {
	cd := NewCypherDriver(nil)

	givenNeoRelatedContent = []neoRelatedContent{
		{UUID: "db90a9db-6cb6-4ba0-8648-c0676087aba2"},
		{UUID: "f78c1482-abab-413e-b753-ca3ce3cb84f0"},
	}

	expectedRelatedContent = []relatedContent{
		{ID: "http://www.ft.com/thing/db90a9db-6cb6-4ba0-8648-c0676087aba2", APIURL: "http://api.ft.com/content/db90a9db-6cb6-4ba0-8648-c0676087aba2"},
		{ID: "http://www.ft.com/thing/f78c1482-abab-413e-b753-ca3ce3cb84f0", APIURL: "http://api.ft.com/content/f78c1482-abab-413e-b753-ca3ce3cb84f0"},
	}

	relatedContent := cd.transformToRelatedContent(givenNeoRelatedContent)

	expected, _ := json.Marshal(expectedRelatedContent)
	actual, _ := json.Marshal(relatedContent)
	assert.JSONEq(t, string(expected), string(actual))
}
