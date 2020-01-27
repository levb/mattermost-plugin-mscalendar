package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/server/utils"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func getCreateEventFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("create", flag.ContinueOnError)
	flagSet.Bool("help", false, "show help")
	flagSet.String("test-subject", "", "Subject of the event (no spaces for now)")
	flagSet.String("test-body", "", "Body of the event (no spaces for now)")
	flagSet.StringSlice("test-location", nil, "Location of the event <displayName,street,city,state,postalcode,country> (comma separated; no spaces)")
	flagSet.String("starttime", time.Now().Format(time.RFC3339), "Start time for the event")
	flagSet.Bool("allday", false, "Set as all day event (starttime/endtime must be set to midnight on different days - 2019-12-19T00:00:00-00:00)")
	flagSet.Int("reminder", 15, "Reminder (in minutes)")
	flagSet.String("endtime", time.Now().Add(time.Hour).Format(time.RFC3339), "End time for the event")
	flagSet.StringSlice("attendees", nil, "A comma separated list of Mattermost UserIDs")

	return flagSet
}

func (c *Command) createEvent(parameters ...string) (string, error) {

	if len(parameters) == 0 {
		return fmt.Sprintf(getCreateEventFlagSet().FlagUsages()), nil
	}

	event, err := parseCreateArgs(parameters)
	if err != nil {
		return err.Error(), nil
	}

	createFlagSet := getCreateEventFlagSet()
	err = createFlagSet.Parse(parameters)
	if err != nil {
		return "", err
	}

	mattermostUserIDs, err := createFlagSet.GetStringSlice("attendees")
	if err != nil {
		return "", err
	}

	calEvent, err := c.MSCalendar.CreateEvent(event, mattermostUserIDs)
	if err != nil {
		return "", err
	}
	resp := "Event Created\n" + utils.JSONBlock(&calEvent)

	return resp, nil
}

func parseCreateArgs(args []string) (*remote.Event, error) {

	event := &remote.Event{}

	createFlagSet := getCreateEventFlagSet()
	err := createFlagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	// check for required flags
	requiredFlags := []string{"test-subject"}
	flags := make(map[string]bool)
	createFlagSet.Visit(
		func(f *flag.Flag) {
			flags[f.Name] = true
		})
	for _, req := range requiredFlags {
		if !flags[req] {
			return nil, errors.Errorf("Missing required flag: `--%s` ", req)
		}
	}

	help, err := createFlagSet.GetBool("help")
	if help == true {
		return nil, errors.Errorf(getCreateEventFlagSet().FlagUsages())
	}

	subject, err := createFlagSet.GetString("test-subject")
	if err != nil {
		return nil, err
	}
	// check that next arg is not a flag "--"
	if strings.HasPrefix(subject, "--") {
		return nil, errors.Errorf("test-subject flag requires an argument")
	}
	event.Subject = subject

	body, err := createFlagSet.GetString("test-body")
	if err != nil {
		return nil, err
	}
	// check that next arg is not a flag "--"
	if strings.HasPrefix(body, "--") {
		return nil, errors.Errorf("body flag requires an argument")
	}
	event.Body = &remote.ItemBody{
		Content: body,
	}

	startTime, err := createFlagSet.GetString("starttime")
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(startTime, "--") {
		return nil, errors.Errorf("starttime flag requires an argument")
	}
	event.Start = &remote.DateTime{
		DateTime: startTime,
		TimeZone: "Pacific Standard Time",
	}

	endTime, err := createFlagSet.GetString("endtime")
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(endTime, "--") {
		return nil, errors.Errorf("endtime flag requires an argument")
	}
	event.End = &remote.DateTime{
		DateTime: endTime,
		TimeZone: "Pacific Standard Time",
	}

	allday, err := createFlagSet.GetBool("allday")
	if err != nil {
		return nil, err
	}
	event.IsAllDay = allday

	reminder, err := createFlagSet.GetInt("reminder")
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(strconv.Itoa(int(reminder)), "--") {
		return nil, errors.Errorf("reminder flag requires an argument")
	}
	event.ReminderMinutesBeforeStart = reminder

	location, err := createFlagSet.GetStringSlice("test-location")
	if err != nil {
		return nil, err
	}
	if len(location) != 0 {
		if len(location) != 6 {
			return nil, errors.Errorf("test-location flag requires 6 parameters, including a comma for empty values")
		}
		event.Location = &remote.Location{
			LocationType: "default",
			DisplayName:  location[0],
			Address: &remote.Address{
				Street:          location[1],
				City:            location[2],
				State:           location[3],
				PostalCode:      location[4],
				CountryOrRegion: location[5],
			},
		}
	}

	return event, nil
}