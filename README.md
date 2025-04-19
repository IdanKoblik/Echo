# DataDrop

### 📡 File Transfer Protocol — `FileChunk`
> This protocol enables reliable **UDP-based file transfer** using **Protocol Buffers (proto3)** for serialization.

📐 **Structure**

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

📦 **Protocol Overview**
1) **Sender**:
    * Reads file, splits it into 1024-byte chunks.
    * Wraps each chunk in a FileChunk message.
    * Sends over UDP and waits for ACK before continuing.

2) **Receiver**:
    * Receives each FileChunk message.
    * Writes data to disk and sends back a FileAck.
    * If it's the last chunk, verifies the file with the checksum.

🧪 **Features**
* ✅ Chunked UDP transfer
* ✅ Reliable delivery via ACKs
* ✅ Checksum validation (SHA-256)
* ✅ Filename and size flexibility
* ✅ Protobuf-powered serialization