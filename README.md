# file2parts

パーツをファイルとして扱えるようになり、使いづらいCMSのUIから開放され、バージョン管理できるようになります。

以下のことができます

- Daisyのデータベースからパーツをファイルとして書き出す
- パーツファイルをデータベースに書き込む
- パーツファイルを監視し、変更されたファイルをデータベースに書き込む

## ダウンロード

releaseタブからダウンロードできます

## 使い方

- ダウンロードから自分の環境にあった実行ファイルをダウンロードする
- 実行ファイルにPATHを通す、もしくはPATHが通ってるディレクトリに配置する
- とりあえずコマンドを打つとカレントディレクトリに設定ファイルが書き出される
- 設定ファイルを書き換えてDBに接続できるようにする
- あとはお好みで

```
Usage of file2parts:
  -confing: config file name
  -d:  dump parts that exist as files from database
  -da: dump all parts from database
  -r:  restore parts to database
  -w:  watch and restore modified part file
  -rw: alias of "-r -w"
```

- file2parts -d データベース上のパーツのうち、存在するパーツファイルを更新
- file2parts -da データベース上のパーツすべてをファイルとして書き出し
- file2parts -r すべてのパーツファイルをデータベースに書き込み
- file2parts -w ファイルを監視し、更新があったパーツファイルをデータベースに書き込む

基本的に -da で全件取得し、バージョン管理したいものだけを手元に残し、それをバージョン管理します。<br>
作業中は -rw で手元で編集しつつデータベースをリアルタイムに更新します。

## ありそうな質問

### Q. 設定ファイルはどう書けばいいですか

A. 実行するとファイルがカレントディレクトリに生成されます。

### Q. Redisにつながりません

A. /etc/redis.conf の

```
bind 127.0.0.1
```

を

```
bind 0.0.0.0
```

に書き換えてください。<br>
セキュリティ気にする場合はGoogle

## 開発用メモ

goの開発環境が必要です

```
# リポジトリ取得
# できなさそうな時はこのあたりを参考 http://jnst.hateblo.jp/entry/2016/10/17/210612
go get github.com/team-lab/daisy-file2parts

# 依存解決
dep ensure

# ソースコードのフォーマット整形
make fmt

# ソース実行
make run

# ビルド、成果物生成
make build
```
