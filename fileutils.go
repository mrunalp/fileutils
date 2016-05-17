package fileutils

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

// CopyFile copies the file at source to dest
func CopyFile(source string, dest string) error {
	si, err := os.Lstat(source)
	if err != nil {
		return err
	}

	// Handle symlinks
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(source)
		if err != nil {
			return err
		}
		if err := os.Symlink(target, dest); err != nil {
			return err
		}
		return nil
	}

	// Handle device files
	st, ok := si.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("could not convert to syscall.Stat_t")
	}

	if st.Mode&syscall.S_IFBLK != 0 || st.Mode&syscall.S_IFCHR != 0 {
		devMajor := int64(major(uint64(st.Rdev)))
		devMinor := int64(minor(uint64(st.Rdev)))
		mode := uint32(si.Mode() & 07777)
		if st.Mode&syscall.S_IFBLK != 0 {
			fmt.Printf("BLOCK")
			mode |= syscall.S_IFBLK
		}
		if st.Mode&syscall.S_IFCHR != 0 {
			fmt.Printf("CHAR")
			mode |= syscall.S_IFCHR
		}
		if err := syscall.Mknod(dest, mode, int(mkdev(devMajor, devMinor))); err != nil {
			return err
		}
		return nil
	}

	// Handle regular files
	if si.Mode().IsRegular() {
		sf, err := os.Open(source)
		if err != nil {
			return err
		}
		defer sf.Close()

		df, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer df.Close()

		_, err = io.Copy(df, sf)
		if err != nil {
			return err
		}
	}

	return nil
}

func major(device uint64) uint64 {
	return (device >> 8) & 0xfff
}

func minor(device uint64) uint64 {
	return (device & 0xff) | ((device >> 12) & 0xfff00)
}

func mkdev(major int64, minor int64) uint32 {
	return uint32(((minor & 0xfff00) << 12) | ((major & 0xfff) << 8) | (minor & 0xff))
}
