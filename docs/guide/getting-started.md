# Getting Started

Echo is a peer-to-peer file transfer tool over UDP, designed for simplicity and speed. You can use it interactively or with command-line flags.

## 1. Download Echo

Visit the [GitHub Releases](https://github.com/IdanKoblik/Echo/releases) and download the binary for your system.

Make it executable if needed:

```bash
chmod +x echo
```

## 2. Run Echo (Interactive Mode)

If you run Echo with no arguments, it will guide you through the setup:

```bash
./echo
```

You’ll be prompted to:

- Select **send** or **receive**
- Enter your **local port**
- Enter the **peer’s address**
- (If sending) Provide a **file path**

## 3. Run Echo (CLI Mode)

To use Echo non-interactively:

### Send a file:

```bash
./echo --mode send --local-port 9000 --remote-addr 192.168.1.5:9001 --file mydoc.pdf
```

### Receive a file:

```bash
./echo --mode receive --local-port 9001 --remote-addr 192.168.1.4:9000
```

## 4. Verify Operation

Make sure:

- The ports are open on both ends.
- Both peers run the compatible version.
- They're on reachable networks (LAN or port-forwarded public IP).

## 5. Help

See available flags:

```bash
./echo --help
```
