package tvdb

import (
	"github.com/pkg/errors"

	"github.com/dashotv/tvdb/openapi/models/operations"
)

// GetArtworkBase wraps the generated openapi.SDK.Artwork.GetArtworkBase call
func (c *Client) GetArtworkBase(id float64) (*GetArtworkBaseResponse, error) {
	r, err := c.sdk.Artwork.GetArtworkBase(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetArtworkExtended wraps the generated openapi.SDK.Artwork.GetArtworkExtended call
func (c *Client) GetArtworkExtended(id float64) (*GetArtworkExtendedResponse, error) {
	r, err := c.sdk.Artwork.GetArtworkExtended(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllArtworkStatuses wraps the generated openapi.SDK.ArtworkStatuses.GetAllArtworkStatuses call
func (c *Client) GetAllArtworkStatuses() (*GetAllArtworkStatusesResponse, error) {
	r, err := c.sdk.ArtworkStatuses.GetAllArtworkStatuses(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllArtworkTypes wraps the generated openapi.SDK.ArtworkTypes.GetAllArtworkTypes call
func (c *Client) GetAllArtworkTypes() (*GetAllArtworkTypesResponse, error) {
	r, err := c.sdk.ArtworkTypes.GetAllArtworkTypes(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAwardCategory wraps the generated openapi.SDK.AwardCategories.GetAwardCategory call
func (c *Client) GetAwardCategory(id float64) (*GetAwardCategoryResponse, error) {
	r, err := c.sdk.AwardCategories.GetAwardCategory(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAwardCategoryExtended wraps the generated openapi.SDK.AwardCategories.GetAwardCategoryExtended call
func (c *Client) GetAwardCategoryExtended(id float64) (*GetAwardCategoryExtendedResponse, error) {
	r, err := c.sdk.AwardCategories.GetAwardCategoryExtended(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllAwards wraps the generated openapi.SDK.Awards.GetAllAwards call
func (c *Client) GetAllAwards() (*GetAllAwardsResponse, error) {
	r, err := c.sdk.Awards.GetAllAwards(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAward wraps the generated openapi.SDK.Awards.GetAward call
func (c *Client) GetAward(id float64) (*GetAwardResponse, error) {
	r, err := c.sdk.Awards.GetAward(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAwardExtended wraps the generated openapi.SDK.Awards.GetAwardExtended call
func (c *Client) GetAwardExtended(id float64) (*GetAwardExtendedResponse, error) {
	r, err := c.sdk.Awards.GetAwardExtended(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetCharacterBase wraps the generated openapi.SDK.Characters.GetCharacterBase call
func (c *Client) GetCharacterBase(id float64) (*GetCharacterBaseResponse, error) {
	r, err := c.sdk.Characters.GetCharacterBase(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllCompanies wraps the generated openapi.SDK.Companies.GetAllCompanies call
func (c *Client) GetAllCompanies(page *float64) (*GetAllCompaniesResponse, error) {
	r, err := c.sdk.Companies.GetAllCompanies(c.ctx, page)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetCompany wraps the generated openapi.SDK.Companies.GetCompany call
func (c *Client) GetCompany(id float64) (*GetCompanyResponse, error) {
	r, err := c.sdk.Companies.GetCompany(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetCompanyTypes wraps the generated openapi.SDK.Companies.GetCompanyTypes call
func (c *Client) GetCompanyTypes() (*GetCompanyTypesResponse, error) {
	r, err := c.sdk.Companies.GetCompanyTypes(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllContentRatings wraps the generated openapi.SDK.ContentRatings.GetAllContentRatings call
func (c *Client) GetAllContentRatings() (*GetAllContentRatingsResponse, error) {
	r, err := c.sdk.ContentRatings.GetAllContentRatings(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllCountries wraps the generated openapi.SDK.Countries.GetAllCountries call
func (c *Client) GetAllCountries() (*GetAllCountriesResponse, error) {
	r, err := c.sdk.Countries.GetAllCountries(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetEntityTypes wraps the generated openapi.SDK.EntityTypes.GetEntityTypes call
func (c *Client) GetEntityTypes() (*GetEntityTypesResponse, error) {
	r, err := c.sdk.EntityTypes.GetEntityTypes(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllEpisodes wraps the generated openapi.SDK.Episodes.GetAllEpisodes call
func (c *Client) GetAllEpisodes(page *float64) (*GetAllEpisodesResponse, error) {
	r, err := c.sdk.Episodes.GetAllEpisodes(c.ctx, page)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetEpisodeBase wraps the generated openapi.SDK.Episodes.GetEpisodeBase call
func (c *Client) GetEpisodeBase(id float64) (*GetEpisodeBaseResponse, error) {
	r, err := c.sdk.Episodes.GetEpisodeBase(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetEpisodeExtended wraps the generated openapi.SDK.Episodes.GetEpisodeExtended call
func (c *Client) GetEpisodeExtended(id float64, meta *operations.Meta) (*GetEpisodeExtendedResponse, error) {
	r, err := c.sdk.Episodes.GetEpisodeExtended(c.ctx, id, meta)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetEpisodeTranslation wraps the generated openapi.SDK.Episodes.GetEpisodeTranslation call
func (c *Client) GetEpisodeTranslation(id float64, language string) (*GetEpisodeTranslationResponse, error) {
	r, err := c.sdk.Episodes.GetEpisodeTranslation(c.ctx, id, language)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllGenders wraps the generated openapi.SDK.Genders.GetAllGenders call
func (c *Client) GetAllGenders() (*GetAllGendersResponse, error) {
	r, err := c.sdk.Genders.GetAllGenders(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllGenres wraps the generated openapi.SDK.Genres.GetAllGenres call
func (c *Client) GetAllGenres() (*GetAllGenresResponse, error) {
	r, err := c.sdk.Genres.GetAllGenres(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetGenreBase wraps the generated openapi.SDK.Genres.GetGenreBase call
func (c *Client) GetGenreBase(id float64) (*GetGenreBaseResponse, error) {
	r, err := c.sdk.Genres.GetGenreBase(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllInspirationTypes wraps the generated openapi.SDK.InspirationTypes.GetAllInspirationTypes call
func (c *Client) GetAllInspirationTypes() (*GetAllInspirationTypesResponse, error) {
	r, err := c.sdk.InspirationTypes.GetAllInspirationTypes(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllLanguages wraps the generated openapi.SDK.Languages.GetAllLanguages call
func (c *Client) GetAllLanguages() (*GetAllLanguagesResponse, error) {
	r, err := c.sdk.Languages.GetAllLanguages(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllLists wraps the generated openapi.SDK.Lists.GetAllLists call
func (c *Client) GetAllLists(page *float64) (*GetAllListsResponse, error) {
	r, err := c.sdk.Lists.GetAllLists(c.ctx, page)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetList wraps the generated openapi.SDK.Lists.GetList call
func (c *Client) GetList(id float64) (*GetListResponse, error) {
	r, err := c.sdk.Lists.GetList(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetListBySlug wraps the generated openapi.SDK.Lists.GetListBySlug call
func (c *Client) GetListBySlug(slug string) (*GetListBySlugResponse, error) {
	r, err := c.sdk.Lists.GetListBySlug(c.ctx, slug)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetListExtended wraps the generated openapi.SDK.Lists.GetListExtended call
func (c *Client) GetListExtended(id float64) (*GetListExtendedResponse, error) {
	r, err := c.sdk.Lists.GetListExtended(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetListTranslation wraps the generated openapi.SDK.Lists.GetListTranslation call
func (c *Client) GetListTranslation(id float64, language string) (*GetListTranslationResponse, error) {
	r, err := c.sdk.Lists.GetListTranslation(c.ctx, id, language)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// PostLogin wraps the generated openapi.SDK.Login.PostLogin call
func (c *Client) PostLogin(request operations.PostLoginRequestBody) (*PostLoginResponse, error) {
	r, err := c.sdk.Login.PostLogin(c.ctx, request)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllMovie wraps the generated openapi.SDK.Movies.GetAllMovie call
func (c *Client) GetAllMovie(page *float64) (*GetAllMovieResponse, error) {
	r, err := c.sdk.Movies.GetAllMovie(c.ctx, page)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetMovieBase wraps the generated openapi.SDK.Movies.GetMovieBase call
func (c *Client) GetMovieBase(id float64) (*GetMovieBaseResponse, error) {
	r, err := c.sdk.Movies.GetMovieBase(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetMovieBaseBySlug wraps the generated openapi.SDK.Movies.GetMovieBaseBySlug call
func (c *Client) GetMovieBaseBySlug(slug string) (*GetMovieBaseBySlugResponse, error) {
	r, err := c.sdk.Movies.GetMovieBaseBySlug(c.ctx, slug)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetMovieExtended wraps the generated openapi.SDK.Movies.GetMovieExtended call
func (c *Client) GetMovieExtended(id float64, meta *operations.QueryParamMeta, short *bool) (*GetMovieExtendedResponse, error) {
	r, err := c.sdk.Movies.GetMovieExtended(c.ctx, id, meta, short)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetMovieTranslation wraps the generated openapi.SDK.Movies.GetMovieTranslation call
func (c *Client) GetMovieTranslation(id float64, language string) (*GetMovieTranslationResponse, error) {
	r, err := c.sdk.Movies.GetMovieTranslation(c.ctx, id, language)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetMoviesFilter wraps the generated openapi.SDK.Movies.GetMoviesFilter call
func (c *Client) GetMoviesFilter(request operations.GetMoviesFilterRequest) (*GetMoviesFilterResponse, error) {
	r, err := c.sdk.Movies.GetMoviesFilter(c.ctx, request)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllMovieStatuses wraps the generated openapi.SDK.MovieStatuses.GetAllMovieStatuses call
func (c *Client) GetAllMovieStatuses() (*GetAllMovieStatusesResponse, error) {
	r, err := c.sdk.MovieStatuses.GetAllMovieStatuses(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllPeople wraps the generated openapi.SDK.People.GetAllPeople call
func (c *Client) GetAllPeople(page *float64) (*GetAllPeopleResponse, error) {
	r, err := c.sdk.People.GetAllPeople(c.ctx, page)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetPeopleBase wraps the generated openapi.SDK.People.GetPeopleBase call
func (c *Client) GetPeopleBase(id float64) (*GetPeopleBaseResponse, error) {
	r, err := c.sdk.People.GetPeopleBase(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetPeopleExtended wraps the generated openapi.SDK.People.GetPeopleExtended call
func (c *Client) GetPeopleExtended(id float64, meta *operations.GetPeopleExtendedQueryParamMeta) (*GetPeopleExtendedResponse, error) {
	r, err := c.sdk.People.GetPeopleExtended(c.ctx, id, meta)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetPeopleTranslation wraps the generated openapi.SDK.People.GetPeopleTranslation call
func (c *Client) GetPeopleTranslation(id float64, language string) (*GetPeopleTranslationResponse, error) {
	r, err := c.sdk.People.GetPeopleTranslation(c.ctx, id, language)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllPeopleTypes wraps the generated openapi.SDK.PeopleTypes.GetAllPeopleTypes call
func (c *Client) GetAllPeopleTypes() (*GetAllPeopleTypesResponse, error) {
	r, err := c.sdk.PeopleTypes.GetAllPeopleTypes(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSearchResults wraps the generated openapi.SDK.Search.GetSearchResults call
func (c *Client) GetSearchResults(request operations.GetSearchResultsRequest) (*GetSearchResultsResponse, error) {
	r, err := c.sdk.Search.GetSearchResults(c.ctx, request)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSearchResultsByRemoteID wraps the generated openapi.SDK.Search.GetSearchResultsByRemoteID call
func (c *Client) GetSearchResultsByRemoteID(remoteID string) (*GetSearchResultsByRemoteIDResponse, error) {
	r, err := c.sdk.Search.GetSearchResultsByRemoteID(c.ctx, remoteID)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllSeasons wraps the generated openapi.SDK.Seasons.GetAllSeasons call
func (c *Client) GetAllSeasons(page *float64) (*GetAllSeasonsResponse, error) {
	r, err := c.sdk.Seasons.GetAllSeasons(c.ctx, page)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeasonBase wraps the generated openapi.SDK.Seasons.GetSeasonBase call
func (c *Client) GetSeasonBase(id float64) (*GetSeasonBaseResponse, error) {
	r, err := c.sdk.Seasons.GetSeasonBase(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeasonExtended wraps the generated openapi.SDK.Seasons.GetSeasonExtended call
func (c *Client) GetSeasonExtended(id float64) (*GetSeasonExtendedResponse, error) {
	r, err := c.sdk.Seasons.GetSeasonExtended(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeasonTypes wraps the generated openapi.SDK.Seasons.GetSeasonTypes call
func (c *Client) GetSeasonTypes() (*GetSeasonTypesResponse, error) {
	r, err := c.sdk.Seasons.GetSeasonTypes(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllSeries wraps the generated openapi.SDK.Series.GetAllSeries call
func (c *Client) GetAllSeries(page *float64) (*GetAllSeriesResponse, error) {
	r, err := c.sdk.Series.GetAllSeries(c.ctx, page)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesArtworks wraps the generated openapi.SDK.Series.GetSeriesArtworks call
func (c *Client) GetSeriesArtworks(id float64, lang *string, type_ *int64) (*GetSeriesArtworksResponse, error) {
	r, err := c.sdk.Series.GetSeriesArtworks(c.ctx, id, lang, type_)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesBase wraps the generated openapi.SDK.Series.GetSeriesBase call
func (c *Client) GetSeriesBase(id float64) (*GetSeriesBaseResponse, error) {
	r, err := c.sdk.Series.GetSeriesBase(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesBaseBySlug wraps the generated openapi.SDK.Series.GetSeriesBaseBySlug call
func (c *Client) GetSeriesBaseBySlug(slug string) (*GetSeriesBaseBySlugResponse, error) {
	r, err := c.sdk.Series.GetSeriesBaseBySlug(c.ctx, slug)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesEpisodes wraps the generated openapi.SDK.Series.GetSeriesEpisodes call
func (c *Client) GetSeriesEpisodes(request operations.GetSeriesEpisodesRequest) (*GetSeriesEpisodesResponse, error) {
	r, err := c.sdk.Series.GetSeriesEpisodes(c.ctx, request)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesExtended wraps the generated openapi.SDK.Series.GetSeriesExtended call
func (c *Client) GetSeriesExtended(id float64, meta *operations.GetSeriesExtendedQueryParamMeta, short *bool) (*GetSeriesExtendedResponse, error) {
	r, err := c.sdk.Series.GetSeriesExtended(c.ctx, id, meta, short)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesFilter wraps the generated openapi.SDK.Series.GetSeriesFilter call
func (c *Client) GetSeriesFilter(request operations.GetSeriesFilterRequest) (*GetSeriesFilterResponse, error) {
	r, err := c.sdk.Series.GetSeriesFilter(c.ctx, request)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesNextAired wraps the generated openapi.SDK.Series.GetSeriesNextAired call
func (c *Client) GetSeriesNextAired(id float64) (*GetSeriesNextAiredResponse, error) {
	r, err := c.sdk.Series.GetSeriesNextAired(c.ctx, id)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesSeasonEpisodesTranslated wraps the generated openapi.SDK.Series.GetSeriesSeasonEpisodesTranslated call
func (c *Client) GetSeriesSeasonEpisodesTranslated(id float64, lang string, page int64, seasonType string) (*GetSeriesSeasonEpisodesTranslatedResponse, error) {
	r, err := c.sdk.Series.GetSeriesSeasonEpisodesTranslated(c.ctx, id, lang, page, seasonType)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetSeriesTranslation wraps the generated openapi.SDK.Series.GetSeriesTranslation call
func (c *Client) GetSeriesTranslation(id float64, language string) (*GetSeriesTranslationResponse, error) {
	r, err := c.sdk.Series.GetSeriesTranslation(c.ctx, id, language)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllSeriesStatuses wraps the generated openapi.SDK.SeriesStatuses.GetAllSeriesStatuses call
func (c *Client) GetAllSeriesStatuses() (*GetAllSeriesStatusesResponse, error) {
	r, err := c.sdk.SeriesStatuses.GetAllSeriesStatuses(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// GetAllSourceTypes wraps the generated openapi.SDK.SourceTypes.GetAllSourceTypes call
func (c *Client) GetAllSourceTypes() (*GetAllSourceTypesResponse, error) {
	r, err := c.sdk.SourceTypes.GetAllSourceTypes(c.ctx)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}

// Updates wraps the generated openapi.SDK.Updates.Updates call
func (c *Client) Updates(since int64, action *operations.Action, page *float64, type_ *operations.Type) (*UpdatesResponse, error) {
	r, err := c.sdk.Updates.Updates(c.ctx, since, action, page, type_)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.Errorf("non-200 response: %d", r.StatusCode)
	}
	return r.Object, nil
}
