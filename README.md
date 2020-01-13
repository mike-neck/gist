gist
===

A tool for pulling gist as git repository into a user configured directory.

Configuration
===

A user can configure `gist` with `$HOME/.gist.yml` file.

* `profile` - determines which context to use.(mandatory)
* `github_access_token` - git hub access token(requires `gist` scope). If this value is not set, an environmental value `GITHUB_ACCESS_TOKEN` will be used.
* `destination_dir` - a directory relative to user home where `gist` clones gist repositories.(default `gist/{profile}` which means gist repositories will be cloned at `$HOME/gist/{profile}`)

```yaml
- profile: default
  destination_dir: /users/foo/my-gists
  # GITHUB_ACCESS_TOKEN will be used for the profile "default".
- profile: privates
  github_access_token: 5f4e3d2c1b0a
  # $HOME/gist/privates will be used for the profile "privates".
```

Commands
===

List cloned gists
---

*WIP*

* command - `list`
* parameters
    * `profile` - Profile to use.(Default: `default`)
    * `output` - Output format.(Default: `json`. Available: `json`, `xml`, `yaml`, `csv`, `tsv`)
    * `limit` - Size of pages.(Default: `20`. If `0` is given, all gists will be shown)
    * `page` - A position of pages.(Default: `1`)
    * `sort` - A sort order of gists.(Default: `pub-desc`. Available: `pub-desc`, `pub-asc`, `id-desc`, `id-asc`)

#### Example

```bash
gist list -profile privates -output csv -limit 10 -page 3
```

Clone gist
---

*WIP*

* command - `clone`
* parameters
    * An id of gist
    * `profile` - Profile to use.(Default: `default`)
    * `user` - A user name of gist.
    * `ssh` - Prefer ssh(Default: `false` = `https`)

#### Example

```bash
gist clone 0a1b2c3d4e5f -ssh
```

Profile
---

Creates or overrides a profile.

* command - `profile`
* parameters
    * `profile` - A name of new profile.(If not specified, `default` will be used.)
    * `token` - GitHub access token for the new profile to use.
    * `dir` - Destination directory for the new profile to use.

```bash
gist profile -name privates -token f5e4d3c2b1a0 
```

Of cause you can add a github access token manually, after creation of a new profile.
