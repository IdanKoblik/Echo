# File Transfer Protocol

Echo uses a sophisticated protocol called `FileChunk` for reliable UDP-based file transfers. The protocol uses Protocol Buffers (proto3) for efficient serialization.

### Protocol Structure

```
LSB                                                                MSB
Byte:   1       2       3       4       5       6       7       8
     +-------+-------+-------+-------+-------+-------+-------+-------+
     |                        version (uint32)                       |
     +-------+-------+-------+-------+-------+-------+-------+-------+
Byte:   9      10      11      12      13      14      15      16
     |                      chunkIndex (uint32)                      |
     +-------+-------+-------+-------+-------+-------+-------+-------+
Byte:  17      18      19      20      21
     |           totalChunks (uint32)        | isLastChunk (1 byte)  |
     +-------+-------+-------+-------+-------+-------+-------+-------+
Byte:  22+ (variable)
     | filename length (varint) |
     +--------------------------+
     | filename (N bytes)       |
     +--------------------------+
     | data length (varint)     |
     +--------------------------+
     | data (M bytes)           |
     +--------------------------+
     | checksum length (varint) |   ← only if `isLastChunk == true`
     +--------------------------+
     | checksum (SHA256 string) |
     +--------------------------+
```

### Protocol Operation

#### Sender Process
1. Reads the file and splits it into 1024-byte chunks
2. Wraps each chunk in a FileChunk message
3. Sends chunks over UDP and waits for acknowledgment
4. Includes SHA-256 checksum with the final chunk

#### Receiver Process
1. Receives each FileChunk message
2. Writes chunk data to disk
3. Sends acknowledgment back to sender
4. Verifies complete file using checksum after receiving final chunk

### Protocol Features
- Chunked UDP transfer for efficient data transmission
- Reliable delivery through acknowledgment system
- File integrity verification via SHA-256 checksums
- Flexible filename and size handling
- Efficient serialization using Protocol Buffers
