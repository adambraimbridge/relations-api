package relations

import (
	"fmt"
	"github.com/Financial-Times/neo-model-utils-go/mapper"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
)

type Driver interface {
	read(UUID string) (res relations, found bool, err error)
	checkConnectivity() error
}

type cypherDriver struct {
	conn neoutils.NeoConnection
}

func NewCypherDriver(conn neoutils.NeoConnection) cypherDriver {
	return cypherDriver{conn}
}

func (cd cypherDriver) checkConnectivity() error {
	return neoutils.Check(cd.conn)
}

func (cd cypherDriver) read(contentUUID string) (relations, bool, error) {
	//neo curated related content a.k.a. (former) story package
	neoCRC := []neoRelatedContent{}
	//neo content package contained contents
	neoCPContains := []neoRelatedContent{}
	//neo contained in contents
	neoCPContainedIn := []neoRelatedContent{}

	//TODO Decide Curation or StoryPackage label to use to get story packages from Neo
	queryCRC := &neoism.CypherQuery{
		Statement: `
                MATCH (c:Content{uuid:{contentUUID}})<-[:IS_CURATED_FOR]-(cc:Curation)
                MATCH (cc)-[rel:SELECTS]->(t:Content)
                RETURN t.uuid as uuid
                ORDER BY rel.order
                `,
		Parameters: neoism.Props{"contentUUID": contentUUID},
		Result:     &neoCRC,
	}

	queryCPContains := &neoism.CypherQuery{
		Statement: `
                MATCH (cp:ContentPackage{uuid:{contentUUID}})-[:CONTAINS]->(cc:ContentCollection)
                MATCH (cc)-[rel:CONTAINS]->(c:Content)
                RETURN c.uuid as uuid
                ORDER BY rel.order
                `,
		Parameters: neoism.Props{"contentUUID": contentUUID},
		Result:     &neoCPContains,
	}

	queryCPContainedIn := &neoism.CypherQuery{
		Statement: `
                MATCH (c:Content{uuid:{contentUUID}})<-[:CONTAINS]-(cc:ContentCollection)
                MATCH (cc)<-[rel:CONTAINS]-(cp:ContentPackage)
                RETURN cp.uuid as uuid
                ORDER BY rel.order
                `,
		Parameters: neoism.Props{"contentUUID": contentUUID},
		Result:     &neoCPContainedIn,
	}

	err := cd.conn.CypherBatch([]*neoism.CypherQuery{queryCRC, queryCPContains, queryCPContainedIn})
	if err != nil {
		return relations{}, false, fmt.Errorf("Error querying Neo for uuid=%s, err=%v", contentUUID, err)
	}

	var found bool

	if len(neoCRC) != 0 || len(neoCPContains) != 0 || len(neoCPContainedIn) != 0 {
		found = true
	}

	mappedCRC := cd.transformToRelatedContent(neoCRC)
	mappedCPC := cd.transformToRelatedContent(neoCPContains)
	mappedCIC := cd.transformToRelatedContent(neoCPContainedIn)
	relations := relations{mappedCRC, mappedCPC, mappedCIC}

	return relations, found, nil
}

func (cd cypherDriver) transformToRelatedContent(neoRelatedContent []neoRelatedContent) []relatedContent {
	mappedRelatedContent := []relatedContent{}
	for _, neoContent := range neoRelatedContent {
		c := relatedContent{
			APIURL: mapper.APIURL(neoContent.UUID, []string{"Content"}, "local"),
			ID:     mapper.IDURL(neoContent.UUID),
		}
		mappedRelatedContent = append(mappedRelatedContent, c)
	}

	return mappedRelatedContent
}
