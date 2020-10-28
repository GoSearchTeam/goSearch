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

## /index/add

**Request Type**: POST

**Sent**: JSON document you want to add

**Received**: `Added Index`

## /index/addMultiple
