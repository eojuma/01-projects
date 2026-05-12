package internal

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// PrintLong prints files in long format (-l).
// showTotal determines if the "total X" block count line is printed.
func PrintLong(path string, files []string, showTotal bool) {
	var total int64

	// ---- calculate total ----
	if showTotal {
		for _, name := range files {
			fullPath := JoinPaths(path, name)
			info, err := os.Lstat(fullPath)
			if err != nil {
				continue
			}

			sys, ok := info.Sys().(*syscall.Stat_t)
			if ok {
				total += sys.Blocks
			}
		}
		// print total in 1K blocks (sys.Blocks is usually 512-byte blocks)
		fmt.Println("total", total/2)
	}

	// ---- print each file ----
	for _, name := range files {
		fullPath := JoinPaths(path, name)
		info, err := os.Lstat(fullPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		sys, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			continue
		}

		var out strings.Builder

		// ---- permissions & file type ----
		// Manually determine the first character to avoid "dcrw" issues
		var typeChar byte = '-'
		m := info.Mode()
		if m.IsDir() {
			typeChar = 'd'
		} else if m&os.ModeCharDevice != 0 {
			typeChar = 'c'
		} else if m&os.ModeDevice != 0 {
			typeChar = 'b'
		} else if m&os.ModeSymlink != 0 {
			typeChar = 'l'
		} else if m&os.ModeNamedPipe != 0 {
			typeChar = 'p'
		} else if m&os.ModeSocket != 0 {
			typeChar = 's'
		}

		// Use the last 9 characters of Go's Mode string for the permissions
		modeStr := info.Mode().String()
		perms := string(typeChar) + modeStr[len(modeStr)-9:]
		out.WriteString(perms + " ")

		// ---- links ----
		out.WriteString(strconv.FormatUint(uint64(sys.Nlink), 10) + " ")

		// ---- user and group ----
		usr, _ := user.LookupId(strconv.Itoa(int(sys.Uid)))
		grp, _ := user.LookupGroupId(strconv.Itoa(int(sys.Gid)))
		out.WriteString(usr.Username + " " + grp.Name + " ")

		// ---- size or major/minor numbers ----
		// If it's a character or block device, show major/minor numbers instead of size
		if (m&os.ModeCharDevice != 0 || m&os.ModeDevice != 0) && !m.IsDir() {
			major := (sys.Rdev >> 8) & 0xfff
			minor := (sys.Rdev & 0xff) | ((sys.Rdev >> 12) & 0xfff00)
			out.WriteString(fmt.Sprintf("%3d, %3d ", major, minor))
		} else {
			// Use a single space after size instead of %8d padding
			out.WriteString(fmt.Sprintf("%d ", info.Size()))
		}

		// ---- modification time ----
		modTime := info.ModTime()
		now := time.Now()
		sixMonthsAgo := now.AddDate(0, -6, 0)

		if modTime.Before(sixMonthsAgo) || modTime.After(now) {
			out.WriteString(modTime.Format("Jan _2  2006") + " ")
		} else {
			out.WriteString(modTime.Format("Jan _2 15:04") + " ")
		}

		displayName := name

		// ---- symlink support ----
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(fullPath)
			if err == nil {
				displayName = "\033[1;36m" + displayName + "\033[0m -> \033[1;34m" + target + "\033[0m"
			}
		}

		// ---- coloring ----
		if info.IsDir() {
			out.WriteString("\033[1;34m" + displayName + "\033[0m")
		} else if info.Mode()&os.ModeSymlink != 0 {
			out.WriteString("\033[1;36m" + displayName + "\033[0m")
		} else if info.Mode()&0o111 != 0 {
			out.WriteString("\033[1;32m" + displayName + "\033[0m")
		} else {
			if strings.Contains(m.Perm().String(), "x") {
				out.WriteString("\033[1;32m" + displayName + "\033[0m")
			} else {
				out.WriteString(displayName)
			}
		}

		fmt.Println(out.String())
	}
}

func PrintShort(path string, files []string) {
	for _, name := range files {
		fullPath := JoinPaths(path, name)
		info, err := os.Lstat(fullPath)
		displayName := name

		if err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				displayName = "\033[1;96m" + name + "\033[0m"
			} else if info.IsDir() {
				displayName = "\033[1;34m" + name + "\033[0m"
			} else {
				if strings.Contains(info.Mode().Perm().String(), "x") {
					displayName = ("\033[1;32m" + name + "\033[0m")
				} else {
					displayName = name
				}
			}
		}
		fmt.Printf("%s  ", displayName)
	}
	fmt.Println()
}
