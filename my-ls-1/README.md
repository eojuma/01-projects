# 📁 my-ls — Project Summary

## 📖 Overview

`my-ls` is a Go-based reimplementation of the Unix `ls` command.
The goal is to replicate the behavior of the original command, which lists files and directories within a specified path (or the current directory by default).

---

## 🎯 Objectives

* Recreate core functionality of the `ls` command
* Handle directory listing and file metadata
* Support multiple flags and combinations, similar to the real `ls`
* Ensure output (especially for `-l`) closely matches the system command

---

## ⚙️ Required Features

Your implementation must support the following flags:

* `-l` → Long listing format (must match system `ls -l`)
* `-a` → Show hidden files (starting with `.`)
* `-r` → Reverse sorting order
* `-t` → Sort files by last modified time

### 📌 Notes

* Flags can be combined (e.g. `-la`, `-ltr`)
* Default behavior (no arguments) → list current directory

---

## 🧩 Constraints

* Must be written in **Go**
* Must follow **good coding practices**
* **Unit tests are recommended**
* ❌ Use of `os/exec` is **strictly prohibited**

---

## 🌟 Bonus Features

You can extend functionality with:

* `-R` → Recursive directory listing
* 🎨 Colored output (e.g. directories in blue)
* Additional flags beyond the required ones

---

## 📦 Allowed Packages

You may only use the following standard library packages:

* `fmt`
* `os`
* `os/user`
* `strconv`
* `strings`
* `syscall`
* `time`
* `math/rand`
* `errors`
* `io/fs`

---

## 💡 Hints

* Plan for recursive behavior (`-R`) early in your design
* Maintain correct **file ordering**
* Refer to the system `ls` manual (`man ls`) for expected behavior
* Pay attention to formatting details in `-l` output

---

## ▶️ Usage

Run your program similarly to:

```bash
./my-ls
./my-ls folder/
./my-ls -l
./my-ls -a
./my-ls -ltr
```

---

## 🚀 Getting Started (How to Run)

### 1. Clone the repository

```bash
git clone https://github.com/your-username/my-ls.git
cd my-ls
```

### 2. Initialize Go module (if not already)

```bash
go mod init my-ls
go mod tidy
```

### 3. Build the project

```bash
go build -o my-ls
```

### 4. Run the program

```bash
go run . [-flags] path/to/list
```

### 5. (Optional) Run without building

```bash
go run . -- -l
```

---

## 📚 Learning Outcomes

This project helps you understand:

* 🐧 Unix system behavior
* 📥 Handling input (arguments, flags)
* 📤 Producing formatted output
* 🔤 String manipulation
* 🧱 Working with structs and file metadata

---

## 👥 Contributors

* Margaret Kerubo(@makerubo)
* Silas Lelei(@ksilas)
* Evans Juma(@evajuma)

---

