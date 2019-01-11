[![CircleCI](https://circleci.com/gh/Financial-Times/relations-api/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/relations-api/tree/master)
[![Coverage Status](https://coveralls.io/repos/github/Financial-Times/relations-api/badge.svg)](https://coveralls.io/github/Financial-Times/relations-api)
# relations-api

Relations Api is an internally used API for retrieving content collection related content.
That is:
- content of CURATED relations
- content of CONTAINS relations for a given content or content collection (content package)

## Usage
### Install
`go get -u github.com/Financial-Times/relations-api`

## Running locally
To run the service locally, you will need to run the following commands first to get the vendored dependencies for this project:
  `dep ensure -vendor-only`

```
Usage: relations-api [OPTIONS]

Options:
  --neo-url="http://localhost:7474/db/data"   neo4j endpoint URL ($NEO_URL)
  --port="8080"                               Port to listen on ($PORT)
  --cache-duration="30s"                      Duration Get requests should be cached for. e.g. 2h45m would set the max-a
ge value to '7440' seconds ($CACHE_DURATION)
```

## Endpoints

### Application specific endpoints:

* /content/{uuid}/relations
* /contentcollection/{uuid}/relations

### Admin specific endpoints:

* /ping
* /build-info
* /__ping
* /__build-info
* /__health
* /__gtg

## Examples

#### For /content/{uuid}/relations endpoint:

`GET https://pre-prod-uk-up.ft.com/__relations-api/content/9b6eb364-0275-11e7-b9ac-52b4e2bf8289/relations`

```
{
       "curatedRelatedContent": [{
           "id": "http://api.ft.com/things/74bd05b4-edca-11e6-abbc-ee7d9c5b3b90",
           "apiUrl": "http://api.ft.com/content/74bd05b4-edca-11e6-abbc-ee7d9c5b3b90"
           }]
        "contains": [{
           "id": "http://api.ft.com/things/74bd05b4-edca-11e6-1234-ee7d9c5b3b90",
           "apiUrl": "http://api.ft.com/content/74bd05b4-edca-11e6-abbc-ee7d9c5b3b90"
           },
           {
           "id": "http://api.ft.com/things/74bd05b4-edca-11e6-1313-ee7d9c5b3b90",
           "apiUrl": "http://api.ft.com/content/74bd05b4-edca-11e6-abbc-ee7d9c5b3b90"
           }]
        "containedIn": [{
           "id": "http://api.ft.com/things/74bd05b4-adsd-1342-abbc-ee7d9c5b3b90",
           "apiUrl": "http://api.ft.com/content/74bd05b4-edca-11e6-abbc-ee7d9c5b3b90"
           }]
   }
```

#### For /contentcollection/{uuid}/relations endpoint (for content package):

`GET https://pre-prod-uk-up.ft.com/__relations-api/content/9b6eb364-0275-11e7-b9ac-52b4e2bf8289/relations`

```
{
        "contains": [{
           "id": "http://api.ft.com/things/74bd05b4-edca-11e6-1234-ee7d9c5b3b90",
           "apiUrl": "http://api.ft.com/content/74bd05b4-edca-11e6-abbc-ee7d9c5b3b90"
           },
           {
           "id": "http://api.ft.com/things/74bd05b4-edca-11e6-1313-ee7d9c5b3b90",
           "apiUrl": "http://api.ft.com/content/74bd05b4-edca-11e6-abbc-ee7d9c5b3b90"
           }]
        "containedIn": [{
           "id": "http://api.ft.com/things/74bd05b4-adsd-1342-abbc-ee7d9c5b3b90",
           "apiUrl": "http://api.ft.com/content/74bd05b4-edca-11e6-abbc-ee7d9c5b3b90"
           }]
   }
```