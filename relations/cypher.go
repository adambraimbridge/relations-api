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
	neoCRC := []neoRelatedContent{}

	//TODO Decide Curation or StoryPackage label to use to get story packages from Neo
	crcQuery := &neoism.CypherQuery{
		Statement: `
                MATCH (c:Content{uuid:{contentUUID}})<-[:IS_CURATED_FOR]-(cc:Curation)
                MATCH (cc)-[rel:SELECTS]->(t:Content)
                RETURN t.uuid as uuid
                ORDER BY rel.order
                `,
		Parameters: neoism.Props{"contentUUID": contentUUID},
		Result:     &neoCRC,
	}

	err := cd.conn.CypherBatch([]*neoism.CypherQuery{crcQuery})
	if err != nil {
		return relations{}, false, fmt.Errorf("Error querying Neo for uuid=%s, err=%v", contentUUID, err)
	}

	var found bool

	if (len(neoCRC)) != 0 {
		found = true
	}

	return cd.transformToRelations(neoCRC), found, nil
}

func (cd cypherDriver) transformToRelations(neoCRC []neoRelatedContent) relations {
	mappedRelatedContent := []relatedContent{}
	for _, neoContent := range neoCRC {
		c := relatedContent{
			APIURL: mapper.APIURL(neoContent.UUID, []string{"Content"}, "local"),
			ID:     mapper.IDURL(neoContent.UUID),
		}
		mappedRelatedContent = append(mappedRelatedContent, c)
	}

	return relations{mappedRelatedContent}
}
