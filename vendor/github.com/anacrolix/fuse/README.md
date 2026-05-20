bazil.org/fuse -- Filesystems in Go
===================================

This fork has **support for FUSE on Mac**. MacFUSE 4.0.0 and 3.3 (or newer) are supported.

The original project **dropped support** for FUSE on Mac when OSXFUSE stopped being an open source project [bazil/fuse#224](https://github.com/bazil/fuse/issues/224). I respect the maintainers decisions of both projects.

In this fork, the following patches to _remove_ support for OSXFUSE have been dropped:

* [60eaf8](https://github.com/bazil/fuse/commit/60eaf8f021ce00e5c52529cdcba1067e13c1c2c2) - Remove macOS support
* [eca21f](https://github.com/bazil/fuse/commit/eca21f36f00e04957de26b2e64e21544fa0e0372) - Comment cleanup after macOS support removal

After forking, a patch to introduce support for macFUSE 4 has been made [#1](https://github.com/zegl/fuse/pull/1).

To use this fork in your project: `go get github.com/zegl/fuse`

---

`bazil.org/fuse` is a Go library for writing FUSE userspace
filesystems.

It is a from-scratch implementation of the kernel-userspace
communication protocol, and does not use the C library from the
project called FUSE. `bazil.org/fuse` embraces Go fully for safety and
ease of programming.

Here’s how to get going:

    go get bazil.org/fuse

Website: http://bazil.org/fuse/

Github repository: https://github.com/bazil/fuse

API docs: http://godoc.org/bazil.org/fuse

Our thanks to Russ Cox for his fuse library, which this project is
based on.
