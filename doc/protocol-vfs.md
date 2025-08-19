# VFS protocol

The VFS (Virtual File System) protocol serves static content from a ZIP archive.
It allows you to mount a ZIP file as a file system and access its contents via a URL.

## Configuration

You need to define the archives in your configuration file (e.g., `config.yml`) under the `vfs` protocol section.
The configuration is a map of an alias (hostname for the URL) to the path of the ZIP file.

Example:

```yaml
protocol:
  vfs:
    archives:
      mydocs: /path/to/your/documents.zip
```

## URL Structure

The URL format is `vfs://<alias>/<path-in-zip>`.

*   `<alias>`: The alias name defined in the configuration (e.g., `mydocs`).
*   `<path-in-zip>`: The path to the file or directory within the ZIP archive.

For example, to access `images/photo.jpg` inside `mydocs.zip`, you would use the URL: `vfs://mydocs/images/photo.jpg`.
This URL can be accessed as nvgd at the following URL:

    http://127.0.0.1:9280/vfs://mydocs/images/photo.jpg

Also, the `vfs://` scheme has an alias `vfs/`, so the following URL has the same meaning.
This is a workaround for the problem that in web applications, consecutive slashes in a URL path are sometimes normalized and combined into one.

    http://127.0.0.1:9280/vfs/mydocs/images/photo.jpg

If the path points to a directory, the protocol will attempt to serve the `index.html` file from that directory.

## Examples

### JupyterLite

Download [`jupyterlite-playground-0.0.1.zip`][jl_zip] from [here][jl_release] and configure it as follows to use JupyterLite (with OPFS support) at `http://127.0.0.1:9280/vfs/jupyterlite/`.
The alias name `jupyterlite` can be changed freely.

```yaml
protocol:
  vfs:
    archives:
      jupyterlite: /path/to/your/jupyterlite-playground-0.0.1.zip
```

[jl_release]:https://github.com/koron/jupyterlite-playground/releases/tag/v0.0.1
[jl_zip]:https://github.com/koron/jupyterlite-playground/releases/download/v0.0.1/jupyterlite-playground-0.0.1.zip

## Remarks

*   Filters cannot be applied to content provided by the vfs protocol.
    This is to prevent misinterpretation of query parameters that may be specified in web applications.
