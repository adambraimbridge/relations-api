package relations

import "github.com/Financial-Times/neo-model-utils-go/mapper"

func transformToRelatedContent(neoRelatedContent []neoRelatedContent) []relatedContent {
	var mappedRelatedContent []relatedContent
	for _, neoContent := range neoRelatedContent {
		c := relatedContent{
			APIURL: mapper.APIURL(neoContent.UUID, []string{"Content"}, "local"),
			ID:     mapper.IDURL(neoContent.UUID),
		}
		mappedRelatedContent = append(mappedRelatedContent, c)
	}

	return mappedRelatedContent
}
