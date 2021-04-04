# Gomphotherium

Gomphotherium (*/ˌɡɒmfəˈθɪəriəm/*; "welded beast"), a command line Mastodon 
client.


## Description

Gomphotherium is a Mastodon client for the command line, offering a CLI as well
as a TUI with a usage similar to [rainbowstream](rainbowstream).

[rainbowstream]: https://github.com/orakaro/rainbowstream


## Installation

Download a binary from the [releases][releases] page. MacOS, Linux, Windows,
FreeBSD, and OpenBSD binaries are available.

Or just build it yourself (requires Go 1.16+):

```bash
make
```

[releases]: https://github.com/mrusme/gomphotherium/releases


## User Manual


### Authentication

To authenticate with your Mastodon instance, run the following command and
follow the instructions:

```sh
gomphotherium authenticate https://YOUR-MASTODON-SERVER-URL-HERE.com
```


### CLI

TODO


### TUI

Launch the TUI with the following command:

```sh
gomphotherium tui
```

**Note:** If you haven't exported the required environment variables that were
shown to you during the [Authentication][#authentication], please do so first
or use the CLI flags (`gomphotherium -h`) instead.


#### Modes

The TUI can be operated in two modes: **Normal** and **Command**.

In **Normal** mode no interaction is possible apart from scrolling and 
refreshing the timeline and quitting Gomphotherium. The shortcuts can be looked
up on the [cheatsheet](#cheatsheet)

In **Command** mode, the command input becomes available and scrolling the
timeline is not possible anymore. Commands can then be issued to interact with
the Mastodon instance.


#### Cheatsheet


##### Shortcuts

`Ctrl+Q`: Quit Gomphotherium

`Ctrl+R`: Refresh timeline

`i`: Enter **Command** mode (while in **Normal** mode)

`Esc`: Leave **Command** mode (while in **Command** mode)


##### Commands


