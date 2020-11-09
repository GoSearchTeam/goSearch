# GoSearch <!-- omit in toc -->
Full-Text Search Engine Written in Go

## Table of Contents <!-- omit in toc -->

- [Performance](#performance)
  - [Data Structures](#data-structures)
    - [Why OrderedMaps?](#why-orderedmaps)
      - [Advanced NoSQL Patterns for Search](#advanced-nosql-patterns-for-search)
- [How it Works](#how-it-works)
    - [Basic Architecture](#basic-architecture)
    - [Ranking and Sorting](#ranking-and-sorting)
    - [Test 1 - A Primitive Test](#test-1---a-primitive-test)

## Performance

### Data Structures

To build the high performing index and search core, we leveraged _inverted indexes_, _OrderedMaps (called LinkedHashMaps in Java)_, and _radix trees_.

#### Why OrderedMaps?

What we needed was a data structure that had native sorting (for listing first n items in O(n) time), but also could insert, update, or delete items in O(1) time. Enter the sorted map (known as a linked hash map in Java). It was not as simple as using `OrderedMaps`, since we needed to sort by the score (what would traditionally be the `value` in the `key, value` pair maps use). To effectively use this we had to deploy advanced NoSQL data modeling techniques.

##### Advanced NoSQL Patterns for Search

Since we are basically making the NoSQL version of a full-text search indexing platform, in hindsight it's no surprise that NoSQL data modeling techniques make an appearance. Yet until it was realized, it was not obvious.

OrderedMaps are sorted by their key, and as a result we could run into conflicts with documents having the same score if we used them in a `(score, docID)` format. Additionally, we want to design our data structure for read speed, willing to sacrifice insert, update, and delete speeds in the process.

To get the read performance we want (O(n) time for iterating over highest scored n items), we needed the `score` to be in the key. But in order to prevent collisions from overwriting, we also needed the `docID` to be in the key. _Enter compound keys._ By leveraging the format of `(score#docID, null)`, we get the best of both worlds. We can sort by `score`, then by `docID`.

We still maintain very high speeds for insert (O(1) time), as well as update/delete (i + O(1), where i is the time it takes to re-score a document). For delete, since we are given the `docID`, what we can do is fetch the document from disk, re-calculate the score of each field, then use those scores to make O(1) delete operations on the OrderedMap by using the `score#docID` key. For updates, we are also given the `docID`, and perform a delete than insert (i + 2(O(1)) time). This keeps all operations very fast, and by using a `null` (`nil` in Go) value we save some memory since we can pull both the `score` and `docID` by splitting the key at the `#`.

## How it Works

#### Basic Architecture

`GoSearch` has the follow basic hierarchy:

`[]App` -> `[]AppIndex` -> `Index Tree` -> `List Item`

Multiple Apps can exit within a node/cluster. An `App` is just a logical separation of `AppIndexes` to separate what data is searched. The goal is that you can connect multiple products or features to the same cluster.

`AppIndexes` contain the indexes for the `App`, along with the field name. An `App` will consume JSON data, and separate each field into an `AppIndex`. This allows us to perform searches on specific fields of a JSON object if requested (e.g. only search the `title` and `author` fields but not the `postContent` of a blog post). This allows for faster search times when less input is required to perform a search.

`AppIndexes` have `Index Trees`, which are `Radix Trees`. [Radix Trees](https://en.wikipedia.org/wiki/Radix_tree) allow for really fast prefix search. It is what enables us to have the `beginsWith()` search method, allowing search on the prefix of word.

Each node in an `Index Tree` is a `List Item`. A `List Item` is a [Roaring Bitmap](https://roaringbitmap.org/) which is a high speed compressed bitmap. These bitmaps are of `uint64` giving us 1.844E19 possible documents to store on disk. It also allows us to efficiently perform array operations with a smaller footprint. These bitmaps represent the names of the documents on disk.

During this process we handle [Ranking and Sorting](#ranking-and-sorting).

#### Ranking and Sorting

Opening up a document to rank and sort is very expensive. In order to handle the ranking and sorting, we look at the frequency of a document as it appears in a search. The more times a document appears from the search of the `AppIndex`, the higher rank it obtains. We then take this stage 1 rank, and begin to open documents. Once the documents are open, we then take into account what information needs to be served back. If only certain fields of a document are requested, we filter out the rest of the fields before sending the data back.

#### Test 1 - A Primitive Test
![unknown-3](/assets/unknown-3.png)
_Note: These tests were performed on a MacBook Pro, while in a Zoom call and running lots of other apps. This test is primitive and should be no indication of full performance capabilities. The requests and GoSearch were on the same device (localhost network). This test is also pre-GoodLists._

A Zoom call was opened at ~2.5 million documents stored, which you can see a visible change in dispersion of latencies.
