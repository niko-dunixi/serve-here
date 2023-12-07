# serve-here 📄

Legit doesn't get any simpler than this!

A straight forward static file server that can be used to serve files
over HTTP during developement. Useful for one-off and quick file-serving.
> This is not indented to replace a Apache, NGINX, or any proper CDN.
> Use proper production best practices outside of local development and
> iteration.

### Installation
```bash
$ go install github.com/niko-dunixi/serve-here
```

### Usage
```bash
$ cd /path/to/my/static/directory
$ serve-here // Your browser will open to this directory
```
