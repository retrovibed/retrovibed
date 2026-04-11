package tvdb

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var tvdbApiKey = os.Getenv("TVDB_API_KEY")
var tvdbToken = os.Getenv("TVDB_API_TOKEN")

var artwork_id_int64 = int64(23165)
var artwork_id_float64 = float64(23165)
var awardcategories_id_int64 = int64(2)
var awardcategories_id_float64 = float64(2)
var awards_id_int64 = int64(1)
var awards_id_float64 = float64(1)
var characters_id_int64 = int64(65407033)
var characters_id_float64 = float64(65407033)
var companies_id_int64 = int64(1305)
var companies_id_float64 = float64(1305)
var episodes_id_int64 = int64(55452)
var episodes_id_float64 = float64(55452)
var episodes_language_string = "eng"
var genres_id_int64 = int64(15)
var genres_id_float64 = float64(15)
var lists_id_int64 = int64(14232)
var lists_id_float64 = float64(14232)
var lists_language_string = "eng"
var lists_slug_string = "the-simpsons"
var movies_id_int64 = int64(504)
var movies_id_float64 = float64(504)
var movies_language_string = "eng"
var movies_slug_string = "midnight-sun"
var people_id_int64 = int64(296524)
var people_id_float64 = float64(296524)
var people_language_string = "eng"
var search_remoteID_string = "tt0096697"
var seasons_id_int64 = int64(2727)
var seasons_id_float64 = float64(2727)
var seasons_language_string = "eng"
var series_id_int64 = int64(71663)
var series_id_float64 = float64(71663)
var series_lang_string = "eng"
var series_language_string = "eng"
var series_page_int64 = int64(0)
var series_seasonType_string = "default"
var series_slug_string = "the-simpsons"
var updates_since_int64 = time.Now().Add(-24 * time.Hour).Unix()
var updates_since_float64 = float64(time.Now().Add(-24*time.Hour).Unix() * 1000)
var userInfo_id_int64 = int64(1)

var operations_PostLoginRequestBody = PostLoginRequestBody{
	Apikey: tvdbApiKey,
}
var operations_GetMoviesFilterRequest = GetMoviesFilterRequest{
	Country: "US",
	Lang:    movies_language_string,
}
var operations_GetSearchResultsRequest = GetSearchResultsRequest{
	Query: String("The Simpsons"),
}
var operations_GetSeriesEpisodesRequest = GetSeriesEpisodesRequest{
	ID:         series_id_float64,
	Page:       series_page_int64,
	SeasonType: series_seasonType_string,
}
var operations_GetSeriesFilterRequest = GetSeriesFilterRequest{
	Country: "US",
	Lang:    series_language_string,
}

func testClient(t *testing.T) *Client {
	assert.NotEmpty(t, tvdbApiKey, "TVDB_API_KEY is empty")

	c, err := Login(tvdbApiKey)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	return c
}
