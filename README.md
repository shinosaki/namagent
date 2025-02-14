# namagent

namagent is a cross-platform live streaming alert CLI application written in Golang.

## Download
[Releases Page](https://github.com/shinosaki/namagent/releases)

## Usage

### alert

Starts a daemon that automatically records streams from the specified user in `config.yaml`.

```bash
namagent alert [options]
```

Options:
  - `--config`: Specify the path to the configuration file (by default, try to load `./config.yaml`)

### recorder

```bash
namagent recorder [options] url
```

Options:
  - `--ffmpeg`: Path to the ffmpeg executable file (by default, `ffmpeg`)
  - `--output`: Path for the saved file

## Config File

- `meta.fetch_interval`: Interval in seconds for periodic retrieval of program data (by default, `10` sec)
- `paths.ffmpeg`: Path to the ffmpeg executable file (by default, `ffmpeg`)
- `paths.output_base_dir`: Destination directory for recordings using **alert** (by default, working directory)
- `following.users.nico`: Array of user IDs to be automatically recorded

### Example
```yaml
meta:
  fetch_interval: 10

paths:
  ffmpeg: "ffmpeg"
  output_base_dir: "."

following:
  users:
    nico:
      - "5599432"
      - "96462240"
```

## LICENSE
[MIT](./LICENSE)

## Author
shinosaki https://shinosaki.com
