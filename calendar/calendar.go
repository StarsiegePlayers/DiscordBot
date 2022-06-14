package calendar

import (
	"context"
	"sync"
	"time"

	"github.com/StarsiegePlayers/DiscordBot/config"
	"github.com/StarsiegePlayers/DiscordBot/module"
	"github.com/StarsiegePlayers/DiscordBot/rpc"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const (
	ServiceName = "calendar"
)

type Service struct {
	module.Base

	wg     sync.WaitGroup
	config *config.CalendarConfig

	service *calendar.Service

	UpcomingEvents *calendar.Events
}

func (s *Service) Init() {
	// queue waiting on config
	s.wg.Add(1)
	s.PubSubSubscribe(rpc.NewConfigLoadedTopic, s.configMessagePubSubHandler)
}

func (s *Service) Start() (err error) {
	// wait for config
	s.wg.Wait()

	ctx, cancelfn := context.WithTimeout(s.Base, 5*time.Second)
	defer cancelfn()

	s.service, err = calendar.NewService(ctx, option.WithCredentialsJSON([]byte(s.config.AuthToken)))
	if err != nil {
		s.Logf("Unable to retrieve Calendar client: %v", err)
	}

	s.updateEvents()
	s.listEvents()

	return nil
}

func (s *Service) updateEvents() {
	t := time.Now().Format(time.RFC3339)
	events, err := s.service.Events.List(s.config.CalendarID).
		ShowDeleted(false).SingleEvents(true).
		TimeMin(t).MaxResults(s.config.NumEventLookAhead).OrderBy("startTime").
		Do()
	if err != nil {
		s.Logf("Unable to retrieve next ten of the user's events: %v", err)
	}

	s.UpcomingEvents = events
}

func (s *Service) listEvents() {
	s.Logln("Upcoming events:")
	if len(s.UpcomingEvents.Items) != 0 {
		for _, item := range s.UpcomingEvents.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			s.Logf("[%v] (%v) %v - %v\n", item.Id, date, item.Summary, item.Description)
		}
	}
}

func (s *Service) Stop() error {
	return nil
}
