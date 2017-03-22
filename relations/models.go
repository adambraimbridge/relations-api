package relations

type relations struct {
	//This is the "new" name for story packages
	CuratedRelatedContents []relatedContent `json:"curatedRelatedContent,omitempty"`
	//This is the content-package list of contained content (for a series/special report)
	Contains []relatedContent `json:"contains,omitempty"`
	//This is used to relate to content-packages that contain this content (or content-package)
	ContainedIn []relatedContent `json:"containedIn,omitempty"`
}

type relatedContent struct {
	ID     string `json:"id,omitempty"`
	APIURL string `json:"apiUrl,omitempty"`
}

type neoRelatedContent struct {
	UUID string `json:"uuid"`
}
