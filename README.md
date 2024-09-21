# media-nexus

## Table of Contents

* [Overview](#overview)
* [Architecture](#architecture)
* [Build and Run](#build-and-run)
* [Design Choices](#design-choices)

## TODO

* finish tests

## Overview

media-nexus is a simple application to create media and tags through an HTTP API.

The API provides the following functionalities:

* create & list tags
  * a tag is simply a name
* create media
  * media is a tuple (name, list of tag IDs, picture)
* search media by tag IDs

### HTTP API

Run the service (cf. [Build and Run](#build-and-run)) and then navigate to `http://localhost:8081/swagger`.

## Architecture

### Services

```mermaid
graph LR
  subgraph mongodb
    tags[(tags)]
    mmd[(media metadata)]
  end
  subgraph s3
    blobs[(media blobs)]
  end

  a[media-nexus instance] --> tags
  a --> mmd
  a --> blobs
```

### Model

```mermaid
graph LR
  subgraph mongodb
    tags[(tags)]
    mmd[(media metadata)]
  end
  subgraph s3
    blobs[(media blobs)]
  end

  subgraph media item
    metadata --> mmd
    data --> blobs
  end

  subgraph tag
    t[data] --> tags
  end
```

## Build and Run

### Prerequisites

* AWS account configured that's able to manage (create, head, use) the configured media bucket
  * should be setup in `~/.aws/config` and `~/.aws/credentials`
  * or through environment variables (e.g. `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`)
* mongodb instance
* (optional) for `make lint`: `golangci-lint`

### Configuration

The configuration file is located at `local-config.yml`.

### Execution

Build it:

```bash
make deps
make compile
```

And then run it:

```bash
AWS_PROFILE=<aws profile> MEDIANEXUS_MONGODBURI=<mongo uri> ./media-nexus local-config.yml
```

### Documentation

```bash
make docs
```

This will regenerate the documentation. Now relaunch the service and navigate to `http://localhost:8081/swagger`.

### Integration Tests

#### Testing Prerequisites

Similar to [build and run](#build-and-run) above, you need the aws profile and mongdb URI
present in the environment. You can export it manually or do the following:

* edit `local-config.env` with your defaults
* execute `set -a && . ./local-config.env && set +a` to export them to the current shell

#### Testing itself

```bash
make test.integration
```

## Design Choices

### Architectural

A requirement is always AWS, because the company's target prod environment is also
on AWS.

#### Hexagonal Architecture

* very well established architectural pattern for micro services
* separates well inputs, outputs, business logic

#### Media Storage

* requirements:
  * must: persistent, redundant
  * must: store blob & metadata
  * must: metadata queryable
    * contains a tag_ids metadata item, that contains a target tag ID
  * must: easily scalable
  * nice to have: public URLs for retrieving

So we want a blob storage of some sort in the cloud. Because that's persistent,
redundant and easily scalable. For details more below.

### Technology Choices

#### S3 for Blobs

* requirements
  * must: persistent, redundant
  * nice to have: AWS
  * nice to have: multi-regional
  * nice to have: random blob metadata, that's queryable

With all that, S3 is a good choice. S3 doesn't have the metadata requirement, though.
That's why we need an additional storage for that.

#### MongoDB for Media Metadata

* requirements:
  * must: persistent, redundant
  * must: schema-less
    * you don't know the future
    * especially metadata will easily change
  * nice to have: AWS

* Schema-less: all document-based DBs
* persistent, redundant: use AWS or some other cloud
* AWS: you have DocumentDB
  * that's fully mongodb-compatible in case we don't want AWS

#### MongoDB for Tags

* requirements:
  * must: persistent, redundant
  * nice to have: AWS
  * schema-less? unsure

The choice made here is a rather pragmatic one because of the already existing idea
of using MongoDB for media metadata. We could think about using SQL here (possibly
not many schema updates, we always want to query the whole list of tags).

But since we already have MongoDB it's just a way simpler design to use MongoDB here
as well. One database less to care about with all the bells and whistles attached
you would need (costs, monitoring, ...). And most probably cheaper: one additional
collection in an existing DB vs. an additional DB.

### Library Choices

#### Swag for API documentation

* requirements:
  * must: be able to write clear, concise documentation
  * nice to have: written close to the code, not somewhere else
    * docs are always outdated, this reduces this
  * nice to have: nice UI

That's quite obvious I would say. Swagger is sooo wide spread, has so much tooling.
And with just a couple of annotations you get a really nice HTML documentation.

swag is such a tool that checks all the boxes. Actively maintained, 10k stars on github.

**So why not use OpenAPI to generate the boilerplate code for the API?**

It's not that easy. It's very opinionated and you have to work around a lot of this.
Or you buy in fully and then possibly have to rewrite everything when doing a major
version upgrade.

Even then, there are still bugs, because it's a huge project. Then again working around
issues.

#### Logrus for Logging

* requirements:
  * must: log levels
  * must: log format
  * must: structured logging

With such a low-level tool that's spread to really all the code we must expect it will
change and alternatives will come. We choose one and put it behind a facade.

I used Logrus before, it has almost 27k stars on github. Of course there are many
alternatives out there. Logrus started to go into maintenance-only mode. So we could
have chosen another one here to be honest. But so it's good to have it behind a facade.

#### Viper for Configuration

* requirements:
  * must: config defined in code
  * must: define defaults
  * must: env vars
  * must: config file

Viper has almost 27k stars on github, is actively maintained and checks all the boxes.

### Next Steps

* discuss incomplete metadata lifetime
  * that's the time that needs to pass before retrying
  * so should be quite short
  * but so probably different for a picture and a e.g. video
  * right now: task is for pictures, so quite short
* discuss: what should be valid characters for tag & media name?
  * then validate them as well
* deadlines on request contexts
* paging on
  * list of tags
  * list of found media items
* more endpoints
  * query media by name
  * delete media
  * delete tags
  * update media (different name, different tags)
* proper cache headers
  * no cache headers right now, but definitely need that
* metrics
  * depends a bit on the environment
  * inside a service mesh: get many metrics already for free
  * outside: need to implement solution based on opentelemetry probably
* traces
  * need to propagate downstream trace headers to upstream request
