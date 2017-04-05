file2parts
==========

パーツをファイルとして扱えるようになり、使いづらいCMSのUIから開放され、バージョン管理できるようになります。

以下のことができます

* Daisyのデータベースからパーツをファイルとして書き出す
* パーツファイルをデータベースに書き込む
* パーツファイルを監視し、変更されたファイルをデータベースに書き込む

## ダウンロード

| OS | URL |
|---|---|
| windows 64-bit x86 | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/windows-amd64/file2parts.exe |
| windows 32-bit x86 | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/windows-386/file2parts.exe |
| Mac OS 64-bit x86 | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/darwin-amd64/file2parts |
| Mac OS 32-bit x86 | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/darwin-386/file2parts |
| Linux 64-bit x86 | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/linux-amd64/file2parts |
| Linux 32-bit x86 | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/linux-386/file2parts |

## 使い方

* ダウンロードから自分の環境にあった実行ファイルをダウンロードする
* 実行ファイルにPATHを通す、もしくはPATHが通ってるディレクトリに配置する
* とりあえずコマンドを打つとカレントディレクトリに設定ファイルが書き出される
* 設定ファイルを書き換えてDBに接続できるようにする
* あとはお好みで

```
Usage of file2parts:
  -d:  dump parts that exist as files from database
  -da: dump all parts from database
  -r:  restore parts to database
  -w:  watch and restore modified part file
  -rw: alias of "-r -w"
```

* file2parts -d データベース上のパーツのうち、存在するパーツファイルを更新
* file2parts -da データベース上のパーツすべてをファイルとして書き出し
* file2parts -r すべてのパーツファイルをデータベースに書き込み
* file2parts -w ファイルを監視し、更新があったパーツファイルをデータベースに書き込む

基本的に -da で全件取得し、バージョン管理したいものだけを手元に残し、それをバージョン管理します。  
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
に書き換えてください。  
セキュリティ気にする場合はGoogle

## 開発用メモ

TODO: Makefile書きかけ… depとか入れたい

```
go get github.com/go-sql-driver/mysql
go get github.com/go-fsnotify/fsnotify
go get gopkg.in/redis.v5
```
