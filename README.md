# GoSearch <!-- omit in toc -->
Full-Text Search Engine Written in Go

## Table of Contents <!-- omit in toc -->

- [Features](#features)
- [Usage](#usage)
  - [Node Configuration](#node-configuration)
  - [Web interface](#web-interface)
    - [CLI Arguments](#cli-arguments)
    - [To-Do:](#to-do)
- [Performance](#performance)
  - [Data Structures](#data-structures)
    - [Why OrderedMaps?](#why-orderedmaps)
      - [Advanced NoSQL Data Patterns for Search](#advanced-nosql-data-patterns-for-search)
- [How it Works](#how-it-works)
    - [Basic Architecture](#basic-architecture)
    - [Ranking and Sorting](#ranking-and-sorting)
    - [Why NoSQL Search?](#why-nosql-search)
    - [Clustering](#clustering)
    - [Test 1 - A Primitive Test](#test-1---a-primitive-test)
    - [Test 2 - A Single Node Comparative Test](#test-2---a-single-node-comparative-test)

## Features

- NoSQL design
- Web socket interface for direct user connection giving extremely low latency search
- Preemptive scoring, sorting, and filtering
- Fault tolerance (replication)
- Global and local clusters (similar to Cassandra)
- Linearly scalable throughput

## Usage

### Node Configuration

In its current form, all configuration is done through the cli.

**Example usage** (node joining cluster): ``./goSearch --cluster-mode --iface="192.168.86.237" --gossip-port=7777 --port=8182 --fellow-nodes="192.168.86.237:4444" --local-cluster="lc1" --global-cluster="glob1"``

### Web interface

Nodes will also each present a web ui from the `/admin` url. In this interface, you can do basic document operations, monitor the performance of the node, and see its logs.

#### CLI Arguments

`--cluster-mode`: `boolean`
If provided, the node knows to either begin a cluster, or join one.

`--iface`: `string`
The interface on which to provide gossip communication. Should be kept over a private network.

`--gossip-port`: `int`
The port on which to provide gossip communication.

`--port`: `int`
The port on which the API will listen to.
_The API currently listens to all interfaces for development purposes, this will become an independent cli flag_.

`--fellow-nodes`: `[]string`
A CSV of other nodes in the form of `gossip_ip:gossip_port`. Only one node needs to be provided to join the cluster. This flag is only required when a node is joining an existing cluster.

`--local-cluster`: `string`
The name of the local cluster that the node is joining (reginal)

`--global-cluster`: `string`
The name of the global cluster that the node is joining. This must match the other nodes listed in the `fellow-nodes` flag.

#### To-Do:

- Add pre-shared key for connecting to cluster
- Add in TLS support for gossip
- Move cluster sharing `addIndex` off HTTP interface
- Add concurrency/goroutine to index operation sharing
- Add other index updates to cluster sharing
- Add other node events to gossip
- Add cluster node heartbeat and dead/suspect node detection and logic
- many more things to make this production ready!
- Build a custom library with latency detection for each node (like datastax cassandra driver)

## Performance

### Data Structures

To build the high performing index and search core, we leveraged _inverted indexes_, _OrderedMaps (commonly known as LinkedHashMaps in Java)_, and _radix trees_.

#### Why OrderedMaps?

What we needed was a data structure that had native sorting (for listing first n items in O(n) time), but also could insert, update, or delete items in O(1) time. Enter the sorted map (known as a linked hash map in Java). It was not as simple as using `OrderedMaps`, since we needed to sort by the score (what would traditionally be the `value` in the `key, value` pair maps use). To effectively use this we had to deploy advanced NoSQL data modeling techniques.

##### Advanced NoSQL Data Patterns for Search

Since we are basically making the NoSQL version of a full-text search indexing platform, in hindsight it's no surprise that NoSQL data modeling techniques make an appearance. Yet until it was realized, it was not obvious.

OrderedMaps are sorted by their key, and as a result we could run into conflicts with documents having the same score if we used them in a `(score, docID)` format. Additionally, we want to design our data structure for read speed, willing to sacrifice insert, update, and delete speeds in the process.

To get the read performance we want (O(n) time for iterating over highest scored n items), we needed the `score` to be in the key. But in order to prevent collisions from overwriting, we also needed the `docID` to be in the key. _Enter compound keys._ By leveraging the format of `(score#docID, len(field))`, we get the best of both worlds. We can sort by `score`, then by `docID`.

We still maintain very high speeds for insert (O(1) time), as well as update/delete (i + O(1), where `i` is the time it takes to re-score a document). For delete, since we are given the `docID`, what we can do is fetch the document from disk, re-calculate the score of each field, then use those scores to make O(1) delete operations on the OrderedMap by using the `score#docID` key. For updates, we are also given the `docID`, and perform a delete than insert (i + 2(O(1)) time). This keeps all operations very fast, and by using a `len(field)` value we save some memory since we can pull both the `score` and `docID` by splitting the key at the `#`, and keep the length of that JSON field (key) stored in memory to do all scoring without loading the document object. At the end, we just sum the term scores of the same document (among other things) to get the final document score, and sort based on that score.

## How it Works

#### Basic Architecture

`GoSearch` has the follow basic hierarchy:

`[]App` -> `[]AppIndex` -> `Index Tree` -> `List Item`

Multiple Apps can exit within a node/cluster. An `App` is just a logical separation of `AppIndexes` to separate what data is searched. The goal is that you can connect multiple products or features to the same cluster.

`AppIndexes` contain the indexes for the `App`, along with the field name. An `App` will consume JSON data, and separate each field into an `AppIndex`. This allows us to perform searches on specific fields of a JSON object if requested (e.g. only search the `title` and `author` fields but not the `postContent` of a blog post). This allows for faster search times when less input is required to perform a search.

`AppIndexes` have `Index Trees`, which are `Radix Trees`. [Radix Trees](https://en.wikipedia.org/wiki/Radix_tree) allow for really fast prefix search. It is what enables us to have the `beginsWith()` search method, allowing search on the prefix of word.

Each node in an `Index Tree` is a `List Item`. A `List Item` is an `OrderedMap`. `OrderedMaps` are a combination of `Doubly Linked Lists` and `maps`, providing update operations in O(1) time, and iteration in O(n). Combined with the pre-sorting and pre-ranking this allows us to search very quickly.

During this process we handle [Ranking and Sorting](#ranking-and-sorting).

#### Ranking and Sorting

The way in which the data is stored and sorted is a proprietary modification of the `Pivoted Normalization Formula`. What gives GoSearch such speed and consistency is that **part of the algorithm for scoring and sorting a search result is done at index time**, meaning we have around half of the formula and sorting completed before a search result even comes in. **At search time we only have to perform a subset of the typical operations on a much smaller dataset as the documents are pre-sorted and partially pre-scored.**

_The above section is intentionally kept simple. See more detail in [Advanced NoSQL Data Patterns for Search](#advanced-nosql-data-patterns-for-search)._

#### Why NoSQL Search?

GoSearch tackles the same problems for search that DBs like Cassandra and DynamoDB tackle for databases. Low latency, eventually consistent, linearly scalable, and global distribution. GoSearch also adds additional features such as direct user connection through web sockets for extremely low latency search.

NoSQL is a double edged sword. One one hand, you have extremely low level access to how the data is handled, meaning you can manipulate it more flexibly and store is less structured. On the other hand, you need to put in the work up front for designing these data models and access patterns in such a way that you can still do direct lookups instead of performing scans over the data.

This low level access to the data (vs. something like SQL) allows us to preemptively sort and score documents before a search occurs, as explained above.

#### Clustering

GoSearch uses clustering to provide linear scalability for throughput. In it's current form, more RAM and disk will need to be added to increase the amount of documents stored on a node.

GoSearch uses a custom Gossip implementation on top of TCP to handle inter-node metadata communication. When a single node joins the cluster using a command like `./goSearch --cluster-mode --iface="192.168.86.237" --gossip-port=7777 --port=8182 --fellow-nodes="192.168.86.237:4444" --local-cluster="lc1" --global-cluster="glob1"`, it only needs one `fellow-nodes` in the list to find all nodes total in the cluster within microseconds (on the same network, cloud results may take up to a few milliseconds). It uses a default TTL of 6 for gossip messages, which is more than enough for 120+ nodes in a cluster.

When a node receives an index operation (add, update, delete), it first performs it locally. If that is successful, then it tells every other node in the cluster to perform the same operation, reaching consensus in the low single digit milliseconds. This ensures that all nodes are eventually consistent with the dataset. Since writes are very low compared to reads for full-text search, this should not be a concern.

Since all nodes have the same dataset, searches are performed entirely locally, resulting in maximum performance. It also means that any node can be searched from, allowing you to always contact the lowest latency node for the fastest possible search result.

The result of the combination of this clustering and data handling model is extremely low latency and high consistency search results across all nodes.

#### Test 1 - A Primitive Test
![unknown-3](/assets/unknown-3.png)
_Note: These tests were performed on a MacBook Pro, while in a Zoom call and running lots of other apps. This test is primitive and should be no indication of full performance capabilities. The requests and GoSearch were on the same device (localhost network). This test is also pre-GoodLists._

A Zoom call was opened at ~2.5 million documents stored, which you can see a visible change in dispersion of latencies.

This test was run before implemented NoSQL techniques and `OrderedMaps`

#### Test 2 - A Single Node Comparative Test
![GoSearch Testing2-es](/assets/GoSearch%20Testing2-es.png)

This test was run on a single node cloud environment. Quickly summarizing some of the test results:

- 100k documents, 1k searches
- Average Document Add Request Time: `0.782ms`
- Average Document Search Request Time: `0.822ms`
- Over 2x performance boost over ElasticSearch running the same test on the same hardware
- Consistently kept HTTP requests under `1ms` in same AZ and VPC
- Request time was far less volatile than ElasticSearch, and average request time actually decreased over time.
- No repeated documents or searches were performed in this test
- Requests were initiated immediately after the previous request responded

This test shows seriously promising results. We can see the advantage of the NoSQL techniques and preemptive scoring and sorting. By doing a lot of the work during the add operation (which we are even faster in!), we save lots of time during the search operations.

Obviously 100k documents is a relatively small amount, and a single node is not practical for production. In further tests we will see how a clustered environment performs.
