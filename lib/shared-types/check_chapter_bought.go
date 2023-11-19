package shared_types

type CheckChapterBoughtRequest struct {
	UserID  	 uint64 `json:"userId"`
	ChapterID 	 uint64 `json:"chapterId"`
}

func (r *CheckChapterBoughtRequest) IsValid() bool {
	return r.ChapterID != 0
}

type CheckChapterBoughtResponse struct {
	Success bool `json:"success"`
}
