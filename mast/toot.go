package mast

import (
  // "context"

  "github.com/mattn/go-mastodon"
)

const (
  VisibilityPublic string    = "public"
  VisibilityPrivate          = "private"
  VisibilityUnlisted         = "unlisted"
  VisibilityDirect           = "direct"
)

type Toot struct {
  client                     *mastodon.Client

  ID                         int
  Status                     mastodon.Status

  IsNotification             bool
  Notification               mastodon.Notification
}

func NewToot(
  mastodonClient *mastodon.Client,
  mastodonStatus *mastodon.Status,
  mastodonNotification *mastodon.Notification,
  id int) Toot {
  toot := Toot{
    client: mastodonClient,

    ID: id,
    Status: *mastodonStatus,
  }

  if mastodonNotification != nil {
    toot.IsNotification = true
    toot.Notification = *mastodonNotification
  }

  return toot
}
