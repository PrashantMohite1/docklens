# Deep Dive: Docker

## Why Docker needs Linux

Containers are not a Docker invention. They are built out of **Linux kernel features**:

- **namespaces** — make a process think it has its own filesystem, network, process list
- **cgroups** — limit CPU/memory
- **overlayfs** — stack image layers into one filesystem

These features **only exist in the Linux kernel.** Windows doesn't have them. So a Linux container literally cannot run directly on the Windows kernel — the machinery it depends on isn't there.

### So how does Docker run on Windows? → A hidden Linux VM

Since Windows can't provide Linux kernel features, Docker Desktop quietly runs a **real Linux machine inside your Windows machine**:

```text
┌───────────────────────────────────────────────────────────┐
│  Windows (your laptop)                                    │
│                                                           │
│   docker.exe  ← the client you type commands into         │
│        │                                                  │
│        │  talks to                                        │
│        ▼                                                  │
│   ┌──────────────────────────────────────────────────┐    │
│   │  Linux VM  (name: "docker-desktop", runs on WSL2)│    │
│   │                                                  │    │
│   │    dockerd (daemon)                              │    │
│   │    /var/lib/docker   ← images/layers live HERE   │    │
│   │    Linux kernel: namespaces, cgroups, overlayfs  │    │
│   └──────────────────────────────────────────────────┘    │
└───────────────────────────────────────────────────────────┘
```

- **WSL2** = "Windows Subsystem for Linux v2" = the technology that runs a genuine Linux kernel inside Windows in a lightweight VM.
- Docker Desktop uses a WSL2 distro literally named `**docker-desktop**`.

So the daemon, the images, the layers, the running containers — **all of that lives inside that Linux VM**, not on your Windows `C:\` drive.

## Image directory structure

When you run `docker save alpine -o alpine.tar` and extract it, this is the full tree:

```text
alpine.tar                         ← portable "suitcase" (a plain tar, not compressed)
   │  extract (tar -xf)
   ▼
alpine/
├── oci-layout                     ← marks the format version ({"imageLayoutVersion":"1.0.0"})
├── index.json                     ← entry point → points to the image index blob
├── manifest.json                  ← convenience list: which config + which layers
└── blobs/
    └── sha256/                    ← every file named by the SHA-256 of its contents
        ├── 28bd5fe8…              ← image INDEX  (which manifest for which CPU arch)
        ├── 79ff19e9…              ← image MANIFEST (points to config + layers)
        ├── d529dd0c…              ← image CONFIG  (env, cmd, arch, diff_ids)
        └── 55afa1ec…              ← the LAYER  (a .tar.gz)
                │  untar (tar -xzf)
                ▼
            rootfs/                ← the EXACT root filesystem we run inside the container
            ├── bin   ├── etc   ├── lib   ├── usr   ├── var
            ├── dev   ├── home  ├── mnt   ├── sbin  ├── ...
```

### How this maps to "how the image is created underneath"

Everything is just **files referencing the hash of the file below it**, built bottom-up:

```text
1. real files (/bin, /etc, ...)  →  tar          →  layer.tar      →  sha256 = diffID
2. layer.tar                     →  gzip          →  layer.tar.gz   →  sha256 = digest   → stored as blob 55afa1ec…
3. write CONFIG json (env, cmd, diff_ids)                          →  stored as blob d529dd0c…
4. write MANIFEST json (points to config + layer blobs)            →  stored as blob 79ff19e9…
5. write INDEX json (points to the manifest)                       →  stored as blob 28bd5fe8…
6. write index.json + oci-layout, attach a tag (alpine:latest)
```

- The **files** come first (a plain directory).
- Each level above (layer tar → config → manifest → index) is just a **json/tar wrapper that references the hash of the thing below it**.
- Change one file → its hash changes → the config's diffID changes → the manifest changes → the index changes. This is why images are **tamper-evident**.

There is no magic build engine at the core. The image format is simply: **folder → tar → gzip → hash → reference in json → store as blobs**.



---



## WSL timeline (short)


| Year      | Milestone                                          |
| --------- | -------------------------------------------------- |
| 2016      | WSL1 announced ("Bash on Ubuntu on Windows")       |
| 2017–2018 | More distros, better syscall coverage              |
| 2019      | WSL2 announced (real Linux kernel, lightweight VM) |
| 2020+     | Docker Desktop moves hard to WSL2 backend          |
| 2021+     | WSLg (Linux GUI apps on Windows)                   |
| Ongoing   | Better networking, systemd support, GPU, etc.      |


- **WSL1** = Linux apps via **syscall translation** (clever, limited)
- **WSL2** = Linux apps via a **real Linux kernel in a tiny VM** (compatible, Docker-ready)
- **Docker Desktop on Windows** = Windows UI/CLI + **Linux engine living in WSL2**

## A "distro" is NOT a whole OS

A Linux system is really **two separate layers**: userspace (the "distro") and the kernel (the engine).

```text
┌─────────────────────────────────────────┐
│  USERSPACE  (this is what a "distro" is)  │
│  - /bin, /usr, /etc, /home                │
│  - apt vs yum vs pacman                   │
│  - bash, python, package versions         │
│  - the "Ubuntu-ness" or "CentOS-ness"     │
└─────────────────────────────────────────┘
                    │  makes syscalls to ↓
┌─────────────────────────────────────────┐
│  KERNEL  (the actual engine)              │
│  - process scheduling, memory, drivers    │
│  - namespaces, cgroups, overlayfs         │
└─────────────────────────────────────────┘
```

- Ubuntu, CentOS, Kali, Debian mostly use the **same Linux kernel** — the kernel isn't "Ubuntu" or "CentOS".
- What makes distros *feel* different is entirely **userspace**: files, package manager (`apt` vs `yum`), configs, versions.
- A distro ≈ **root filesystem + package manager + defaults**, sitting on a Linux kernel.
- "Only one distro at a time" on bare metal is **not a Linux law** — it's just that booting picks **one** root filesystem to be `/`.

## How WSL2 runs multiple distros at once

WSL2 = **one real Linux kernel** in a lightweight VM. That single kernel mounts **many isolated root filesystems** — each one is a distro.

```text
┌──────────────── WSL2 lightweight VM ────────────────┐
│                                                     │
│   ONE shared Linux kernel (6.18.x)                  │
│        ▲          ▲          ▲          ▲           │
│        │          │          │          │           │
│   ┌────────┐ ┌────────┐ ┌──────────┐ ┌──────────┐  │
│   │ Ubuntu │ │  Kali  │ │ docker-  │ │ docker-  │  │
│   │ rootfs │ │ rootfs │ │ desktop  │ │ desktop- │  │
│   │        │ │        │ │ rootfs   │ │ data     │  │
│   └────────┘ └────────┘ └──────────┘ └──────────┘  │
│   each is just a separate Linux filesystem          │
└─────────────────────────────────────────────────────┘
```

- All distros **share the one kernel** but keep their **own separate userspace** → they run simultaneously.
- Each distro is stored as its own virtual disk on Windows (`ext4.vhdx`):

```text
Ubuntu           → ...\CanonicalGroup...\ext4.vhdx
docker-desktop   → ...\docker-desktop\ext4.vhdx
```

- `wsl -d Ubuntu` → starts shared kernel/VM → mounts Ubuntu's `ext4.vhdx` as root → drops into Ubuntu userspace. Switch distro = same kernel, different rootfs.

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

