# goSearch <!-- omit in toc -->
Full-Text Search Engine Written in Go

## Table of Contents <!-- omit in toc -->

- [Performance](#performance)
    - [Data Structures](#data-structures)
    - [Test 1](#test-1)

## Performance

#### Data Structures

To build the high performing index and search core, we leveraged _inverted indexes_, _roaring bitmaps_, and _radix trees_.

These allowed the impressive testing performance seen below.

#### Test 1
![Messages Image(571780674)](/assets/Messages%20Image(571780674).png)
_Note: These tests were performed on a MacBook Pro, while in a Zoom call and running lots of other apps. This test is primitive and should be no indication of full performance capabilities._
