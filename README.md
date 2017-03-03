file2parts
==========

パーツをファイルとして扱えるようになり、使いづらいCMSのUIから開放され、バージョン管理できるようになります。

以下のことができます

* Daisyのデータベースからパーツをファイルとして書き出す
* パーツファイルをデータベースに書き込む
* パーツファイルを監視し、変更されたファイルをデータベースに書き込む

## ダウンロード

| OS | URL |
|:--:|:---:|
| windows | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/windows-amd64/file2parts.exe |
| Mac OS | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/darwin-amd64/file2parts |
| Linux | http://raw.github.team-lab.local/twakayama/daisy-file2parts/master/build/linux-amd64/file2parts |

※すべて amd64 向けです

## 使い方

* バイナリにパスを通す
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

## メモ

TODO: Makefile書きかけ… depとか入れたい

```
go get github.com/go-sql-driver/mysql
go get github.com/go-fsnotify/fsnotify
```
