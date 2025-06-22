TODO: 翻訳する
TODO: TLS利用の手順を詳細化する

# Secure Contexts

NVGDのいくつかの機能(OPFS等)は[保護されたコンテキスト][mdnsc]を必要とします。
保護されたコンテキストとは、ローカルホストへの接続、TLS(HTTPS)接続、もしくはブラウザのフラグ指定で許可されたオリジンへの接続です。
NVGDはその性質上、閉じたLAN内で運用されることを想定しています。
そのためTLS用の正規の証明書を利用することが困難なケースが多く、
その場合にはオレオレ証明書を用いたTLSを使うことになります。

本文章ではブラウザのフラグ指定で指定したオリジンを強制的に保護されたことにする方法と、NVGDにオレオレ証明書でTLSを提供させる方法を解説します。

[mdnsc]:https://developer.mozilla.org/ja/docs/Web/Security/Secure_Contexts

## ブラウザのフラグ指定でオリジンを強制的に保護されていることにする

1.  ブラウザで `chrome://flags/#unsafely-treat-insecure-origin-as-secure` を開く

    この設定はChrome用

    *   Edgeの場合は `edge://flags/#unsafely-treat-insecure-origin-as-secure`
    *   Firefoxには同等の設定は2025-06-22時点で存在しない

2.  入力エリアに保護されていることにするURLを記載する

    *   URLはスキーマ及びポート番号も記載する
    *   複数のURLを指定する場合はカンマ記号 `,` で区切る

    記述例:

    ```
    http://192.168.0.100:9280,http://192.168.0.101:9280,http://dev.mydomain.org:9280
    ```

3. 機能を有効化する(要再起動)

## オレオレ証明書を作成しインストールする

1.  オレオレ証明書を作成しインストールする
2.  NVGDをTLSモードで起動する
3.  ブラウザにオレオレ認証局を追加する
