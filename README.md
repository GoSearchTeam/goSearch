# GoSearch <!-- omit in toc -->
Full-Text Search Engine Written in Go

## Table of Contents <!-- omit in toc -->

- [Performance](#performance)
    - [Data Structures](#data-structures)
    - [Basic Architecture](#basic-architecture)
    - [Ranking and Sorting](#ranking-and-sorting)
    - [Test 1 - A Primitive Test](#test-1---a-primitive-test)

## Performance

#### Data Structures

To build the high performing index and search core, we leveraged _inverted indexes_, _roaring bitmaps_, and _radix trees_.

These allowed the impressive testing performance seen below.

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
_Note: These tests were performed on a MacBook Pro, while in a Zoom call and running lots of other apps. This test is primitive and should be no indication of full performance capabilities. The requests and GoSearch were on the same device (localhost network)._

A Zoom call was opened at ~2.5 million documents stored, which you can see a visible change in dispersion of latencies.
