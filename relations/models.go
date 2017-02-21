package relations

type relations struct {
	//This is the "new" name for story packages
	CuratedRelatedContents []relatedContent `json:"curatedRelatedContent,omitempty"`
}

type relatedContent struct {
	ID     string `json:"id,omitempty"`
	APIURL string `json:"apiUrl,omitempty"`
}

type neoRelatedContent struct {
	UUID string `json:"uuid"`
}
