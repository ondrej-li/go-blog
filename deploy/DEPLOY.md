# go.ondralinek.cz

Static Hugo build of the Go blog. No PHP, no database, no runtime — nginx
serves plain files, the source lives in git and the build runs on the VPS.

| | |
|---|---|
| Host | go.ondralinek.cz |
| Web root | `/var/www/go-blog` (deploy:www-data) |
| Source checkout | `~/go-blog` (deploy user, clone of github.com/ondrej-li/go-blog) |
| vhost | `/etc/nginx/sites-available/go-blog` → `sites-enabled/go-blog` |
| headers | `/etc/nginx/snippets/go-blog-headers.conf` (HSTS + CSP + …) |
| TLS | `/etc/letsencrypt/live/go.ondralinek.cz/` (Let's Encrypt, expires 2027-02-07) |
| renewal | `certbot.timer`; deploy hook reloads nginx |
| logs | `/var/log/nginx/go-blog.access.log`, `go-blog.error.log` |
| build tool | Hugo **0.146.0** in `/usr/local/bin/hugo` (pinned) |
| auto-deploy | `go-blog-autodeploy.timer` (5 min) → `deploy/autodeploy.sh` |

Everything nginx-related is versioned in this repository under `deploy/`; the
server copies are installed by `deploy/deploy.sh`.

## Deploy an update

On the VPS, as `deploy`:

```bash
cd ~/go-blog
./deploy/deploy.sh                 # git pull, build, publish, nginx config
SKIP_NGINX=1 ./deploy/deploy.sh    # content only
```

The web root is replaced with `rsync --delete`, so anything that must survive a
deploy belongs outside `/var/www/go-blog`.

## Automatic deploys

`go-blog-autodeploy.timer` runs `deploy/autodeploy.sh` every 5 minutes:

1. `git fetch` (nothing is pulled or built when `origin/master` is unchanged)
2. new commit → `deploy/deploy.sh` (pull, build, rsync, reload nginx if needed)
3. no new commit → no-op, except one rebuild per hour (`MAX_AGE_SECONDS`) so
   future-dated posts appear once their date passes

Properties:

- `flock` prevents overlapping runs (timer vs. a manual `deploy.sh`)
- a dirty server checkout is refused instead of being silently overwritten
- if the build fails, `deploy.sh` exits before rsync, so the previous build
  stays live; the failure is mailed to `admin@ondralinek.cz` (through local
  Postfix, at most once per commit per hour)
- journal output: `journalctl -u go-blog-autodeploy -n 100`

```bash
sudo systemctl start go-blog-autodeploy.service   # deploy now, don't wait for the timer
systemctl list-timers go-blog-autodeploy.timer    # next run
sudo systemctl disable --now go-blog-autodeploy.timer  # turn it off
```

The units are versioned in `deploy/systemd/` and (re)installed by `deploy.sh`;
`SKIP_SYSTEMD=1` skips that step.

## Publishing an article

Commit a new Markdown file under `content/post/` and push to `master`. Hugo
generates `/post/<file-name>/` from the front matter; no manual build is needed,
the timer picks the commit up within 5 minutes.

The published copy in `docs/` (GitHub Pages) is *not* used by this host — the
VPS builds from `content/` and the theme in `themes/`. Run `hugo` locally only
to preview.

## Why Hugo is pinned to 0.146.0

The theme and the existing articles were written for this Hugo generation; newer
releases changed URL and taxonomy output. `deploy.sh` refuses to build with a
different version, so upgrade the VPS and your machine together and check the
URLs first:

```bash
# install / upgrade on the VPS
curl -sSL -o /tmp/hugo.tar.gz \
  https://github.com/gohugoio/hugo/releases/download/v0.146.0/hugo_0.146.0_linux-arm64.tar.gz
tar xzf /tmp/hugo.tar.gz -C /tmp hugo
sudo install -m 0755 /tmp/hugo /usr/local/bin/hugo
```

## DNS cutover (GitHub Pages → this server)

`go.ondralinek.cz` is a **CNAME to `ondrej-li.github.io`** (GitHub Pages) until
the cutover. Replace it — a CNAME cannot coexist with the A/AAAA records:

| Type | Name | Value |
|---|---|---|
| A | go.ondralinek.cz | `2.31.33.21` |
| AAAA | go.ondralinek.cz | `2a01:4f8:c016:e828::1` |

DNS is managed at WebSupport (`ns1.websupport.cz`). The GitHub Pages site stays
untouched, so reverting the record is the rollback plan.

## TLS

Issued on 2026-11-09 and live (HTTP/2, HSTS, HTTP→HTTPS 301). Run the cutover
script only after DNS points here (e.g. when rebuilding the host):

```bash
cd ~/go-blog
./deploy/enable-tls.sh
```

It checks DNS via public resolvers, requests the certificate with the **webroot**
authenticator (`-w /var/www/go-blog`, so the versioned vhost is never rewritten),
installs `deploy/nginx/go-blog.conf` and verifies both schemes.
`deploy/deploy.sh` keeps the ACME challenge directory intact via `rsync
--exclude`, and selects the HTTP or HTTPS vhost depending on whether the
certificate exists.

```bash
sudo certbot renew --cert-name go.ondralinek.cz --dry-run --no-random-sleep-on-renew
```

## Verification

```bash
curl -sSI https://go.ondralinek.cz/ | head -5
curl -sS  https://go.ondralinek.cz/post/ -o /dev/null -w '%{http_code}\n'
curl -sSI https://go.ondralinek.cz/post/mage/ | grep -i cache-control   # no-cache
curl -sSI https://go.ondralinek.cz/css/theme.min.css | grep -i cache-control
curl -sSI http://go.ondralinek.cz/post/ | grep -i location             # -> https
```
