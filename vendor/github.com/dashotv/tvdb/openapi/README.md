# github.com/dashotv/tvdb/openapi

<div align="left">
    <a href="https://speakeasyapi.dev/"><img src="https://custom-icon-badges.demolab.com/badge/-Built%20By%20Speakeasy-212015?style=for-the-badge&logoColor=FBE331&logo=speakeasy&labelColor=545454" /></a>
    <a href="https://opensource.org/licenses/MIT">
        <img src="https://img.shields.io/badge/License-MIT-blue.svg" style="width: 100px; height: 28px;" />
    </a>
</div>


## 🏗 **Welcome to your new SDK!** 🏗

It has been generated successfully based on your OpenAPI spec. However, it is not yet ready for production use. Here are some next steps:
- [ ] 🛠 Make your SDK feel handcrafted by [customizing it](https://www.speakeasyapi.dev/docs/customize-sdks)
- [ ] ♻️ Refine your SDK quickly by iterating locally with the [Speakeasy CLI](https://github.com/speakeasy-api/speakeasy)
- [ ] 🎁 Publish your SDK to package managers by [configuring automatic publishing](https://www.speakeasyapi.dev/docs/advanced-setup/publish-sdks)
- [ ] ✨ When ready to productionize, delete this section from the README

<!-- Start SDK Installation [installation] -->
## SDK Installation

```bash
go get github.com/dashotv/tvdb/openapi
```
<!-- End SDK Installation [installation] -->

<!-- Start SDK Example Usage [usage] -->
## SDK Example Usage

### Example

```go
package main

import (
	"context"
	"github.com/dashotv/tvdb/openapi"
	"github.com/dashotv/tvdb/openapi/models/shared"
	"log"
)

func main() {
	s := openapi.New(
		openapi.WithSecurity(shared.Security{
			BearerAuth: "<YOUR_BEARER_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.ArtworkStatuses.GetAllArtworkStatuses(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.Object != nil {
		// handle response
	}
}

```
<!-- End SDK Example Usage [usage] -->

<!-- Start Available Resources and Operations [operations] -->
## Available Resources and Operations

### [ArtworkStatuses](docs/sdks/artworkstatuses/README.md)

* [GetAllArtworkStatuses](docs/sdks/artworkstatuses/README.md#getallartworkstatuses) - Returns list of artwork status records.

### [ArtworkTypes](docs/sdks/artworktypes/README.md)

* [GetAllArtworkTypes](docs/sdks/artworktypes/README.md#getallartworktypes) - Returns a list of artworkType records

### [Artwork](docs/sdks/artwork/README.md)

* [GetArtworkBase](docs/sdks/artwork/README.md#getartworkbase) - Returns a single artwork base record.
* [GetArtworkExtended](docs/sdks/artwork/README.md#getartworkextended) - Returns a single artwork extended record.

### [Awards](docs/sdks/awards/README.md)

* [GetAllAwards](docs/sdks/awards/README.md#getallawards) - Returns a list of award base records
* [GetAward](docs/sdks/awards/README.md#getaward) - Returns a single award base record
* [GetAwardExtended](docs/sdks/awards/README.md#getawardextended) - Returns a single award extended record

### [AwardCategories](docs/sdks/awardcategories/README.md)

* [GetAwardCategory](docs/sdks/awardcategories/README.md#getawardcategory) - Returns a single award category base record
* [GetAwardCategoryExtended](docs/sdks/awardcategories/README.md#getawardcategoryextended) - Returns a single award category extended record

### [Characters](docs/sdks/characters/README.md)

* [GetCharacterBase](docs/sdks/characters/README.md#getcharacterbase) - Returns character base record

### [Companies](docs/sdks/companies/README.md)

* [GetAllCompanies](docs/sdks/companies/README.md#getallcompanies) - returns a paginated list of company records
* [GetCompany](docs/sdks/companies/README.md#getcompany) - returns a company record
* [GetCompanyTypes](docs/sdks/companies/README.md#getcompanytypes) - returns all company type records

### [ContentRatings](docs/sdks/contentratings/README.md)

* [GetAllContentRatings](docs/sdks/contentratings/README.md#getallcontentratings) - returns list content rating records

### [Countries](docs/sdks/countries/README.md)

* [GetAllCountries](docs/sdks/countries/README.md#getallcountries) - returns list of country records

### [EntityTypes](docs/sdks/entitytypes/README.md)

* [GetEntityTypes](docs/sdks/entitytypes/README.md#getentitytypes) - returns the active entity types

### [Episodes](docs/sdks/episodes/README.md)

* [GetAllEpisodes](docs/sdks/episodes/README.md#getallepisodes) - Returns a list of episodes base records with the basic attributes.<br> Note that all episodes are returned, even those that may not be included in a series' default season order.
* [GetEpisodeBase](docs/sdks/episodes/README.md#getepisodebase) - Returns episode base record
* [GetEpisodeExtended](docs/sdks/episodes/README.md#getepisodeextended) - Returns episode extended record
* [GetEpisodeTranslation](docs/sdks/episodes/README.md#getepisodetranslation) - Returns episode translation record

### [Genders](docs/sdks/genders/README.md)

* [GetAllGenders](docs/sdks/genders/README.md#getallgenders) - returns list of gender records

### [Genres](docs/sdks/genres/README.md)

* [GetAllGenres](docs/sdks/genres/README.md#getallgenres) - returns list of genre records
* [GetGenreBase](docs/sdks/genres/README.md#getgenrebase) - Returns genre record

### [InspirationTypes](docs/sdks/inspirationtypes/README.md)

* [GetAllInspirationTypes](docs/sdks/inspirationtypes/README.md#getallinspirationtypes) - returns list of inspiration types records

### [Languages](docs/sdks/languages/README.md)

* [GetAllLanguages](docs/sdks/languages/README.md#getalllanguages) - returns list of language records

### [Lists](docs/sdks/lists/README.md)

* [GetAllLists](docs/sdks/lists/README.md#getalllists) - returns list of list base records
* [GetList](docs/sdks/lists/README.md#getlist) - returns an list base record
* [GetListBySlug](docs/sdks/lists/README.md#getlistbyslug) - returns an list base record search by slug
* [GetListExtended](docs/sdks/lists/README.md#getlistextended) - returns a list extended record
* [GetListTranslation](docs/sdks/lists/README.md#getlisttranslation) - Returns list translation record

### [Login](docs/sdks/login/README.md)

* [PostLogin](docs/sdks/login/README.md#postlogin) - create an auth token. The token has one month validation length.

### [Movies](docs/sdks/movies/README.md)

* [GetAllMovie](docs/sdks/movies/README.md#getallmovie) - returns list of movie base records
* [GetMovieBase](docs/sdks/movies/README.md#getmoviebase) - Returns movie base record
* [GetMovieBaseBySlug](docs/sdks/movies/README.md#getmoviebasebyslug) - Returns movie base record search by slug
* [GetMovieExtended](docs/sdks/movies/README.md#getmovieextended) - Returns movie extended record
* [GetMovieTranslation](docs/sdks/movies/README.md#getmovietranslation) - Returns movie translation record
* [GetMoviesFilter](docs/sdks/movies/README.md#getmoviesfilter) - Search movies based on filter parameters

### [MovieStatuses](docs/sdks/moviestatuses/README.md)

* [GetAllMovieStatuses](docs/sdks/moviestatuses/README.md#getallmoviestatuses) - returns list of status records

### [People](docs/sdks/people/README.md)

* [GetAllPeople](docs/sdks/people/README.md#getallpeople) - Returns a list of people base records with the basic attributes.
* [GetPeopleBase](docs/sdks/people/README.md#getpeoplebase) - Returns people base record
* [GetPeopleExtended](docs/sdks/people/README.md#getpeopleextended) - Returns people extended record
* [GetPeopleTranslation](docs/sdks/people/README.md#getpeopletranslation) - Returns people translation record

### [PeopleTypes](docs/sdks/peopletypes/README.md)

* [GetAllPeopleTypes](docs/sdks/peopletypes/README.md#getallpeopletypes) - returns list of peopleType records

### [Search](docs/sdks/search/README.md)

* [GetSearchResults](docs/sdks/search/README.md#getsearchresults) - Our search index includes series, movies, people, and companies. Search is limited to 5k results max.
* [GetSearchResultsByRemoteID](docs/sdks/search/README.md#getsearchresultsbyremoteid) - Search a series, movie, people, episode, company or season by specific remote id and returns a base record for that entity.

### [Seasons](docs/sdks/seasons/README.md)

* [GetAllSeasons](docs/sdks/seasons/README.md#getallseasons) - returns list of seasons base records
* [GetSeasonBase](docs/sdks/seasons/README.md#getseasonbase) - Returns season base record
* [GetSeasonExtended](docs/sdks/seasons/README.md#getseasonextended) - Returns season extended record
* [GetSeasonTranslation](docs/sdks/seasons/README.md#getseasontranslation) - Returns season translation record
* [GetSeasonTypes](docs/sdks/seasons/README.md#getseasontypes) - Returns season type records

### [Series](docs/sdks/series/README.md)

* [GetAllSeries](docs/sdks/series/README.md#getallseries) - returns list of series base records
* [GetSeriesArtworks](docs/sdks/series/README.md#getseriesartworks) - Returns series artworks base on language and type. <br> Note&#58; Artwork type is an id that can be found using **/artwork/types** endpoint.
* [GetSeriesBase](docs/sdks/series/README.md#getseriesbase) - Returns series base record
* [GetSeriesBaseBySlug](docs/sdks/series/README.md#getseriesbasebyslug) - Returns series base record searched by slug
* [GetSeriesEpisodes](docs/sdks/series/README.md#getseriesepisodes) - Returns series episodes from the specified season type, default returns the episodes in the series default season type
* [GetSeriesExtended](docs/sdks/series/README.md#getseriesextended) - Returns series extended record
* [GetSeriesFilter](docs/sdks/series/README.md#getseriesfilter) - Search series based on filter parameters
* [GetSeriesNextAired](docs/sdks/series/README.md#getseriesnextaired) - Returns series base record including the nextAired field. <br> Note&#58; nextAired was included in the base record endpoint but that field will deprecated in the future so developers should use the nextAired endpoint.
* [GetSeriesSeasonEpisodesTranslated](docs/sdks/series/README.md#getseriesseasonepisodestranslated) - Returns series base record with episodes from the specified season type and language. Default returns the episodes in the series default season type.
* [GetSeriesTranslation](docs/sdks/series/README.md#getseriestranslation) - Returns series translation record

### [SeriesStatuses](docs/sdks/seriesstatuses/README.md)

* [GetAllSeriesStatuses](docs/sdks/seriesstatuses/README.md#getallseriesstatuses) - returns list of status records

### [SourceTypes](docs/sdks/sourcetypes/README.md)

* [GetAllSourceTypes](docs/sdks/sourcetypes/README.md#getallsourcetypes) - returns list of sourceType records

### [Updates](docs/sdks/updates/README.md)

* [Updates](docs/sdks/updates/README.md#updates) - Returns updated entities.  methodInt indicates a created record (1), an updated record (2), or a deleted record (3).  If a record is deleted because it was a duplicate of another record, the target record's information is provided in mergeToType and mergeToId.

### [UserInfo](docs/sdks/userinfo/README.md)

* [GetUserInfo](docs/sdks/userinfo/README.md#getuserinfo) - returns user info
* [GetUserInfoByID](docs/sdks/userinfo/README.md#getuserinfobyid) - returns user info by user id

### [Favorites](docs/sdks/favorites/README.md)

* [CreateUserFavorites](docs/sdks/favorites/README.md#createuserfavorites) - creates a new user favorite
* [GetUserFavorites](docs/sdks/favorites/README.md#getuserfavorites) - returns user favorites
<!-- End Available Resources and Operations [operations] -->

<!-- Start Error Handling [errors] -->
## Error Handling

Handling errors in this SDK should largely match your expectations.  All operations return a response object or an error, they will never return both.  When specified by the OpenAPI spec document, the SDK will return the appropriate subclass.

| Error Object       | Status Code        | Content Type       |
| ------------------ | ------------------ | ------------------ |
| sdkerrors.SDKError | 4xx-5xx            | */*                |

### Example

```go
package main

import (
	"context"
	"errors"
	"github.com/dashotv/tvdb/openapi"
	"github.com/dashotv/tvdb/openapi/models/sdkerrors"
	"github.com/dashotv/tvdb/openapi/models/shared"
	"log"
)

func main() {
	s := openapi.New(
		openapi.WithSecurity(shared.Security{
			BearerAuth: "<YOUR_BEARER_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.ArtworkStatuses.GetAllArtworkStatuses(ctx)
	if err != nil {

		var e *sdkerrors.SDKError
		if errors.As(err, &e) {
			// handle error
			log.Fatal(e.Error())
		}
	}
}

```
<!-- End Error Handling [errors] -->

<!-- Start Server Selection [server] -->
## Server Selection

### Select Server by Index

You can override the default server globally using the `WithServerIndex` option when initializing the SDK client instance. The selected server will then be used as the default on the operations that use it. This table lists the indexes associated with the available servers:

| # | Server | Variables |
| - | ------ | --------- |
| 0 | `https://api4.thetvdb.com/v4` | None |

#### Example

```go
package main

import (
	"context"
	"github.com/dashotv/tvdb/openapi"
	"github.com/dashotv/tvdb/openapi/models/shared"
	"log"
)

func main() {
	s := openapi.New(
		openapi.WithServerIndex(0),
		openapi.WithSecurity(shared.Security{
			BearerAuth: "<YOUR_BEARER_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.ArtworkStatuses.GetAllArtworkStatuses(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.Object != nil {
		// handle response
	}
}

```


### Override Server URL Per-Client

The default server can also be overridden globally using the `WithServerURL` option when initializing the SDK client instance. For example:
```go
package main

import (
	"context"
	"github.com/dashotv/tvdb/openapi"
	"github.com/dashotv/tvdb/openapi/models/shared"
	"log"
)

func main() {
	s := openapi.New(
		openapi.WithServerURL("https://api4.thetvdb.com/v4"),
		openapi.WithSecurity(shared.Security{
			BearerAuth: "<YOUR_BEARER_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.ArtworkStatuses.GetAllArtworkStatuses(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.Object != nil {
		// handle response
	}
}

```
<!-- End Server Selection [server] -->

<!-- Start Custom HTTP Client [http-client] -->
## Custom HTTP Client

The Go SDK makes API calls that wrap an internal HTTP client. The requirements for the HTTP client are very simple. It must match this interface:

```go
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
```

The built-in `net/http` client satisfies this interface and a default client based on the built-in is provided by default. To replace this default with a client of your own, you can implement this interface yourself or provide your own client configured as desired. Here's a simple example, which adds a client with a 30 second timeout.

```go
import (
	"net/http"
	"time"
	"github.com/myorg/your-go-sdk"
)

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	sdkClient  = sdk.New(sdk.WithClient(httpClient))
)
```

This can be a convenient way to configure timeouts, cookies, proxies, custom headers, and other low-level configuration.
<!-- End Custom HTTP Client [http-client] -->

<!-- Start Authentication [security] -->
## Authentication

### Per-Client Security Schemes

This SDK supports the following security scheme globally:

| Name         | Type         | Scheme       |
| ------------ | ------------ | ------------ |
| `BearerAuth` | http         | HTTP Bearer  |

You can configure it using the `WithSecurity` option when initializing the SDK client instance. For example:
```go
package main

import (
	"context"
	"github.com/dashotv/tvdb/openapi"
	"github.com/dashotv/tvdb/openapi/models/shared"
	"log"
)

func main() {
	s := openapi.New(
		openapi.WithSecurity(shared.Security{
			BearerAuth: "<YOUR_BEARER_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.ArtworkStatuses.GetAllArtworkStatuses(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.Object != nil {
		// handle response
	}
}

```
<!-- End Authentication [security] -->

<!-- Start Special Types [types] -->
## Special Types


<!-- End Special Types [types] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->

# Development

## Maturity

This SDK is in beta, and there may be breaking changes between versions without a major version update. Therefore, we recommend pinning usage
to a specific package version. This way, you can install the same version each time without breaking changes unless you are intentionally
looking for the latest version.

## Contributions

While we value open-source contributions to this SDK, this library is generated programmatically. Any manual changes added to internal files will be overwritten on the next generation. 
We look forward to hearing your feedback. Feel free to open a PR or an issue with a proof of concept and we'll do our best to include it in a future release. 

### SDK Created by [Speakeasy](https://docs.speakeasyapi.dev/docs/using-speakeasy/client-sdks)
