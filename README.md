# Introduction

This is basically a quick way of organizing videos based on filtering rules.

# Command-line Options

```
$ nyaa_copy -h
usage: nyaa_copy [-h|--help] [-c|--config "<value>"] [--cache "<value>"]
                 [-e|--expired] [-d|--dry-run] [--debug] [--version]

                 A custom copying program that can relabel videos, add season
                 information, and track the number of times a given filter was
                 used.

Arguments:

  -h  --help     Print help information
  -c  --config   Use custom configuration file. Default:
                 /home/jweatherly/.config/nyaa_copy/config.yaml
      --cache    Use custom cache file. Default:
                 /home/jweatherly/.cache/nyaa_copy/cache.json
  -e  --expired  Delete old filter files and exit
  -d  --dry-run  Do not actually copy files, just simulate. Default: false
      --debug    Set logging level to DEBUG
      --version  Display version and exit
```

# Configuration

The main configuration file:

```yaml
directories:
  shows: /directory/with/filters
  videos: /directory/with/videos
  destination: /directory/to/copy/to
allowed_extensions:
  - .mkv
  - .mp4
  - .avi
```

- `directories`: Contains all of the directory locations required for the program
  - `shows`: The location for the show filter files
  - `videos`: The location of the videos to copy/organize
  - `destination`: The destination "root" directory.  This will be the base path for every video.
- `allowed_extensions`: Ignore any file that does not have one of these extensions

# Filter File

Each show is defined by a `yaml` file that contains what it's looking for, where it is located in the standard configuration file's `destination` directory, season and episode offset information.

```yaml
search: "Test Show"
destination: "Test Show"
rename: "Test Show (English)"
season: 2
offset: -13
```

- `search`: The name of the show in the file.  For example, if the show is `Test Show - 20.mkv`, then the `search` should be `Test Show`.
- `destination`: The folder in the configuration's `destination` directory where the show will live.
- `rename`: Rename the show when it's copied over.
- `season|1`: The season of the show.
- `offset|0`: A number to add to the episode number so the math makes sense.  For example, if episode 20 of **Test Show** is really the 7th episode of Season 2, then you'd set `season` to **2** and `offset` to **-13**.
