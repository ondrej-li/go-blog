# CLAUDE.md — go-blog

A Hugo blog in Czech about Go, published at **go.ondralinek.cz**.
[`deploy/DEPLOY.md`](deploy/DEPLOY.md) is the deployment map (host, web root,
vhost, TLS, the pinned Hugo); this file is the working map.

## Layout

| Path | What |
|---|---|
| `content/post/*.md` | the posts. Czech, front matter with `title`, `date`, `tags` |
| `themes/hugo-code-editor-theme/` | the theme, vendored. Do not edit it to style one post |
| `static/` | files served as-is |
| `config.toml` | site config: baseURL, language `cs-CZ`, theme, permalinks |
| `makefile` | `make publish` → `hugo -d ./docs` |
| `docs/` | **generated** output, committed, and what GitHub Pages serves |
| `deploy/` | nginx vhost + headers, `deploy.sh`, `autodeploy.sh` (the VPS timer) |

## Build

Hugo **0.146.0** — the same version CI pins and the VPS has in
`/usr/local/bin/hugo`.

```sh
hugo -d ./docs        # same as `make publish`
hugo server           # preview on :1313
```

CI (`.github/workflows/ci.yml`, `ubuntu-latest`) builds with the pinned Hugo and
then fails if `docs/` differs from what that build produces. So: **edit the post,
run `make publish`, commit `content/` and `docs/` together.** A red CI here means
exactly one thing — the published copy is stale.

## Deploy

Two channels, both automatic:

- **The VPS.** `go-blog-autodeploy.timer` on `deploy@ondralinek.cz` pulls `master`
  every five minutes, renders with Hugo and rsyncs into `/var/www/go-blog`.
  Nothing to do but push.
- **GitHub Pages.** Serves `docs/` from `master`. This is why `docs/` is committed.

Pushing to `master` is therefore a deploy. There is no staging.

## Conventions

- Commit messages: sentence case, imperative-ish, and the post's subject in the
  body when the change is content. `docs/` churn goes in the same commit as the
  post that caused it.
- Czech in the posts; English in code and config comments.
- Images belong in `static/` (or a post bundle), and keep them small — the box
  serves them without a CDN.
- Do not reformat `docs/`: it is generated HTML, and hand-editing it will be
  reverted by the next `make publish`.
- The theme is vendored from upstream; local changes should be rare and explained
  in the commit, or they will be lost at the next theme update.

## Before this site goes live

This repository is served publicly at https://go.ondralinek.cz. A change that adds or alters a
page, a form, a script, a header or a vhost passes the launch gate before it
ships:

```sh
web-launch-check https://go.ondralinek.cz
```

- `FAIL` blocks the publish; `WARN` is a decision to state in the reply, not
  something to leave unmentioned.
- The `web-launch-checklist` skill has the 25 items the checker cannot judge:
  privacy and terms content, colour contrast, real-device speed, spam
  protection, analytics and error monitoring, one clear call to action, and the
  signup/login/reset flows.
- Run it again after publishing: it checks the live vhost, headers and
  certificate, not the working tree.
- A preview vhost passes the same gate before DNS moves.

## Deploying this site

https://go.ondralinek.cz is published by the **deploy** account from this repository's
`master` branch — the `site-deploy` timer, every five minutes. Merge is
deploy: there is no publish command, and `sudo vps-site publish` refuses this
web root (one writer per root, and it is not this account).

Verify after merging, and give it up to ten minutes:

```sh
site-status --site go-blog
curl -fsS https://go.ondralinek.cz/deployed.json           # {"commit": "...", "deployed_at": "..."}
git -C ~/repos/go-blog rev-parse origin/master
```

The three commits must match. If it has not moved after two timer cycles, read
`journalctl -u site-deploy -n 50`, report the commit that should be live and the
one that is, and do not rsync into `/var/www` or edit a live file to compensate.
The `static-site-deploy` skill has the full procedure.
