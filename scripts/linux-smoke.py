#!/usr/bin/env python3
"""Optional local Linux/BusyBox/GNU check. Requires curl, Go, bsdtar, and QEMU.

Downloads Alpine from its official HTTPS mirror, verifies rootfs/ISO SHA-256,
extracts GNU packages into a temporary RAM-only guest, and boots without networking.
No cloud resources, host mounts, disks, or privileged host operations are used.
"""
from pathlib import Path
import gzip
import hashlib
import io
import json
import os
import re
import shutil
import stat
import subprocess
import tarfile
import tempfile

ROOT = Path(__file__).resolve().parents[1]
MIRROR = "https://dl-cdn.alpinelinux.org/alpine"
BRANCH = "v3.24"
RELEASE = "3.24.3"
ARTIFACTS = {
    f"alpine-minirootfs-{RELEASE}-x86_64.tar.gz": "f8e1cafba9b1b283b5d30d06f9885c8d2bc161ffe9083ade8e2e157813a53407",
    f"alpine-virt-{RELEASE}-x86_64.iso": "cdc1bc0efc64d8d611bf3733f54868172a545888dfd726f518bf279bcb0ff1af",
}


def download(url, destination):
    subprocess.run(["curl", "-fsSL", "--retry", "2", url, "-o", str(destination)], check=True)
    return destination


def unpack(archive, root):
    # Only unpack official archive data into a new, disposable guest root.
    with tarfile.open(archive, "r:gz", ignore_zeros=True) as tar:
        for member in tar:
            name = member.name.removeprefix("./")
            if not name or name.startswith("."):
                continue
            relative = Path(name)
            if relative.is_absolute() or ".." in relative.parts:
                raise ValueError(f"Unsafe archive path: {name}")
            dest = root / relative
            for parent in dest.parents:
                if parent == root:
                    break
                if parent.is_symlink():
                    raise ValueError(f"Archive crosses symlink: {name}")
            if member.isdir():
                dest.mkdir(parents=True, exist_ok=True)
            elif member.isfile():
                dest.parent.mkdir(parents=True, exist_ok=True)
                if dest.is_symlink():
                    dest.unlink()
                dest.write_bytes(tar.extractfile(member).read())
                dest.chmod(member.mode)
            elif member.issym():
                dest.parent.mkdir(parents=True, exist_ok=True)
                if dest.is_symlink() or dest.is_file():
                    dest.unlink()
                dest.symlink_to(member.linkname)


def cpio(root, destination):
    stream = io.BytesIO()
    inode = 1

    def entry(name, mode, data=b"", major=0, minor=0):
        nonlocal inode
        encoded = name.encode() + b"\0"
        fields = [inode, mode, 0, 0, 1, 0, len(data), 0, 0, major, minor, len(encoded), 0]
        stream.write(b"070701" + b"".join(f"{n:08x}".encode() for n in fields))
        stream.write(encoded)
        stream.write(b"\0" * (-stream.tell() % 4))
        stream.write(data)
        stream.write(b"\0" * (-stream.tell() % 4))
        inode += 1

    for path in sorted(root.rglob("*")):
        metadata = path.lstat()
        name = str(path.relative_to(root))
        if path.is_symlink():
            data = os.readlink(path).encode()
        elif path.is_file():
            data = path.read_bytes()
        else:
            data = b""
        entry(name, metadata.st_mode, data)
    entry("dev/console", stat.S_IFCHR | 0o600, major=5, minor=1)
    entry("dev/null", stat.S_IFCHR | 0o666, major=1, minor=3)
    entry("TRAILER!!!", 0)
    destination.write_bytes(gzip.compress(stream.getvalue(), mtime=0))


def main():
    qemu = shutil.which("qemu-system-x86_64")
    if not qemu:
        raise SystemExit("Install QEMU separately before running this optional check.")
    bsdtar = shutil.which("bsdtar")
    if not bsdtar:
        raise SystemExit("Install bsdtar (libarchive tools) to extract the Alpine ISO.")
    cache = ROOT / ".cache/linux-smoke"
    cache.mkdir(parents=True, exist_ok=True)
    for filename, digest in ARTIFACTS.items():
        artifact = cache / filename
        if not artifact.exists():
            download(f"{MIRROR}/{BRANCH}/releases/x86_64/{filename}", artifact)
        if hashlib.sha256(artifact.read_bytes()).hexdigest() != digest:
            raise SystemExit(f"SHA-256 mismatch: {filename}")

    with tempfile.TemporaryDirectory(prefix="shell-explain-linux-") as temporary:
        work = Path(temporary)
        root = work / "root"
        root.mkdir()
        unpack(cache / f"alpine-minirootfs-{RELEASE}-x86_64.tar.gz", root)
        subprocess.run([bsdtar, "-xf", str(cache / f"alpine-virt-{RELEASE}-x86_64.iso"), "-C", str(work), "boot/vmlinuz-virt"], check=True)
        index_file = download(f"{MIRROR}/{BRANCH}/main/x86_64/APKINDEX.tar.gz", work / "APKINDEX.tar.gz")
        with tarfile.open(index_file) as tar:
            index = tar.extractfile("APKINDEX").read().decode()
        packages = {}
        providers = {}
        for paragraph in index.split("\n\n"):
            record = dict(line.split(":", 1) for line in paragraph.splitlines() if ":" in line)
            if "P" not in record:
                continue
            packages[record["P"]] = record
            for provided in record.get("p", "").split():
                providers[provided.split("=")[0]] = record["P"]
        pending = ["grep", "sed", "findutils"]
        selected = set()
        while pending:
            name = pending.pop()
            if name in selected:
                continue
            record = packages[name]
            selected.add(name)
            for dependency in record.get("D", "").split():
                dependency = re.split(r"[=<>~]", dependency)[0]
                if dependency.startswith("!"):
                    continue
                provider = providers.get(dependency, dependency)
                if provider in packages:
                    pending.append(provider)
        for name in sorted(selected):
            record = packages[name]
            filename = f"{name}-{record['V']}.apk"
            artifact = download(f"{MIRROR}/{BRANCH}/main/x86_64/{filename}", work / filename)
            unpack(artifact, root)
        print("Alpine GNU package versions:", json.dumps({n: packages[n]["V"] for n in sorted(selected)}), flush=True)
        env = {**os.environ, "GOOS": "linux", "GOARCH": "amd64", "CGO_ENABLED": "0"}
        for target, command in [("explain", "./cmd/explain"), ("compat-check", "./cmd/compat-check")]:
            subprocess.run(["go", "build", "-trimpath", "-o", str(root / target), command], cwd=ROOT, env=env, check=True)
        init = root / "init"
        init.write_text("#!/bin/sh\nexport PATH=/usr/bin:/bin:/sbin\nmount -t proc proc /proc\n/compat-check\nresult=$?\necho GUEST_EXIT=$result\npoweroff -f\n")
        init.chmod(0o755)
        initrd = work / "guest.cpio.gz"
        cpio(root, initrd)
        result = subprocess.run([qemu, "-machine", "q35", "-accel", "tcg", "-m", "512", "-nographic", "-no-reboot", "-nic", "none", "-kernel", str(work / "boot/vmlinuz-virt"), "-initrd", str(initrd), "-append", "console=ttyS0 rdinit=/init quiet"], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, timeout=180)
        print(result.stdout)
        if "LINUX_SMOKE_PASS" not in result.stdout or "GUEST_EXIT=0" not in result.stdout:
            raise SystemExit("Linux smoke check failed")


if __name__ == "__main__":
    main()
