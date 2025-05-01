# Installation Guide

This guide explains how to install Echo using prebuilt binaries from the [GitHub Releases](https://github.com/IdanKoblik/Echo/releases) page.

## Windows

1. Go to the [Releases page](https://github.com/IdanKoblik/Echo/releases).
2. Download the latest `echo.exe`.
3. (Optional) Add it to your system PATH for global use.
4. Run it via:

```powershell
.\echo.exe
```

## macOS

1. Go to the [Releases page](https://github.com/IdanKoblik/Echo/releases).
2. Download the `echo-macos` binary.
3. Give it execute permission:

```bash
chmod +x echo-macos
./echo-macos
```

## Linux

1. Go to the [Releases page](https://github.com/IdanKoblik/Echo/releases).
2. Download the `echo-linux` binary.
3. Make it executable:

```bash
chmod +x echo-linux
./echo-linux
```

(Optional) Move to `/usr/local/bin` to use it globally:

```bash
sudo mv echo-linux /usr/local/bin/echo
```