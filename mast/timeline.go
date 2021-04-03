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
  timelineType               TimelineType

  LastRenderedIndex          int

  Account                    mastodon.Account
  Toots                      []Toot
  TootIndexStatusIDMappings  map[string]int
  KnownUsers                 []string
}

func NewTimeline(mastodonClient *mastodon.Client) Timeline {
  timeline := Timeline{
    client: mastodonClient,
    timelineType: TimelineHome,

    LastRenderedIndex: -1,
    TootIndexStatusIDMappings: make(map[string]int),
  }

  return timeline
}

func (timeline *Timeline) Switch(timelineType TimelineType) {
  if timeline.timelineType != timelineType {
    timeline.timelineType = timelineType
    timeline.Toots = []Toot{}
    timeline.TootIndexStatusIDMappings = make(map[string]int)
    timeline.LastRenderedIndex = -1
  }
}

func (timeline *Timeline) GetCurrentType() (TimelineType) {
  return timeline.timelineType
}

func (timeline *Timeline) Load() (error) {
  var statuses []*mastodon.Status
  var err error

  account, err := timeline.client.GetAccountCurrentUser(context.Background())
  if err != nil {
    return err
  }

  timeline.Account = *account

  switch timeline.timelineType {
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

  oldestStatusIndex := len(statuses) - 1
  for i := oldestStatusIndex; i >= 0; i-- {
    status := statuses[i]

    id := string(status.ID)
    _, exists := timeline.TootIndexStatusIDMappings[id]
    if exists == false {
      tootIndex := len(timeline.Toots)
      timeline.Toots = append(timeline.Toots, NewToot(timeline.client, status, tootIndex))
      timeline.KnownUsers = append(timeline.KnownUsers, status.Account.Acct) // TODO: Deduplicate
      timeline.TootIndexStatusIDMappings[id] = tootIndex
    }
  }

  return nil
}

func (timeline *Timeline) Toot(status *string, inReplyTo int, filesToUpload []string, visibility *string, sensitive bool, spoiler *string) (*mastodon.Status, error) {
  newToot := mastodon.Toot{
    Status: *status,
    Visibility: *visibility,
    Sensitive: sensitive,
    SpoilerText: *spoiler,
  }

  if inReplyTo > -1 {
    newToot.InReplyToID = timeline.Toots[inReplyTo].Status.ID
  }

  if len(filesToUpload) > 0 {
    var mediaIDs []mastodon.ID

    for _, fileToUpload := range filesToUpload {
      attachment, err := timeline.client.UploadMedia(context.Background(), fileToUpload)
      if err != nil {
        return nil, err
      }

      mediaIDs = append(mediaIDs, attachment.ID)
    }

    newToot.MediaIDs = mediaIDs
  }

  return timeline.client.PostStatus(context.Background(), &newToot)
}
