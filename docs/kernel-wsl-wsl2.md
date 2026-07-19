# Kernel, WSL & WSL2

How the Linux and Windows kernels differ, and how WSL / WSL2 bridge them.
This is the foundation under Docker Desktop on Windows (see `deep-dive-docker.md`).

## Contents

- [What is a kernel?](#what-is-a-kernel)
- [Windows kernel vs Linux kernel](#windows-kernel-vs-linux-kernel)
- [How WSL solves it](#how-wsl-solves-it)
  - [WSL1 — translation](#wsl1--translation)
  - [WSL2 — real Linux kernel in a lightweight VM](#wsl2--real-linux-kernel-in-a-lightweight-vm)
  - [WSL1 vs WSL2](#wsl1-vs-wsl2)
- [A "distro" is NOT a whole OS](#a-distro-is-not-a-whole-os)
- [How WSL2 runs multiple distros at once](#how-wsl2-runs-multiple-distros-at-once)
  - [Same trick as Docker](#same-trick-as-docker)
  - [Quick proof](#quick-proof)
- [WSL timeline (short)](#wsl-timeline-short)
- [Why this matters for DockLens](#why-this-matters-for-docklens)

Related doc: [Deep Dive: Docker](./deep-dive-docker.md)

## What is a kernel?

The **kernel** is the core of an operating system — the only layer allowed to
talk to hardware. Apps ask it for everything through **system calls (syscalls)**.

```text
┌──────────────────────────────────────────┐
│  Apps  (browser, python, docker, ...)    │
└──────────────────────────────────────────┘
                      │  syscalls: open file, get memory, send packet
                      ▼
┌──────────────────────────────────────────┐
│  KERNEL                                  │
│    - process scheduling (who runs when)  │
│    - memory management                   │
│    - filesystems                         │
│    - networking stack                    │
│    - device drivers                      │
└──────────────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────┐
│  Hardware  (CPU, RAM, disk, network)     │
└──────────────────────────────────────────┘
```

- Apps never touch hardware directly — they make **syscalls** (open file, get memory, send packet).
- The kernel serves those requests while enforcing isolation & security.
- **Key insight:** a compiled program targets **one specific kernel's syscalls**. That single fact is why a Linux binary can't just run on Windows.

## Windows kernel vs Linux kernel

Two completely different kernels — different authors, different syscalls. A binary
built for one cannot talk to the other.

|                     | **Linux kernel**              | **Windows kernel (NT)**     |
| ------------------- | ----------------------------- | --------------------------- |
| Origin              | Linus Torvalds, 1991, open    | Microsoft (NT), 1993, closed |
| Open a file         | `open()`                      | `NtCreateFile()`            |
| Create a process    | `fork()` + `execve()`         | `CreateProcess()`           |
| Container primitives| **namespaces, cgroups, overlayfs** | none (different model) |
| Default filesystem  | ext4, xfs, ...                | NTFS                        |

Two problems this creates:

1. **Different syscalls** — a Linux binary calls functions the Windows kernel simply doesn't have.
2. **Missing features** — containers need **namespaces, cgroups, overlayfs**; these are Linux-kernel-specific and Windows has no reusable equivalent.

## How WSL solves it

WSL bridges the gap in two very different ways across its generations — this is the
clearest way to understand WSL1 vs WSL2.

### WSL1 — translation

Keep the Windows kernel; add a layer that converts Linux syscalls into NT syscalls on the fly.

```text
Linux app
   │   Linux syscall:  open()
   ▼
WSL1 translation layer          ← catches & converts
   │   NT syscall:     NtCreateFile()
   ▼
Windows NT kernel
```

- No VM, lightweight — but syscall coverage is **incomplete**.
- **Docker didn't work well:** there's nothing to translate `namespaces`/`cgroups`/`overlayfs` *to* — Windows lacks them.
- Mental model: a **live interpreter** translating sentence-by-sentence — fine for common phrases, breaks on complex ones.

### WSL2 — real Linux kernel in a lightweight VM

Stop translating. Microsoft ships a **real Linux kernel** and runs it in a fast, lightweight VM.

```text
┌───────────────────────────────────────────────────┐
│  Windows  — Windows apps talk to the NT kernel    │
│                                                   │
│  lightweight VM (started/stopped by wsl):         │
│  ┌─────────────────────────────────────────────┐  │
│  │  REAL Linux kernel (6.18.x)                 │  │
│  │    → real namespaces / cgroups / overlayfs  │  │
│  │  Linux apps make Linux syscalls to a real   │  │
│  │  Linux kernel  — no translation needed      │  │
│  └─────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────┘
```

- Linux apps now hit an **actual Linux kernel** → **full syscall compatibility**.
- Container features **exist for real** → Docker works properly.
- Mental model: a **native speaker in the room**, not an interpreter.

### WSL1 vs WSL2

|                          | WSL1 (translate) | WSL2 (real kernel in VM) |
| ------------------------ | ---------------- | ------------------------ |
| Linux kernel present?    | No               | **Yes (real)**           |
| Syscall compatibility    | Partial          | **Full**                 |
| Docker / containers      | Poor             | **Works**                |
| Access to Windows `C:\`  | Fast (native)    | Slower (crosses VM edge) |
| Startup / feel           | Very light       | Light VM, still fast     |

- **WSL1** = Linux apps via **syscall translation** (clever, limited)
- **WSL2** = Linux apps via a **real Linux kernel in a tiny VM** (compatible, Docker-ready)
- **Docker Desktop on Windows** = Windows UI/CLI + **Linux engine living in WSL2**

## A "distro" is NOT a whole OS

A Linux system is really **two layers**: userspace (the "distro") and the kernel (the engine).

```text
┌───────────────────────────────────────────┐
│  USERSPACE  (this is what a "distro" is)  │
│    - /bin, /usr, /etc, /home              │
│    - apt vs yum vs pacman                 │
│    - bash, python, package versions       │
│    - the "Ubuntu-ness" or "CentOS-ness"   │
└───────────────────────────────────────────┘
                      │  makes syscalls to ↓
┌───────────────────────────────────────────┐
│  KERNEL  (the actual engine)              │
│    - process scheduling, memory, drivers  │
│    - namespaces, cgroups, overlayfs       │
└───────────────────────────────────────────┘
```

- Ubuntu, CentOS, Kali, Debian mostly share the **same Linux kernel** — the kernel isn't "Ubuntu" or "CentOS".
- What makes distros *feel* different is entirely **userspace**: files, package manager (`apt` vs `yum`), configs, versions.
- A distro ≈ **root filesystem + package manager + defaults**, sitting on a Linux kernel.
- "Only one distro at a time" on bare metal is **not a Linux law** — booting just picks **one** root filesystem to be `/`.

## How WSL2 runs multiple distros at once

WSL2 = **one real Linux kernel** in a lightweight VM. That single kernel mounts
**many isolated root filesystems** — each one is a distro.

```text
┌───────────────────────────────────────────────────────────────────────┐
│   ONE shared Linux kernel (6.18.x)                                    │
│   shared by every distro below                                        │
│                                                                       │
│   ┌────────────┐   ┌────────────┐   ┌────────────┐   ┌────────────┐   │
│   │  Ubuntu    │   │  Kali      │   │  docker-   │   │  docker-   │   │
│   │  rootfs    │   │  rootfs    │   │  desktop   │   │  desktop-  │   │
│   └────────────┘   └────────────┘   │  rootfs    │   │  data      │   │
│                                     └────────────┘   └────────────┘   │
│                                                                       │
│   each box = one distro (its own root filesystem)                     │
└───────────────────────────────────────────────────────────────────────┘
```

- All distros **share the one kernel** but keep their **own separate userspace** → they run simultaneously.
- Each distro is stored as its own virtual disk on Windows (`ext4.vhdx`):

```text
Ubuntu           → ...\CanonicalGroup...\ext4.vhdx
docker-desktop   → ...\docker-desktop\ext4.vhdx
```

- `wsl -d Ubuntu` → start shared kernel/VM → mount Ubuntu's `ext4.vhdx` as root → drop into its userspace. Switch distro = same kernel, different rootfs.

### Same trick as Docker

- **WSL distro** = a root filesystem running against the **shared WSL2 kernel**.
- **Docker container** = a root filesystem (image `rootfs/`) running against the **host kernel**, isolated by namespaces/cgroups.
- Same fundamental idea: **swap the userspace, share the kernel**. That's why `docker-desktop` is itself just a WSL2 distro.

### Quick proof

```powershell
wsl -l -v                      # multiple distros, all VERSION 2
wsl -d docker-desktop uname -r # same kernel version...
wsl -d Ubuntu uname -r         # ...as this one → shared kernel
```

## WSL timeline (short)

| Year      | Milestone                                          |
| --------- | -------------------------------------------------- |
| 2016      | WSL1 announced ("Bash on Ubuntu on Windows")       |
| 2017–2018 | More distros, better syscall coverage              |
| 2019      | WSL2 announced (real Linux kernel, lightweight VM) |
| 2020+     | Docker Desktop moves hard to WSL2 backend          |
| 2021+     | WSLg (Linux GUI apps on Windows)                   |
| Ongoing   | Better networking, systemd support, GPU, etc.      |

## Why this matters for DockLens

- On Windows, the real container truth (image blobs, overlay mounts, namespaces, cgroup limits) lives **inside the WSL2 Linux world**, not on `C:\`.
- DockLens must reason about that Linux-side reality even though the commands start from Windows.
