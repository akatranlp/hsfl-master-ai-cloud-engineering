package chapters_controller

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

type Controller interface {
	GetChaptersForBook(http.ResponseWriter, *http.Request)
	GetChapterForBook(http.ResponseWriter, *http.Request)

	PostChapter(http.ResponseWriter, *http.Request)
	PatchChapter(http.ResponseWriter, *http.Request)
	DeleteChapter(http.ResponseWriter, *http.Request)

	ValidateChapterId(http.ResponseWriter, *http.Request)

	LoadChapterMiddleware(http.ResponseWriter, *http.Request, router.Next)
}
