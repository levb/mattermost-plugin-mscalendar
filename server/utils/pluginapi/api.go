// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package pluginapi

import (
	"strings"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/store"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type API struct {
	api plugin.API
}

func New(api plugin.API) *API {
	return &API{
		api: api,
	}
}

func (a *API) GetMattermostUserStatus(userID string) (*model.Status, error) {
	return a.api.GetUserStatus(userID)
}

func (a *API) GetMattermostUserStatusesByIds(userIDs []string) ([]*model.Status, error) {
	return a.api.GetUserStatusesByIds(userIDs)
}

func (a *API) UpdateMattermostUserStatus(userID, status string) (*model.Status, error) {
	return a.api.UpdateUserStatus(userID, status)
}

// IsPluginAdmin returns true if the user is authorized to use the workflow plugin's admin-level APIs/commands.
func (a *API) IsSysAdmin(mattermostUserID string) (bool, error) {
	user, err := a.api.GetUser(mattermostUserID)
	if err != nil {
		return false, err
	}
	return strings.Contains(user.Roles, "system_admin"), nil
}

func (a *API) GetMattermostUserByUsername(mattermostUsername string) (*model.User, error) {
	for strings.HasPrefix(mattermostUsername, "@") {
		mattermostUsername = mattermostUsername[1:]
	}
	u, err := a.api.GetUserByUsername(mattermostUsername)
	if err != nil {
		return nil, err
	}
	if u.DeleteAt != 0 {
		return nil, store.ErrNotFound
	}
	return u, nil
}

func (a *API) GetMattermostUser(mattermostUserID string) (*model.User, error) {
	mmuser, err := a.api.GetUser(mattermostUserID)
	if err != nil {
		return nil, err
	}
	if mmuser.DeleteAt != 0 {
		return nil, store.ErrNotFound
	}
	return mmuser, nil
}

func (a *API) CleanKVStore() error {
	appErr := a.api.KVDeleteAll()
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a *API) SendEphemeralPost(channelID, userID, message string) {
	ephemeralPost := &model.Post{
		ChannelId: channelID,
		UserId:    userID,
		Message:   message,
	}
	_ = a.api.SendEphemeralPost(userID, ephemeralPost)
}