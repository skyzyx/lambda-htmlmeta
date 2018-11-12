# HTML Metadata Parsing for AWS Lambda

Will take a URL passed as a query string, and can generate output for OpenGraph, OEmbed, Microdata, and Microformats 2.

This is a _Serverless_ app, written in Go ([Golang]), running in AWS Lambda, with API Gateway in front of it, and AWS CloudFront in front of that (for caching).

**This is an experiment.** Uptime is not guaranteed, and there is no SLA. _But_ all-in-all, it should be reasonably reliable.

> **NOTE:** I've tried to build this in a way that is very cheap to run. But if you find yourself using this more than just occasionally, consider kicking me down a few bucks to offset the cost of running this service. <https://cash.me/$rparman>

## Usage

The `https://metadata.ryanparman.com` hostname is a CloudFront caching distribution, in front of API Gateway, in front of a Lambda function.

There are four endpoints:

* `GET /oembed` — This will return an oEmbed response, formatted as JSON.
* `GET /opengraph` — This will return all OpenGraph-formatted metadata as JSON.
* `GET /microdata` — This will return all Microdata-formatted metadata (e.g., schema.org) as JSON.
* `GET /microformats2` — This will return all Microformats 2-formatted metadata as JSON.

Additionally, it accepts one query-string parameter.

| Parameter | Description |
| --------- | ----------- |
| `url` | (Required) The URL to parse the metadata from. |

## Developing/Deploying

### Golang

Go (when spoken) or [Golang] (when written) is a strongly-typed language from Google that "blends the simplicity of Python with the performance of C". Static binaries can be compiled for all major platforms, and many minor ones.

It is recommended that you install Golang using your system's package manager. If you don't have one (or if the version is too old), you can [install Golang from its website](https://golang.org/doc/install). Reading the [Getting Started](https://golang.org/doc/) documentation is a valuable exercise.

```bash
brew update && brew install golang
```

### Glide

Golang dependencies are managed with [Glide]. You should install them before compiling this project.

```bash
curl https://glide.sh/get | sh
glide install
```

### GoMetaLinter

[GoMetaLinter] pulls together many popular linting tools, and can run them on a project.

```bash
gometalinter.v2 --install
```

### Serverless

[Serverless] is a platform that wraps AWS Lambda and AWS CloudFormation, simplifying the deployment of Lambda apps. Serverless is written in Node.js, so you need to install that first.

I recommend you install the [Node Version Manager][nvm], and use that to install the latest Node.js and npm. Once that's complete, install `serverless`.

```bash
npm i -g serverless
```

### Developing

This app is small, and is self-contained in `main.go`.

_By default_, it expects to be running in AWS Lambda, receiving HTTP requests coming from API Gateway.

If you are performing local development/testing, run `make build` to build for the local platform, then `bin/htmlinfo -mock` to run it. The local build reads from `mock.json` and treats it as an incoming request from API Gateway. You can change the query-string parameters to have the app respond to the documented query string parameters.

```bash
make build

# OpenGraph
META_OPENGRAPH=true bin/htmlinfo -mock

# oEmbed
META_OEMBED=true bin/htmlinfo -mock

# Microdata
META_MICRODATA=true bin/htmlinfo -mock

# Microformats 2
META_MICROFORMATS2=true bin/htmlinfo -mock
```

Make sure that you run the linter to catch any issues.

```bash
make lint
```

### Deployment

`serverless` uses the same [local credentials](https://docs.aws.amazon.com/cli/latest/topic/config-vars.html) that the AWS CLI tools and the AWS SDKs use. If you haven't configured those yet, do that first.

Run `make package` to build a binary for AWS Lambda. Then, `serverless deploy` to deploy the app to your environment.

```bash
make deploy
```

> **IMPORTANT:** Since this deployment uses CloudFront, deployments take _much_ longer than normal deployments. First-time deployments can take around 30 minutes. Subsequent deployments can still take several minutes.

  [Glide]: https://glide.sh
  [Golang]: https://golang.org
  [GoMetaLinter]: https://github.com/alecthomas/gometalinter
  [nvm]: https://github.com/creationix/nvm
  [Serverless]: https://serverless.com/framework/docs/getting-started/
