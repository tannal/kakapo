# LSM Memtable & SST Based Key-Value Store

## Overview

This project implements a simple yet efficient key-value store based on the Log-Structured Merge-tree (LSM) architecture. It utilizes in-memory tables (memtables) and disk-based Sorted String Tables (SSTs) to provide fast write and read operations. This store is ideal for applications requiring high write throughput and efficient read operations.

## Features

- **LSM Tree Architecture**: Leverages the LSM tree design for optimized write performance.
- **Memtable**: In-memory data structure for fast write operations.
- **SST Implementation**: Persistent storage of data in sorted string tables on disk.
- **Compaction**: Periodic merging and compaction of SST files to optimize read performance and storage efficiency.
- **Concurrency Support**: Basic threading model to handle concurrent read/write operations.
- **Crash Recovery**: Basic mechanism to recover from unexpected shutdowns.
