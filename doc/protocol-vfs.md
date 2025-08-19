# VFS protocol

The VFS (Virtual File System) protocol serves static content from a ZIP archive. It allows you to mount a ZIP file as a file system and access its contents via a URL.

## Configuration

You need to define the archives in your configuration file (e.g., `config.yml`) under the `vfs` protocol section. The configuration is a map of an alias (hostname for the URL) to the path of the ZIP file.

Example:

```yaml
protocol:
  vfs:
    archives:
      mydocs: /path/to/your/documents.zip
      site: /var/www/site.zip
```

## URL Structure

The URL format is `vfs://<alias>/<path-in-zip>`.

*   `<alias>`: The alias name defined in the configuration (e.g., `mydocs`).
*   `<path-in-zip>`: The path to the file or directory within the ZIP archive.

For example, to access `images/photo.jpg` inside `mydocs.zip`, you would use the URL: `vfs://mydocs/images/photo.jpg`.

If the path points to a directory, the protocol will attempt to serve the `index.html` file from that directory.

---

日本語訳 (Japanese translation):

# VFS プロトコル

VFS (Virtual File System) プロトコルは、ZIPアーカイブから静的コンテンツを提供します。ZIPファイルをファイルシステムとしてマウントし、そのコンテンツにURL経由でアクセスすることができます。

## 設定

設定ファイル（例: `config.yml`）の `vfs` プロトコルセクションで、アーカイブを定義する必要があります。設定は、エイリアス（URLのホスト名）からZIPファイルのパスへのマップです。

例:

```yaml
protocol:
  vfs:
    archives:
      mydocs: /path/to/your/documents.zip
      site: /var/www/site.zip
```

## URL構造

URLの形式は `vfs://<alias>/<path-in-zip>` です。

*   `<alias>`: 設定で定義されたエイリアス名（例: `mydocs`）。
*   `<path-in-zip>`: ZIPアーカイブ内のファイルまたはディレクトリへのパス。

例えば、`mydocs.zip` 内の `images/photo.jpg` にアクセスするには、次のURLを使用します: `vfs://mydocs/images/photo.jpg`

パスがディレクトリを指している場合、プロトコルはそのディレクトリ内の `index.html` ファイルを提供しようとします。
