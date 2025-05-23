# Changelog

## 2025-b4
- Added a --benchmark flag to enable performance benchmarking during file transfers
- Implemented parallel sending of file chunks using Go routines
- The number of parallel workers is dynamically determined based on the user's network characteristics, specifically RTT (Round Trip Time) and available upload bandwidth.

## 2025-b3
- Add support for windows 7
- Created new arch for echo to support diff client creation
- Add hash check for each chunk

## 2025-b2
- Add p2p bidirectional from user A to user B where they both can exchange data between each other 

## 2025-b1
- First prototype of the project

## test
- test
