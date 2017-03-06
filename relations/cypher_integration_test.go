package relations

import (
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/content-collection-rw-neo4j/collection"
	"github.com/Financial-Times/content-rw-neo4j/content"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type payloadData struct {
	uuid   string
	path   string
	id     string
	apiURL string
}

var (
	leadContent = payloadData{"3fc9fe3e-af8c-4a4a-961a-e5065392bb31", "./fixtures/Content-with-SP-3fc9fe3e-af8c-4a4a-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-4a4a-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-4a4a-961a-e5065392bb31"}
	relatedContent1 = payloadData{"3fc9fe3e-af8c-1a1a-961a-e5065392bb31", "./fixtures/Content-3fc9fe3e-af8c-1a1a-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-1a1a-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-1a1a-961a-e5065392bb31"}
	relatedContent2 = payloadData{"3fc9fe3e-af8c-2a2a-961a-e5065392bb31", "./fixtures/Content-3fc9fe3e-af8c-2a2a-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-2a2a-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-2a2a-961a-e5065392bb31"}
	relatedContent3 = payloadData{"3fc9fe3e-af8c-3a3a-961a-e5065392bb31", "./fixtures/Content-3fc9fe3e-af8c-3a3a-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-3a3a-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-3a3a-961a-e5065392bb31"}
	storyPackage = payloadData{"63559ba7-b48d-4467-b2b0-ce956f9e9494", "./fixtures/StoryPackage-63559ba7-b48d-4467-b2b0-ce956f9e9494.json",
		"", ""}
	allData = []payloadData{leadContent, relatedContent1, relatedContent2, relatedContent3, storyPackage}
)

func TestRetrieveCuratedRelatedContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Short flag is set. Skipping integration test")
	}
	expectedResponse := relations{
		CuratedRelatedContents: []relatedContent{
			{relatedContent1.id, relatedContent1.apiURL},
			{relatedContent2.id, relatedContent2.apiURL},
			{relatedContent3.id, relatedContent3.apiURL},
		},
	}
	db := getDatabaseConnection(t)
	contents := []payloadData{leadContent, relatedContent1, relatedContent2, relatedContent3}
	cleanDB(t, db, allData)

	writeContent(t, db, contents)
	writeStoryPackage(t, db, []payloadData{storyPackage})
	defer cleanDB(t, db, allData)

	driver := NewCypherDriver(db)
	actualCRC, found, err := driver.read(leadContent.uuid)
	assert.NoError(t, err, "Unexpected error for content %s", leadContent.uuid)
	assert.True(t, found, "Found no relations for content %s", leadContent.uuid)
	assert.Equal(t, len(expectedResponse.CuratedRelatedContents), len(actualCRC.CuratedRelatedContents), "Didn't get the same number of curated related content")
	assertListContainsAll(t, actualCRC.CuratedRelatedContents, expectedResponse.CuratedRelatedContents)
}

func writeContent(t testing.TB, db neoutils.NeoConnection, data []payloadData) baseftrwapp.Service {
	contentRW := content.NewCypherContentService(db)
	assert.NoError(t, contentRW.Initialise())
	for _, d := range data {
		writeJSONWithService(t, contentRW, d.path)
	}
	return contentRW
}

func writeStoryPackage(t testing.TB, db neoutils.NeoConnection, data []payloadData) collection.Service {
	contentCollectionRW := collection.NewContentCollectionService(db)
	assert.NoError(t, contentCollectionRW.Initialise())
	for _, d := range data {
		writeJSONWithContentCollectionService(t, contentCollectionRW, d.path)
	}
	return contentCollectionRW
}

func writeJSONWithService(t testing.TB, service baseftrwapp.Service, pathToJSONFile string) {
	path, err := filepath.Abs(pathToJSONFile)
	require.NoError(t, err)
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()
	require.NoError(t, err)
	dec := json.NewDecoder(f)
	inst, _, err := service.DecodeJSON(dec)
	require.NoError(t, err)
	err = service.Write(inst)
	require.NoError(t, err)
}

func writeJSONWithContentCollectionService(t testing.TB, service collection.Service, pathToJSONFile string) {
	path, err := filepath.Abs(pathToJSONFile)
	require.NoError(t, err)
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()
	require.NoError(t, err)
	dec := json.NewDecoder(f)
	inst, _, err := service.DecodeJSON(dec)
	require.NoError(t, err)
	err = service.Write(inst, "StoryPackage")
	require.NoError(t, err)
}

func assertListContainsAll(t *testing.T, list interface{}, items ...interface{}) {
	if reflect.TypeOf(items[0]).Kind().String() == "slice" {
		expected := reflect.ValueOf(items[0])
		expectedLength := expected.Len()
		assert.Len(t, list, expectedLength)
		for i := 0; i < expectedLength; i++ {
			assert.Contains(t, list, expected.Index(i).Interface())
		}
	} else {
		assert.Len(t, list, len(items))
		for _, item := range items {
			assert.Contains(t, list, item)
		}
	}
}

func getDatabaseConnection(t testing.TB) neoutils.NeoConnection {
	db, err := getDBConn()
	require.NoError(t, err, "Failed to connect to Neo4j")
	return db
}

func getDBConn() (neoutils.NeoConnection, error) {
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false
	return neoutils.Connect(url, conf)
}

func cleanDB(t testing.TB, db neoutils.NeoConnection, data []payloadData) {
	qs := make([]*neoism.CypherQuery, len(data))
	for i, d := range data {
		qs[i] = &neoism.CypherQuery{
			Statement: fmt.Sprintf(`
			MATCH (a:Thing {uuid: "%s"})
			OPTIONAL MATCH (a)<-[iden:IDENTIFIES]-(i:Identifier)
			DELETE iden, i
			DETACH DELETE a`, d.uuid)}
	}
	err := db.CypherBatch(qs)
	assert.NoError(t, err)
}
