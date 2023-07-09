package tclient

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
)

// Settings is the settings of a team received from the API
type Settings struct {
	Focus int64 `json:"focus"`
	Pause int64 `json:"pause"`
}

// UpdateSettingsResponse is the response received when updating settings
type UpdateSettingsResponse struct {
	Href     string   `json:"href"`
	Settings Settings `json:"Settings"`
}

// UpdateSettings updates the settings for a team
func (c *Client) UpdateSettings(
	ctx context.Context,
	team string,
	focus int64,
	pause int64) (*UpdateSettingsResponse, error) {
	u, err := url.JoinPath(c.httpBaseURL, urlTeamSlug, team, urlSettingsSlug)
	if err != nil {
		return nil, err
	}

	body := Settings{
		Focus: focus,
		Pause: pause,
	}

	var bBody bytes.Buffer

	if err := c.createRequestBody(&body, &bBody); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, &bBody)
	if err != nil {
		return nil, err
	}

	res := UpdateSettingsResponse{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
