# namagent

namagent is a cross-platform live streaming alert CLI application written in Golang.

## Features
- Recording
- Alert
  - 定期的に最新の番組情報を取得します
  - 認証情報を設定している場合、WebPush通知を購読します（フォロー中のユーザの番組開始後、いち早く録画を開始できます）

## Support Sites
- live.nicovideo.jp: ニコニコ生放送
  - 動画データの取得: WebSocketサーバから取得したHLSのURLとCookieデータを、ユーザが指定したコマンドに渡して起動します
  - コメントデータの取得: MessageServerから取得したChatデータをJSON形式で保存します
  - WebPush通知の購読: フォロー中のユーザが番組を開始すると、ほぼ同時にユーザが指定したコマンドを起動します

## Download
[Releases Page](https://github.com/shinosaki/namagent/releases)

## Usage

### Global Options
- `--config`: Specify the path to the configuration file (by default, try to load `./config.yaml`)

### alert
Configファイルで指定した、対象ユーザが番組を開始した場合、自動でrecorderを実行します。

```bash
namagent alert [options]
```

### recorder
ProgramId （ニコ生の場合: `lv123...`）を含む文字列を引数として渡すと、番組データを取得し、command_templateのコマンドを実行します。

```bash
namagent recorder [options] [URL or ProgramId]
```

## Configs

See more example of [config.example.yaml](./config.example.yaml)

```yaml
alert:
  # 定期的な番組情報の取得間隔, 下限値は10s
  check_interval: 10s
auth:
  nico:
    # ニコニコのuser_sessionクッキーの値
    # 指定した場合、ログイン状態で番組データを取得します
    # また、NicoPush (WebPush通知)の購読を開始します
    user_session: user_session....
following:
  # ニコニコの対象ユーザのリスト
  # ニコニコでフォロー中でも、リストに含まれていない場合は録画されません
  nico:
    - 96254336
recorder:
  # 拡張子, デフォルトは".ts"
  extension: mp4

  # （拡張子を除く）出力ファイル名のテンプレート
  output_template: '{{.AuthorId}}/{{.StartedAt.Format "20060102"}}-{{.ProgramId}}-{{printf "%.20s" .AuthorName}}-{{printf "%.50s" .ProgramTitle}}'

  # recorderで呼び出されるコマンドのテンプレート
  # リストとして指定してください
  command_template: [
    "ffmpeg",
      "-cookies",  '{{formatCookies .Cookies "\n"}}',
      "-i",        "{{.URL}}",
      "-c",        "copy",
      "-movflags", "faststart",
      "{{.Output}}.{{.Extension}}"
  ]
```

テンプレートは、Goの[`text/template`](https://pkg.go.dev/text/template)の形式に従います。

`output_template`には、[`namagent.Template`](./pkg/namagent/types.go)の値を使用できます。

`command_template`には、[`namagent.StreamData`](./pkg/namagent/types.go)の値を使用できます。

`formatCookies`関数は、[`[]*http.Cookie`](https://pkg.go.dev/net/http#Cookie)の`key=value; domain=<domain>; path=<path>`をセパレータ文字列で結合します。詳細は、[internal/utils/template.go](./internal/utils/template.go)を参照してください。

## LICENSE
[MIT](./LICENSE)

## Author
shinosaki https://shinosaki.com
