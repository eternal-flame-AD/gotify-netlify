# gotify-netlify

An example gotify plugin for receiveing webhooks from netlify.

## Development

1. Clone this repo
1. `export GO111MODULE=on` if you are in GOPATH
1. `make download-tools`
1. Make changes
1. `make build` to build the plugin for the master branch of gotify

## Use the release

If you found your gotify version is included in the build, you can download the shared object and put that into your plugin dir.

If you did not find you gotify version, follow these steps to build one for your own:

1. Download a zip file of the source code of current release at the releases page and extract it.
1. `export GO111MODULE=on` if you are in GOPATH
1. `make download-tools`
1. `make GOTIFY_VERSION=v1.2.1 build` to build the plugin for your gotify version (`GOTIFY_VERSION` could be a tag, a branch or a commit).