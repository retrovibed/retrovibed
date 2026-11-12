// Package me manages the authenticated profile itself. the functions take an http client
// authorized for the cloud service, see authn.AutoJWTClient.
package me

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/env"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

// StatusError the api responded with an error status.
type StatusError struct {
	Code   int
	Status string
}

func (t *StatusError) Error() string {
	return t.Status
}

// BadRequest reports whether err is the api rejecting the request, for example an invalid email.
func BadRequest(err error) bool {
	var status *StatusError
	return errors.As(err, &status) && status.Code == http.StatusBadRequest
}

// Forbidden reports whether err is the api refusing the change, for example disabling the current profile.
func Forbidden(err error) bool {
	var status *StatusError
	return errors.As(err, &status) && status.Code == http.StatusForbidden
}

// NotFound reports whether err is the api not finding the profile for the ssh identity.
func NotFound(err error) bool {
	var status *StatusError
	return errors.As(err, &status) && status.Code == http.StatusNotFound
}

func endpoint() string {
	return fmt.Sprintf("https://%s/me", env.Deeppool())
}

// Profiles lists every profile associated with the ssh identity of the profile of the token, including
// disabled profiles.
func Profiles(ctx context.Context, c *http.Client) ([]*authn.Profile, error) {
	var out authn.ProfileSearchResponse

	if err := do(ctx, c, http.MethodGet, endpoint()+"/profile/", nil, &out, "unable to list profiles"); err != nil {
		return nil, err
	}

	return out.Items, nil
}

// UpdateProfile updates the profile of the token, the profile is never part of the request.
// the display and email replace the stored values, they are sent as given.
func UpdateProfile(ctx context.Context, c *http.Client, display string, email string) (*authn.Profile, error) {
	var out authn.ProfileUpdateResponse

	encoded, err := json.Marshal(&authn.ProfileUpdateRequest{Profile: &authn.Profile{Display: display, Email: email}})
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to encode request")
	}

	if err = do(ctx, c, http.MethodPatch, endpoint()+"/", encoded, &out, "unable to update profile"); err != nil {
		return nil, err
	}

	return out.Profile, nil
}

// EnableProfile enables the profile, it must share an ssh identity with the profile of the token.
func EnableProfile(ctx context.Context, c *http.Client, id string) (*authn.Profile, error) {
	return profileAccess(ctx, c, http.MethodPut, id, "unable to enable profile")
}

// DisableProfile disables the profile, it must share an ssh identity with the profile of the token.
// the profile of the token cannot be disabled.
func DisableProfile(ctx context.Context, c *http.Client, id string) (*authn.Profile, error) {
	return profileAccess(ctx, c, http.MethodDelete, id, "unable to disable profile")
}

func profileAccess(ctx context.Context, c *http.Client, method string, id string, msg string) (*authn.Profile, error) {
	var out authn.ProfileUpdateResponse

	if err := do(ctx, c, method, fmt.Sprintf("%s/profile/%s", endpoint(), url.PathEscape(id)), nil, &out, msg); err != nil {
		return nil, err
	}

	return out.Profile, nil
}

func do(ctx context.Context, c *http.Client, method string, uri string, body []byte, out any, msg string) error {
	req, err := http.NewRequestWithContext(ctx, method, uri, bytes.NewReader(body))
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(req)
	if err != nil {
		return errorsx.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return errorsx.Wrap(&StatusError{Code: resp.StatusCode, Status: resp.Status}, msg)
	}

	if err = json.NewDecoder(resp.Body).Decode(out); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	return nil
}
