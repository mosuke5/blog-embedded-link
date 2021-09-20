# Blog Embedded Link Generator (BELG)
[![Unit testing](https://github.com/mosuke5/blog-embedded-link/actions/workflows/test.yaml/badge.svg)](https://github.com/mosuke5/blog-embedded-link/actions/workflows/test.yaml)

## 背景
マークダウンでブログやドキュメントを書くことがおおい。
埋め込みリンクをやりたいが、良いツールがない。

iframely  
一番使い勝手がいいが、iframeのため重い。Static Web向きではない。  
google analyticsのevent埋め込みなどができず、クリック数を集計できないなど問題あり。

## ツール概要
URLを与えると、そのサイトの情報を取得し、埋め込みリンクを生成する。
Static Webでの速さを重視し、プレーンなHTMLとして出力する。

必要な情報

- ページタイトル
- サイトの概要
- イメージ画像（あれば）
- サイトタイトル
- サイトのファビコン

こんなかんじ  
![image](embeded-link-image.png)

## ツール仕様
URLを引数に渡して実行するとHTMLの出力をする

```
$ belg https://xxxxxxxxx/aaaa/bbbb
<div class="belg-link">
  <div class="belg-left">
    <img src="https://opengraph.githubassets.com/0499c14d87df16fe94ed2a14fc292954e6ee3df56759374ffd5fb6626a6d59d9/ndabAP/vue-go-example" />
  </div>
  <div class="belg-right">
    <div class="belg-title">GitHub - ndabAP/vue-go-example: Vue.js and Go example project</div>
    <div class="belg-description">Vue.js and Go example project. Contribute to ndabAP/vue-go-example development by creating an account on GitHub.</div>
    <div class="belg-site-name">GitHub</div>
  </div>
</div>
```
