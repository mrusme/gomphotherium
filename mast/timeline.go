package mast

import (
  "context"

  "github.com/mattn/go-mastodon"
)

type TimelineType int
const (
  TimelineHome TimelineType  = 0
  TimelineLocal              = 1
  TimelinePublic             = 2
  TimelineNotifications      = 3
  TimelineEnd                = 4
)

type Timeline struct {
  client                     *mastodon.Client

  LastRenderedIndex          int

  Type                       TimelineType
  Account                    mastodon.Account
  Toots                      []mastodon.Status
  TootIDs                    map[string]int
}

func NewTimeline(mastodonClient *mastodon.Client) Timeline {
  timeline := Timeline{
    client: mastodonClient,

    LastRenderedIndex: -1,
    TootIDs: make(map[string]int),
  }

  return timeline
}

func (timeline *Timeline) Load(timelineType TimelineType) (error) {
  var statuses []*mastodon.Status
  var err error

  account, err := timeline.client.GetAccountCurrentUser(context.Background())
  if err != nil {
    return err
  }

  timeline.Account = *account

  switch timelineType {
  case TimelineHome:
    statuses, err = timeline.client.GetTimelineHome(context.Background(), nil)
  case TimelineLocal:
    statuses, err = timeline.client.GetTimelinePublic(context.Background(), true, nil)
  case TimelinePublic:
    statuses, err = timeline.client.GetTimelinePublic(context.Background(), false, nil)
  case TimelineNotifications:
    notifications, err := timeline.client.GetNotifications(context.Background(), nil)
    if err != nil {
      return err
    }

    for _, notification := range notifications {
      statuses = append(statuses, notification.Status)
    }
  }

  if err != nil {
    return err
  }

  for _, status := range statuses {
    id := string(status.ID)
    _, exists := timeline.TootIDs[id]
    if exists == false {
      timeline.Toots = append(timeline.Toots, *status)
      timeline.TootIDs[id] = (len(timeline.Toots) - 1)
    }
  }

  return nil
}
