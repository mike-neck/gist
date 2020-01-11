gist
===

A tool for pulling gist as git repository into a user configured directory.

Configuration
===

A user can configure `gist` with `$HOME/.gist.yml` file.

* `profile` - determines which context to use.(mandatory)
* `github_access_token` - git hub access token(requires `gist` scope).(mandatory)
* `destination_dir` - a directory relative to user home where `gist` clones gist repositories.(default `gist/{profile}` which means gist repositories will be cloned at `$HOME/gist/{profile}`)

```yaml
- profile: default
  github_access_token: a0b1c2d3e4f5
  destination_dir: my-gists
- profile: privates
  github_access_token: 5f4e3d2c1b0a
```

Commands
===

List cloned gists
---

*WIP*

* command - `list`
* parameters
    * `output` - Output format.(Default: `json`. Available: `json`, `xml`, `yaml`, `csv`, `tsv`)
    * `limit` - Size of pages.(Default: `20`. If `0` is given, all gists will be shown)
    * `page` - A position of pages.(Default: `1`)
    * `sort` - A sort order of gists.(Default: `pub-desc`. Available: `pub-desc`, `pub-asc`, `id-desc`, `id-asc`)

#### Example

```bash
gist list -output csv -limit 10 -page 3
```

Clone gist
---

*WIP*

* command - `clone`
* parameters
    * An id of gist
    * `user` - A user name of gist.

#### Example

```bash
gist clone 0a1b2c3d4e5f
```
