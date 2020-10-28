# API Endpoints

## /index/listItems

**Request Type**: GET

**Received**:  Array of objects, containing index with their associated values

**Example**:

```JSON

[
    {
        "IndexName": "index1",
        "IndexValues": [
            "value1",
            "value2"
        ]
    },
    {
        "IndexName": "index2",
        "IndexValues": [
            "valueA",
            "valueB",
            "value1"
        ]
    }
]
```

## /index/listIndexes

**Request Type**: GET

**Received**:  Array of unique indexes across all documents

## /index/search
**Request Type**: POST

**Sent**: JSON object with the following attributes

* query - search query
* fields - array of indices to be searched
* beginsWith (optional) - if true matches values beginning with query, if false matches values exactly with query (defaults to false)
**Received**: JSON object with a docIDs field (array of document ID's) and a document field (array of documents which match the search)

## /index/add

**Request Type**: POST

**Sent**: JSON document you want to add

**Received**: `Added Index`

## /index/addMultiple
**Request Type**: POST

**Sent**: JSON object with items field containing documents to add
```JSON
{
    "items": [
        {
            "index1": "value1",
            "index2": "value2"
        },
        {
            "index1": "valueA",
            "index2": "valueB"
        }
    ]
}
```
**Received**: Added Indexes
