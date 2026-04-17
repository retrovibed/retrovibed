package metaapi

import "github.com/retrovibed/retrovibed/retroapi/authn"

type Account = authn.Account
type AccountLookupRequest = authn.AccountLookupRequest
type AccountLookupResponse = authn.AccountLookupResponse
type AccountUpdateRequest = authn.AccountUpdateRequest
type AccountUpdateResponse = authn.AccountUpdateResponse

type ProfileStatus = authn.ProfileStatus

const ProfileStatus_NONE = authn.ProfileStatus_NONE
const ProfileStatus_DISABLED = authn.ProfileStatus_DISABLED
const ProfileStatus_ENABLED = authn.ProfileStatus_ENABLED
const ProfileStatus_PENDING = authn.ProfileStatus_PENDING

type Profile = authn.Profile
type ProfileSearchRequest = authn.ProfileSearchRequest
type ProfileSearchResponse = authn.ProfileSearchResponse
type ProfileLookupRequest = authn.ProfileLookupRequest
type ProfileLookupResponse = authn.ProfileLookupResponse
type ProfileUpdateRequest = authn.ProfileUpdateRequest
type ProfileUpdateResponse = authn.ProfileUpdateResponse
type ProfileDisableRequest = authn.ProfileDisableRequest
type ProfileDisableResponse = authn.ProfileDisableResponse
type ProfileCreateRequest = authn.ProfileCreateRequest
type ProfileCreateResponse = authn.ProfileCreateResponse

type Bearer = authn.Bearer
type Token = authn.Token
type AuthzRequest = authn.AuthzRequest
type AuthzResponse = authn.AuthzResponse
type AuthzGrantRequest = authn.AuthzGrantRequest
type AuthzGrantResponse = authn.AuthzGrantResponse
type AuthzRevokeRequest = authn.AuthzRevokeRequest
type AuthzRevokeResponse = authn.AuthzRevokeResponse
type AuthzProfileRequest = authn.AuthzProfileRequest
type AuthzProfileResponse = authn.AuthzProfileResponse

type LoginOptions = authn.LoginOptions
type LoginOptionsResponse = authn.LoginOptionsResponse
type Identity = authn.Identity
type Authn = authn.Authn
type Authed = authn.Authed
type Session = authn.Session
