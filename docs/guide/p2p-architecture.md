# P2P Architecture

Echo's peer-to-peer architecture is the foundation of its performance and privacy advantages. This page explains how Echo's P2P system works and why it matters for file transfers.

## What is Peer-to-Peer?

In a traditional client-server model, all data must flow through central servers:

```
Sender → Central Server → Recipient
```

This creates bottlenecks, single points of failure, and privacy concerns.

With Echo's peer-to-peer model, data flows directly between devices:

```
Sender ↔ Recipient
```

This direct connection offers several advantages:
- Higher transfer speeds
- Enhanced privacy
- Reduced dependency on external infrastructure
- Works in offline environments

## Echo's Network workflow

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
                                        | WriteToUDP(senderAddr) <----------
<------------------------------------- Wait for ACK (handleAck)
If last chunk:
- Validate checksum
```