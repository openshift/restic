package fuse

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"
)

func mount(mountPoint string, conf *mountConfig, ready chan<- struct{}, errp *error) (*os.File, error) {
	locations := conf.osxfuseLocations
	if locations == nil {
		locations = []OSXFUSEPaths{
			OSXFUSELocationV4,
			OSXFUSELocationV3,
		}
	}

	var binLocation string
	for _, loc := range locations {
		if _, err := os.Stat(loc.Mount); os.IsNotExist(err) {
			// try the other locations
			continue
		}
		binLocation = loc.Mount
		break
	}
	if binLocation == "" {
		return nil, ErrOSXFUSENotFound
	}

	local, remote, err := unixgramSocketpair()
	if err != nil {
		return nil, err
	}

	defer local.Close()
	defer remote.Close()

	cmd := exec.Command(binLocation,
		"-o", conf.getOptions(),

		// Tell osxfuse-kext how large our buffer is. It must split
		// writes larger than this into multiple writes.
		//
		// OSXFUSE seems to ignore InitResponse.MaxWrite, and uses
		// this instead.
		"-o", "iosize="+strconv.FormatUint(maxWrite, 10),

		mountPoint)

	cmd.ExtraFiles = []*os.File{remote} // fd would be (index + 3)
	cmd.Env = append(os.Environ(),
		"_FUSE_CALL_BY_LIB=",
		"_FUSE_DAEMON_PATH="+os.Args[0],
		"_FUSE_COMMFD=3",
		"_FUSE_COMMVERS=2",
		"MOUNT_OSXFUSE_CALL_BY_LIB=",
		"MOUNT_OSXFUSE_DAEMON_PATH="+os.Args[0])

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	fd, err := getConnection(local)
	if err != nil {
		return nil, err
	}

	go func() {
		// wait inside a goroutine or otherwise it would block forever for unknown reasons
		if err := cmd.Wait(); err != nil {
			err = fmt.Errorf("mount_osxfusefs failed: %v. Stderr: %s, Stdout: %s",
				err, errOut.String(), out.String())
			*errp = err
		}
		close(ready)
	}()

	dup, err := syscall.Dup(int(fd.Fd()))
	if err != nil {
		return nil, err
	}

	syscall.CloseOnExec(int(fd.Fd()))
	syscall.CloseOnExec(dup)

	return os.NewFile(uintptr(dup), "macfuse"), err
}

func unixgramSocketpair() (l, r *os.File, err error) {
	fd, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, nil, os.NewSyscallError("socketpair",
			err.(syscall.Errno))
	}
	l = os.NewFile(uintptr(fd[0]), "socketpair-half1")
	r = os.NewFile(uintptr(fd[1]), "socketpair-half2")
	return
}

func getConnection(local *os.File) (*os.File, error) {
	var data [4]byte
	control := make([]byte, 4*256)

	// n, oobn, recvflags, from, errno  - todo: error checking.
	_, oobn, _, _,
	err := syscall.Recvmsg(
		int(local.Fd()), data[:], control[:], 0)
	if err != nil {
		return nil, err
	}

	message := *(*syscall.Cmsghdr)(unsafe.Pointer(&control[0]))
	fd := *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&control[0])) + syscall.SizeofCmsghdr))

	if message.Type != syscall.SCM_RIGHTS {
		return nil, fmt.Errorf("getConnection: recvmsg returned wrong control type: %d", message.Type)
	}
	if oobn <= syscall.SizeofCmsghdr {
		return nil, fmt.Errorf("getConnection: too short control message. Length: %d", oobn)
	}
	if fd < 0 {
		return nil, fmt.Errorf("getConnection: fd < 0: %d", fd)
	}

	return os.NewFile(uintptr(fd), "macfuse"), nil
}
