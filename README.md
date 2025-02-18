# namagent

namagent is a cross-platform live streaming alert CLI application written in Golang.

## Download
[Releases Page](https://github.com/shinosaki/namagent/releases)

## Usage

## Global Options
- `--config`: Specify the path to the configuration file (by default, try to load `./config.yaml`)

### alert

```bash
namagent alert [options]
```

### recorder

```bash
namagent recorder [options] [URL or ProgramId]
```

## Configs

See more example of [config.yaml](./config.yaml)

- `following.nico`: Array of user IDs to be automatically recorded
- `alert.check_interval_sec`: Interval in seconds for periodic monitoring (by default, `10` sec)
- `recorder.output_template`: Template for the output file name (by default, `{yyyymmdd}-{id}-{providerId}-{title}`)
  - Example Outputs:
    - `20250101-lv1234-1234-title.ts`
    - `20250101-lv1234-1234-title.json`
- `recorder.command_template`: Command array executed during recording, with placeholders for cookies, URL, and output file (by default: `["ffmpeg", "-cookies", "{cookies}", "-i", "{url}", "-c", "copy", "{output}"]`)

## LICENSE
[MIT](./LICENSE)

## Author
shinosaki https://shinosaki.com
