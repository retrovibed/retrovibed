package deeppool

import (
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
)

type Archiver = deeppool.Archiver
type Retrieval = deeppool.Retrieval
type Ranger = deeppool.Ranger
type Media = deeppool.Media
type MediaSearchRequest = deeppool.MediaSearchRequest
type MediaSearchResponse = deeppool.MediaSearchResponse
type MediaCreateRequest = deeppool.MediaCreateRequest
type MediaCreateResponse = deeppool.MediaCreateResponse
type MediaUploadRequest = deeppool.MediaUploadRequest
type MediaUploadResponse = deeppool.MediaUploadResponse
type MediaDownloadRequest = deeppool.MediaDownloadRequest
type MediaCompletedRequest = deeppool.MediaCompletedRequest
type MediaCompletedResponse = deeppool.MediaCompletedResponse
type MediaFindRequest = deeppool.MediaFindRequest
type MediaFindResponse = deeppool.MediaFindResponse
type MediaDeleteRequest = deeppool.MediaDeleteRequest
type MediaDeleteResponse = deeppool.MediaDeleteResponse

var NewArchiver = deeppool.NewArchiver
var NewRetrieval = deeppool.NewRetrieval
var NewRanger = deeppool.NewRanger
