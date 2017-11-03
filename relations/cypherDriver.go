package relations

import (
	"fmt"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
)

type Driver interface {
	findContentRelations(UUID string) (res relations, found bool, err error)
	findContentCollectionRelations(UUID string) (res relations, found bool, err error)
	checkConnectivity() error
}

type cypherDriver struct {
	conn neoutils.NeoConnection
}

func NewCypherDriver(conn neoutils.NeoConnection) *cypherDriver {
	return &cypherDriver{conn}
}

func (cd *cypherDriver) checkConnectivity() error {
	return neoutils.Check(cd.conn)
}

func (cd *cypherDriver) findContentRelations(contentUUID string) (relations, bool, error) {
	neoCRC := []neoRelatedContent{}
	neoCPContains := []neoRelatedContent{}
	neoCPContainedIn := []neoRelatedContent{}

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

	mappedCRC := transformToRelatedContent(neoCRC)
	mappedCPC := transformToRelatedContent(neoCPContains)
	mappedCIC := transformToRelatedContent(neoCPContainedIn)
	relations := relations{mappedCRC, mappedCPC, mappedCIC}

	return relations, found, nil
}

func (cd *cypherDriver) findContentCollectionRelations(contentCollectionUUID string) (relations, bool, error) {
	neoCPContainedIn := []neoRelatedContent{}
	neoCPContains := []neoRelatedContent{}

	queryCPContainedIn := &neoism.CypherQuery{
		Statement: `
                MATCH (cc:ContentCollection{uuid:{contentCollectionUUID}})<-[:CONTAINS]-(cp:ContentPackage)
                RETURN cp.uuid as uuid
                ORDER BY rel.order
                `,
		Parameters: neoism.Props{"contentCollectionUUID": contentCollectionUUID},
		Result:     &neoCPContainedIn,
	}

	queryCPContains := &neoism.CypherQuery{
		Statement: `
                MATCH (cc:ContentCollection{uuid:{contentCollectionUUID}})-[:CONTAINS]->(c:Content)
                RETURN c.uuid as uuid
                ORDER BY rel.order
                `,
		Parameters: neoism.Props{"contentCollectionUUID": contentCollectionUUID},
		Result:     &neoCPContains,
	}

	err := cd.conn.CypherBatch([]*neoism.CypherQuery{queryCPContains, queryCPContainedIn})
	if err != nil {
		return relations{}, false, fmt.Errorf("Error querying Neo for uuid=%s, err=%v", contentCollectionUUID, err)
	}

	var found bool
	if len(neoCPContains) != 0 || len(neoCPContainedIn) != 0 {
		found = true
	}

	mappedCPC := transformToRelatedContent(neoCPContains)
	mappedCIC := transformToRelatedContent(neoCPContainedIn)
	relations := relations{nil, mappedCPC, mappedCIC}

	return relations, found, nil
}
