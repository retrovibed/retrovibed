# retrodsearch

A [PeerTube](https://joinpeertube.org/) search plugin for
[retrovibed](https://github.com/retrovibed/retrovibed)'s
`searchplugin.Registry`. It queries a PeerTube instance's public JSON search
API (by default [sepiasearch.org](https://sepiasearch.org), a cross-instance
search index run by the PeerTube project) and emits results as newline
delimited JSON, one `ddiscapi.Import` per line.

## Install

Use retrovibed's builtin plugin installer, which compiles this module to
wasip1/wasm and drops it into the search plugin directory retrovibed already
watches. `retrodsearch` lives inside the `retrovibed` repo as its own Go
module (a sibling of `retroapi` and `shallows`), so point the installer at
its local directory:

```
retrovibe ddisc search plugin install ./retrodsearch --name peertube
```

Point it at a different PeerTube/SepiaSearch-compatible instance by passing
`-e PEERTUBE_DOMAIN=https://your-instance.example`.

## Development

This is also a normal CLI tool, useful for testing against a live instance
without installing it as a plugin:

```
go run . --query ubuntu
go run . --query ubuntu --adult
```

## Flags

| flag | default | description |
| --- | --- | --- |
| `--query` | (required) | search text |
| `--category` | `all` | one of `all`, `movies`, `tv`, `music`, `games` |
| `--domain` | `https://sepiasearch.org` | base url of the PeerTube/SepiaSearch instance (env `PEERTUBE_DOMAIN`) |
| `--adult` | `false` | allow adult (NSFW) results; set per-search by `searchplugin.Registry`, not install-time configuration |
| `--max-results` | `128` | maximum results to fetch detail pages for |
| `--attempts` | `5` | maximum retry attempts per request |
| `--rate` | `1` | maximum requests per second against the target domain |
| `--workers` | `4` | number of concurrent detail-page fetches |
| `--output` | `-` | output destination; `-` for stdout |
