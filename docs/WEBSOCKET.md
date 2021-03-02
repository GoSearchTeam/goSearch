# Websocket Usage

The websocket endpoint is designed for super fast querying of data, so its feature set it limited.

## Authenticating

First, a client must fetch an authentication token from the origin at the `/ws/auth` route. This will generate a token that must be used later.

## Connecting

The `/ws?token=<API_TOKEN>` endpoint is used, replacing `<API_TOKEN>` with the token from the `/ws/auth` route.

## Querying

Clients can send a stringified JSON message in the same format as the API `/index/search` endpoint to query:

```js
{
  ?query: String, // Optional, text to search for, must be provided if not using `beingsWith`
  ?fields: [String], // Optional, specific fields to search under
  ?beginsWith: String // Optional, must be provided if not using `query`, searches for documents with fields beginning with a string
}
```

This will emit back a stringified JSON object that is a list of documents:

_(parsed)_
```js
{
  [
    {document},
    {document},
    {document},
    {document},
    ...
  ]
}
```
