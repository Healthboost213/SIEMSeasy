# Search & Parsing

This document goes through the ins and outs of SIEMSeasy's search and parsing
functionality. It should serve as basic documentation and high-level architectural
decision making rather than a line-by-line code breakdown.

### File Tailing

File tailing was an interesting challenge to tackle particularly. It is important
to be memory efficient when reading through files and only read what we need.
Therefore, we can set-up a simple buffer to hold the file data. We can use the
buffer to scan through for event boundaries and seperate the chunks of line data.

I thought of originally just scanning for `CRLF` or `LF` but that didn't really pan
out as expected. The original bit of code I had written was working very
inconsistently and wouldn't have been effective for logs that span multiple lines
(e.g. Java or Python stack traces).

Therefore, I decided to use RegEx. I configured it to do multi-line scanning by
default. It finds matches all across the buffer and returns their indexes. With a
little bit of byte slicing, I can narrow down the last line and easily extract it.
If the size of one event somehow exceeds the size of the buffer, it will keep
adjusting the offsets and reading backwards until a match is found.

RegEx may not be the fastest method to extract event data but it is not terribly
slow either. Even if we are tailing 50-100 log files on just one endpoint, it
wouldn't really be the main bottleneck for performance. It might also seem
counterintuitive against 