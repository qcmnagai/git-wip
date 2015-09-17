# git-wip
[![Circle CI](https://circleci.com/gh/qcmnagai/git-wip/tree/master.svg?style=shield)](https://circleci.com/gh/qcmnagai/git-wip/tree/master.svg?style=shield)

## Description

github上でWIP (Work In Progress) 用のPull-Requestを簡単に作成する為のGitのサブコマンドです。
以下のことを自動で実行してくれます。

- origin/masterから新しくブランチを作成してリモートにプッシュ（Optional）
- 下記フォーマットのPull-Requestを作成する
    - タイトル: ```wip <Issue Title> (closes #<Issue Number>) ```
    - コメント: ```ref. #<Issue Number>```

## Usage

既にPull-Request用のブランチがリモート＆ローカル上に存在している場合
```
git wip -i [issue number] -b [branch name]
```

新しくPull-Request用のブランチを作成する場合
```
git wip -i [issue number] -c [new branch name]
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/qcmnagai/git-wip
```

## Contribution

1. Fork ([https://github.com/qcmnagai/git-wip/fork](https://github.com/qcmnagai/git-wip/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[qcmnagai](https://github.com/qcmnagai)
