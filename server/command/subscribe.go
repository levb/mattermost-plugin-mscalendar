// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/utils"
)

func (c *Command) subscribe(parameters ...string) (string, error) {
	switch {
	case len(parameters) == 0:
		storedSub, err := c.MSCalendar.CreateMyEventSubscription()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Subscription %s created.", storedSub.Remote.ID), nil

	case len(parameters) == 1 && parameters[0] == "list":
		subs, err := c.MSCalendar.ListRemoteSubscriptions()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Subscriptions:%s", utils.JSONBlock(subs)), nil

	case len(parameters) == 1 && parameters[0] == "show":
		storedSub, err := c.MSCalendar.LoadMyEventSubscription()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Subscription:%s", utils.JSONBlock(storedSub)), nil

	case len(parameters) == 1 && parameters[0] == "renew":
		storedSub, err := c.MSCalendar.RenewMyEventSubscription()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Subscription %s renewed until %s", storedSub.Remote.ID, storedSub.Remote.ExpirationDateTime), nil

	case len(parameters) == 1 && parameters[0] == "delete":
		err := c.MSCalendar.DeleteMyEventSubscription()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("User's subscription  deleted"), nil

	case len(parameters) == 2 && parameters[0] == "delete":
		err := c.MSCalendar.DeleteOrphanedSubscription(parameters[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Subscription %s deleted", parameters[1]), nil

	}
	return "bad syntax", nil
}
