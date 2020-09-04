# goSearch <!-- omit in toc -->
Full-Text Search Engine Written in Go

## Table of Contents <!-- omit in toc -->

- [Performance](#performance)
    - [Data Structures](#data-structures)

## Performance

#### Data Structures

After extensively testing golang's built-in linked lists library, `container/list`, raw linked lists (smaller feature set), skip lists, btrees, and b+trees, we display some impressive results:

![speed test](/assets/unknown-3.png)
_Fig. 1 - built-in linked lists vs. raw linked lists vs. skip list vs. btree vs. b+tree operations (100 million records)_

As seen in _Fig. 1_, the speed advantages of b+trees for our use case are enormous. With `100 Million Records`, we get a random search time of only ~6 microseconds using a b+tree.

This is the major advantage that b+trees present, as a significantly large majority of operations will be search, as opposed to data modification operations.

There is a clear slowdown when it comes to modifying the data in the tree, but by sacrificing time on the speed of modification, we by far make that time up in search performance.

With this performance, we can recognize for about `170 thousand` searches per second.

We run an _unlocked_ tree, meaning that CRUD operations can happen async. (in parallel). This increases speed while presenting no downsides for our use case.
