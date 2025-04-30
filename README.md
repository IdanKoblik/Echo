<h1 align="center" style="display: flex; align-items: center; justify-content: center;">
    Echo
</h1>

<div align="center">
  <a href="https://coveralls.io/github/IdanKoblik/Echo?branch=main">
    <img src="https://coveralls.io/repos/github/IdanKoblik/Echo/badge.svg?branch=main" alt="Coverage Status">
  </a>

  <img src="https://img.shields.io/github/actions/workflow/status/IdanKoblik/Echo/main.yml" alt="GitHub Actions Workflow Status">

  <a href="https://sonarcloud.io/summary/new_code?id=IdanKoblik_Echo">
    <img src="https://sonarcloud.io/api/project_badges/measure?project=IdanKoblik_Echo&metric=alert_status" alt="Quality Gate Status">
  </a>

  <a href="https://github.com/IdanKoblik/echo/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/IdanKoblik/echo.svg" alt="License">
  </a>
</div>
<br>

> ⚠️ **Development Status**: This project is currently in development stage.

<br>
Echo is a lightweight, peer-to-peer (P2P) file transfer system designed for reliability, speed, and simplicity. Whether you're sending files across your local network or over the internet, Echo ensures smooth and secure transfers without relying on central servers.

### 🌟 Key Benefits
* P2P Architecture – Transfers go directly between devices, no middleman required.
* User-Friendly – Designed with simplicity in mind, making it easy to send and receive files.
* Reliable Over UDP – Ensures reliable delivery using acknowledgments and chunk-based transfer.

### 📡 File Transfer Protocol — `FileChunk`
> This protocol enables reliable **UDP-based file transfer** using **Protocol Buffers (proto3)** for serialization.
<br>

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

**Workflow**
```
User A (Sender)                          User B (Receiver)
-----------------                       -------------------
Select "Send a file"                    Select "Receive a file"
Input local port, remote addr           Input local port, remote addr
|
Open file and split into chunks
|
for each chunk: ----------------------> ReadFromUDP (waits for chunk)
                                        |
Marshal FileChunk (with data)           |
WriteToUDP(remoteAddr) ---------------->
                                        |
                                        | Unmarshal chunk
                                        | Write chunk to disk
                                        | Marshal ACK
                                        | WriteToUDP(senderAddr) <----------------
<------------------------------------- Wait for ACK (handleAck)
If last chunk:
- Validate checksum
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