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
	leadContentSP = payloadData{"3fc9fe3e-af8c-4a4a-961a-e5065392bb31", "./fixtures/Content-with-SP-3fc9fe3e-af8c-4a4a-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-4a4a-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-4a4a-961a-e5065392bb31"}
	leadContentCP = payloadData{"3fc9fe3e-af8c-1b1b-961a-e5065392bb31", "./fixtures/Content-with-CP-3fc9fe3e-af8c-1b1b-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-1b1b-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-1b1b-961a-e5065392bb31"}
	relatedContent1 = payloadData{"3fc9fe3e-af8c-1a1a-961a-e5065392bb31", "./fixtures/Content-3fc9fe3e-af8c-1a1a-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-1a1a-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-1a1a-961a-e5065392bb31"}
	relatedContent2 = payloadData{"3fc9fe3e-af8c-2a2a-961a-e5065392bb31", "./fixtures/Content-3fc9fe3e-af8c-2a2a-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-2a2a-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-2a2a-961a-e5065392bb31"}
	relatedContent3 = payloadData{"3fc9fe3e-af8c-3a3a-961a-e5065392bb31", "./fixtures/Content-3fc9fe3e-af8c-3a3a-961a-e5065392bb31.json",
		"http://api.ft.com/things/3fc9fe3e-af8c-3a3a-961a-e5065392bb31", "http://api.ft.com/content/3fc9fe3e-af8c-3a3a-961a-e5065392bb31"}
	storyPackage = payloadData{"63559ba7-b48d-4467-b2b0-ce956f9e9494", "./fixtures/StoryPackage-63559ba7-b48d-4467-b2b0-ce956f9e9494.json",
		"", ""}
	contentPackage = payloadData{"63559ba7-b48d-4467-b2b0-ce956f9e9494", "./fixtures/ContentPackage-63559ba7-b48d-4467-1b1b-ce956f9e9494.json",
		"", ""}
	allData = []payloadData{leadContentSP, leadContentCP, relatedContent1, relatedContent2, relatedContent3, storyPackage, contentPackage}
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
	contents := []payloadData{leadContentSP, relatedContent1, relatedContent2, relatedContent3}
	cleanDB(t, db, allData)

	writeContent(t, db, contents)
	writeContentCollection(t, db, []payloadData{storyPackage}, "StoryPackage")
	defer cleanDB(t, db, allData)

	driver := NewCypherDriver(db)
	actualCRC, found, err := driver.read(leadContentSP.uuid)
	assert.NoError(t, err, "Unexpected error for content %s", leadContentSP.uuid)
	assert.True(t, found, "Found no relations for content %s", leadContentSP.uuid)
	assert.Equal(t, len(expectedResponse.CuratedRelatedContents), len(actualCRC.CuratedRelatedContents), "Didn't get the same number of curated related content")
	assertListContainsAll(t, actualCRC.CuratedRelatedContents, expectedResponse.CuratedRelatedContents)
}

//func TestRetrieveContainsContent(t *testing.T) {
//	if testing.Short() {
//		t.Skip("Short flag is set. Skipping integration test")
//	}
//	expectedResponse := relations{
//		Contains: []relatedContent{
//			{relatedContent1.id, relatedContent1.apiURL},
//			{relatedContent2.id, relatedContent2.apiURL},
//		},
//	}
//	db := getDatabaseConnection(t)
//	contents := []payloadData{leadContentCP, relatedContent1, relatedContent2}
//	cleanDB(t, db, allData)
//
//	writeContent(t, db, contents)
//	writeContentCollection(t, db, []payloadData{contentPackage}, "ContentPackage")
//	defer cleanDB(t, db, allData)
//
//	driver := NewCypherDriver(db)
//	actualResponse, found, err := driver.read(leadContentCP.uuid)
//	assert.NoError(t, err, "Unexpected error for content %s", leadContentCP.uuid)
//	assert.True(t, found, "Found no relations for content %s", leadContentCP.uuid)
//	assert.Equal(t, len(expectedResponse.Contains), len(actualResponse.Contains), "Didn't get the same number of content in contains")
//	assertListContainsAll(t, actualResponse.Contains, expectedResponse.Contains)
//}
//
//func TestRetrieveContainedInContent(t *testing.T) {
//	if testing.Short() {
//		t.Skip("Short flag is set. Skipping integration test")
//	}
//	expectedResponse := relations{
//		ContainedIn: []relatedContent{
//			{leadContentCP.id, leadContentCP.apiURL},
//		},
//	}
//	db := getDatabaseConnection(t)
//	contents := []payloadData{leadContentCP, relatedContent1, relatedContent2}
//	cleanDB(t, db, allData)
//
//	writeContent(t, db, contents)
//	writeContentCollection(t, db, []payloadData{contentPackage}, "ContentPackage")
//	defer cleanDB(t, db, allData)
//
//	driver := NewCypherDriver(db)
//	actualResponse, found, err := driver.read(relatedContent2.uuid)
//	assert.NoError(t, err, "Unexpected error for content %s", relatedContent2.uuid)
//	assert.True(t, found, "Found no relations for content %s", relatedContent2.uuid)
//	assert.Equal(t, len(expectedResponse.ContainedIn), len(actualResponse.ContainedIn), "Didn't get the same number of containedIn content")
//	assertListContainsAll(t, actualResponse.ContainedIn, expectedResponse.ContainedIn)
//}

func writeContent(t testing.TB, db neoutils.NeoConnection, data []payloadData) baseftrwapp.Service {
	contentRW := content.NewCypherContentService(db)
	assert.NoError(t, contentRW.Initialise())
	for _, d := range data {
		writeJSONWithService(t, contentRW, d.path)
	}
	return contentRW
}

func writeContentCollection(t testing.TB, db neoutils.NeoConnection, data []payloadData, ccType string) collection.Service {
	contentCollectionRW := collection.NewContentCollectionService(db)
	assert.NoError(t, contentCollectionRW.Initialise())
	for _, d := range data {
		writeJSONWithContentCollectionService(t, contentCollectionRW, d.path, ccType)
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

func writeJSONWithContentCollectionService(t testing.TB, service collection.Service, pathToJSONFile string, ccType string) {
	path, err := filepath.Abs(pathToJSONFile)
	require.NoError(t, err)
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()
	require.NoError(t, err)
	dec := json.NewDecoder(f)
	inst, _, err := service.DecodeJSON(dec)
	require.NoError(t, err)
	err = service.Write(inst, ccType)
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
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false
	db, err := neoutils.Connect(url, conf)
	require.NoError(t, err, "Failed to connect to Neo4j")
	return db
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
