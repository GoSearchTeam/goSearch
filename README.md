# GoSearch <!-- omit in toc -->
Full-Text Search Engine Written in Go

## Table of Contents <!-- omit in toc -->

- [Performance](#performance)
  - [Data Structures](#data-structures)
    - [Wait, what is a GoodList?](#wait-what-is-a-goodlist)
- [How it Works](#how-it-works)
    - [Basic Architecture](#basic-architecture)
    - [Ranking and Sorting](#ranking-and-sorting)
    - [Test 1 - A Primitive Test](#test-1---a-primitive-test)

## Performance

### Data Structures

To build the high performing index and search core, we leveraged _inverted indexes_, _GoodLists_, and _radix trees_.

These allowed the impressive testing performance seen below.

#### Wait, what is a GoodList?

Glad you asked! I (Dan Goodman) invented it! _See what a did with my name there? 😉 It's also just a good list._ A `GoodList` is basically a `Doubly Linked List` with sorting at modify time. In other words, when ever you insert/update or delete an item, it keeps track on where it is supposed to be moved in the linked list. That way it can modify items and sort them in `O(n)` time, just like normal linked list manipulations would occur in.

This has very practical application for search for multiple reasons:

1. All of the benefits of a doubly linked list (fast traversal, memory optimized during manipulations).
2. When using it as the data type stored in the `Radix Tree`, we can fetch the `X` documents with the highest term frequency by taking the `X` first items in the `GoodList`, with their term frequency, making searches silly fast.
3. We can conserve memory by storing the frequency with the document ID, without having to repeat values or use arrays.
4. We can traverse from the front or end (get documents with highest frequency or lowest frequency)

It's pretty fast too:

```
// Legend: (id, frequency)
Adding 3
-(3, 1)- in 18.503µs
Adding 4
-(3, 1)-(4, 1)- in 1.1µs
Adding 4
-(4, 2)-(3, 1)- in 1.039µs
Adding 5
-(4, 2)-(3, 1)-(5, 1)- in 879ns
Adding 5
-(4, 2)-(5, 2)-(3, 1)- in 916ns
Adding 5
-(5, 3)-(4, 2)-(3, 1)- in 880ns
Adding 4
-(5, 3)-(4, 3)-(3, 1)- in 844ns
Adding 4
-(4, 4)-(5, 3)-(3, 1)- in 849ns
Adding 3
-(4, 4)-(5, 3)-(3, 2)- in 877ns
Adding 3
-(4, 4)-(5, 3)-(3, 3)- in 872ns
Adding 3
-(4, 4)-(3, 4)-(5, 3)- in 871ns
Adding 3
-(3, 5)-(4, 4)-(5, 3)- in 855ns
Adding 1
-(3, 5)-(4, 4)-(5, 3)-(1, 1)- in 846ns
Adding 6
-(3, 5)-(4, 4)-(5, 3)-(1, 1)-(6, 1)- in 883ns
Adding 7
-(3, 5)-(4, 4)-(5, 3)-(1, 1)-(6, 1)-(7, 1)- in 878ns
Adding 6
-(3, 5)-(4, 4)-(5, 3)-(6, 2)-(1, 1)-(7, 1)- in 873ns
Adding 5
-(3, 5)-(4, 4)-(5, 4)-(6, 2)-(1, 1)-(7, 1)- in 839ns
Adding 4
-(3, 5)-(4, 5)-(5, 4)-(6, 2)-(1, 1)-(7, 1)- in 839ns
Adding 4
-(4, 6)-(3, 5)-(5, 4)-(6, 2)-(1, 1)-(7, 1)- in 867ns
```

_See `goodMap.go` for the code, maybe I'll make it into it's own package soon._

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
