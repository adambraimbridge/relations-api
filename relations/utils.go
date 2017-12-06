package relations

import "github.com/Financial-Times/neo-model-utils-go/mapper"

func transformToRelatedContent(neoRelatedContent []neoRelatedContent) []relatedContent {
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

func transformContainedInToCCRelations(containedIn []neoRelatedContent) string {
	var leadArticleUuid string
	if len(containedIn) != 0 {
		leadArticleUuid = containedIn[0].UUID
	}
	return leadArticleUuid
}

func transformContainsToCCRelations(neoRelatedContent []neoRelatedContent) []string {
	var contains []string
	for _, neoContent := range neoRelatedContent {
		contains = append(contains, neoContent.UUID)
	}
	return contains
}
