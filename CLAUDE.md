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
