# Database Configuration

This document discusses database configuration and management. It also discusses 
how data is handled and stored. Do remember that SIEMSeasy does not have it's own 
customised indexer. Instead it relies on third-party databases and tooling.

### ClickHouse

This was the biggest architecture decision that I had to make but also the one
that I could make the biggest difference. Most traditional SIEM utilities opt for 
ELK stack because of its powerful search functionality but it requires excessive 
system resources (>16GB) just to run due to heavy Java heap allocations. 

The core promise of SIEMSeasy is that it offers a much faster and lightweight
experience. Hence, I picked ClickHouse. ClickHouse is a SQL database but it stores
data column-wise rather than row-wise. This allows it to process and run aggregate
functions much more quickly on ClickHouse compared to Postgres or MySQL.

ClickHouse has a lot of quality-of-life features for string processing and text
searching. It supports RegEx matching built-in while also offering other nice-to-haves
such as bloom filters and more. Despite all this though, ClickHouse still isn't as
robust as ElasticSearch when it comes to text searching but it offsets it with
lower size on disk and faster query times on average.

I'm planning to fine-tine and configure ClickHouse from the application in such a
way where it will be possible to match the speed of ElasticSearch or at least reach
very close to it for text processing which is important for logs. 