file2parts
==========

以下のことができます

* Daisyのデータベースからパーツをファイルとして書き出す
* ファイルのパーツをデータベースに書き込む
* ファイルを監視し、変更されたファイルをデータベースに書き込む

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

## メモ

TODO: Makefile書きかけ… depとか入れたい

```
go get github.com/go-sql-driver/mysql
go get github.com/go-fsnotify/fsnotify
```
