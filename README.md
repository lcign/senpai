![](res/icon.128.png) ![](res/logo.png)

*Welcome home, desune~*

**A modern terminal IRC client.**

![a screenshot of your senpai feat. simon!](senpai.png)

senpai is an IRC client that works best with bouncers:

- no logs are kept,
- history is fetched from the server via [CHATHISTORY],
- networks are fetched from the server via [bouncer-networks],
- messages can be searched in logs via [SEARCH],
- files can be uploaded via [FILEHOST] (with drag & drop!)

## senpai-on-ghostty — what's different

This is a fork of [senpai](https://git.sr.ht/~delthas/senpai) with extra features built around [Ghostty](https://ghostty.org) on macOS.

### URL indicators

Every link in the chat gets a visual prefix:

| Indicator | Type |
|-----------|------|
| 🖼 | Image (jpg, png, gif…) — click to preview inline |
| 🎬 | Video (mp4, mov, webm…) — click to open in QuickTime |
| 📄 | Markdown file (.md) — click to open in built-in viewer |
| 🔗 | Any other URL — click to open in browser |

### Scroll shortcuts

`Shift+↑` and `Shift+↓` scroll the chat history up and down without leaving the input field.

### Click to move cursor

Click anywhere on the input line to move the text cursor to that position.

### Copy mode

Press `F9` to enter copy mode and select text from the chat history.

| Key | Action |
|-----|--------|
| `F9` (or `Option+S`*) | Enter / exit copy mode |
| `↑` / `↓` | Move cursor line by line |
| click | Jump cursor to clicked line |
| `v` | Start / extend selection |
| `y` | Copy selected text to clipboard |
| `Esc` | Exit without copying |

Selected lines are highlighted in blue. The copied text goes to the system clipboard (`pbcopy` on macOS). Paste into the input with `Cmd+V` as usual.

### Video preview

Click a 🎬 link and senpai downloads the file to a temp location and opens it with the system default player (QuickTime on macOS) — no browser involved. The temp file is removed automatically when you close the player.

### Markdown viewer

Click a 📄 link to open the file in a built-in full-screen viewer — no browser needed.

GitHub and GitLab blob URLs are automatically rewritten to their raw equivalents before downloading, so links like `https://github.com/user/repo/blob/main/README.md` work directly.

| Key | Action |
|-----|--------|
| `↑` / `↓` | Scroll 3 lines |
| `PgUp` / `PgDn` | Scroll 20 lines |
| `Esc` | Close viewer |

Rendered elements: headings (`#`, `##`, `###`), **bold**, `inline code`, fenced code blocks, bullet lists.

### /x0 upload

`/x0 <file>` uploads a file to [x0.at](https://x0.at) and pastes the resulting URL into the input.

Drag & drop a file onto the terminal to auto-fill the command. Shared URLs from x0.at are also parsed and shown with the appropriate 🖼 / 🎬 / 📄 indicator.

### /LIST improvements

`/LIST` now sorts channels by activity: the most recently active channel (highest user count at query time) appears first, so the busiest channels are always at the top of the list.

### Topic expand toggle

Long channel topics can be expanded inline — click the topic bar or use the keyboard shortcut to toggle between the truncated and full view.

---

### Ghostty configuration

By default on macOS, the Option key produces special characters (Option+S = ß) and is not forwarded to the app as Alt. To enable `Option+S` for copy mode, add to `~/.config/ghostty/config`:

```
macos-option-as-alt = left
```

`F9` works immediately without any configuration change.

---

## Quick demo

To try out senpai "online", a live SSH demo is available at:
```shell
ssh -p 6666 delthas.fr
```

Your nick will be set to your SSH username.

*(This connects to the Ergo test network.)*

## Installing

- From your system package repositories: [`senpai`](https://repology.org/project/senpai-irc/versions)
- Windows binary: [senpai-0.5.0.exe](https://git.sr.ht/~delthas/senpai/refs/download/v0.5.0/senpai-0.5.0.exe)
- From source (requires Go):
```shell
git clone https://git.sr.ht/~delthas/senpai
cd senpai
make
sudo make install
```

## Running

From your terminal:
```shell
senpai
```
Senpai will guide you through a configuration assistant on your first run.

Then, type `/join #senpai` on [Libera.Chat] and have a... chat!

See `doc/senpai.1.scd` for more information and `doc/senpai.5.scd` for more
configuration options!

## Debugging errors, testing servers

To debug IRC traffic, run senpai with the `-debug` argument (or put `debug true`) in your config, it will then print in the `home` buffer all the data it sends and receives.

## Issue tracker

Browse tickets at <https://todo.sr.ht/~delthas/senpai>.

To create a ticket, visit the page above, or simply send an email to: [u.delthas.senpai@todo.sr.ht](mailto:u.delthas.senpai@todo.sr.ht) (does not require an account)

## Contributing

Sending patches to senpai is done [by email](https://lists.sr.ht/~delthas/senpai-dev), this is simple and built-in to Git.

### Using pyonji

[pyonji](https://git.sr.ht/~emersion/pyonji) streamlines the Git email contribution workflow.

Install, then after adding your changes to a commit, run `pyonji`.

### Using traditional git tools

Set up your system once by following the steps Installation and Configuration of [git-send-email.io](https://git-send-email.io/)

Then, run once in this repository:
```shell
git config sendemail.to "~delthas/senpai-dev@lists.sr.ht"
```

Then, to send a patch, make your commit, then run:
```shell
git send-email --base=HEAD~1 --annotate -1 -v1
```

It should then appear on [the mailing list](https://lists.sr.ht/~delthas/senpai-dev/patches).

## License

This senpai is open source! Please use it under the ISC license.

Copyright (C) 2021 The senpai Contributors

*senpai was created by taiite, who later handed development over to delthas. This is not a fork, but a continuation of the project initially hosted at https://sr.ht/~taiite/senpai/*

[bouncer-networks]: https://git.sr.ht/~emersion/soju/tree/master/item/doc/ext/bouncer-networks.md
[CHATHISTORY]: https://ircv3.net/specs/extensions/chathistory
[SEARCH]: https://github.com/ircv3/ircv3-specifications/pull/496
[FILEHOST]: https://codeberg.org/emersion/soju/src/branch/master/doc/ext/filehost.md
[Libera.Chat]: https://libera.chat/
[ml]: https://lists.sr.ht/~delthas/senpai-dev
