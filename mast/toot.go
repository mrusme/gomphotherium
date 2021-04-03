package mast

import (
  // "context"

  "github.com/mattn/go-mastodon"
)

type Toot struct {
  client                     *mastodon.Client

  ID                         int
  Status                     mastodon.Status
}

func NewToot(mastodonClient *mastodon.Client, mastodonStatus *mastodon.Status , id int) Toot {
  toot := Toot{
    client: mastodonClient,

    ID: id,
    Status: *mastodonStatus,
  }

  return toot
}
