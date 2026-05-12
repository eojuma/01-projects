# 🚀 7-Day Build Plan for `my-ls`

A step-by-step plan to build your own version of the Unix `ls` command in Go.

---

## 📊 Progress Tracker

- [ ] Day 1 — Basic ls
- [ ] Day 2 — Path Support
- [ ] Day 3 — Hidden Files (-a)
- [ ] Day 4 — Sorting (-r)
- [ ] Day 5 — Time Sorting (-t)
- [ ] Day 6 — Long Format (-l)
- [ ] Day 7 — Polish + Bonus

---

## 📅 DAY 1 — Basic ls (Foundation)

**Goal:**
Print files in the current directory

**Tasks:**

- [ ] Read directory contents
- [ ] Print file names only

**Learn + Use:**

- `os.ReadDir()`
- loops

---

## 📅 DAY 2 — Add Path Support

**Goal:**
Support:

```bash
./my-ls folder/
```

**Tasks:**

- [ ] Read from `os.Args`
- [ ] If no path → use `"."`
- [ ] If path exists → list that directory

**Learn:**

- `os.Args`
- basic error handling

---

## 📅 DAY 3 — Add `-a` (Hidden Files)

**Goal:**
Show hidden files

**Tasks:**

- [ ] Detect flag `-a`
- [ ] Include files starting with `.`

**Learn:**

- `strings.HasPrefix()`

---

## 📅 DAY 4 — Sorting (`-r` and Default)

**Goal:**
Sort files

**Tasks:**

- [ ] Default → sort by name (A → Z)
- [ ] Add `-r` → reverse order

**Learn:**

- Sorting slices
- Comparing strings

---

## 📅 DAY 5 — Time Sorting (`-t`)

**Goal:**
Sort by last modified time

**Tasks:**

- [ ] Use `os.Stat()` to get file time
- [ ] Sort by newest first

**Learn:**

- `.ModTime()`
- sorting with time

---

## 📅 DAY 6 — `-l` (Long Format)

**Goal:**
Replicate `ls -l`

**Tasks:**

- [ ] File permissions
- [ ] Number of links
- [ ] Owner (user)
- [ ] Group
- [ ] File size
- [ ] Last modified time
- [ ] File name

**Learn:**

- `os.FileInfo`
- `os/user`
- `syscall`
- formatting with `fmt.Printf`

---

## 📅 DAY 7 — Polish + Bonus

**Goal:**
Improve + optional features

---

### ⭐ Bonus 1: `-R` (Recursive)

**Tasks:**

- [ ] List subdirectories recursively

**Learn:**

- recursion

---

### 🎨 Bonus 2: Colors

**Tasks:**

- [ ] Add color output
- [ ] Blue → directories
- [ ] Default → files

---

## 🧠 Notes

- Keep functions small and focused
- Test each feature before moving to the next day
- Build incrementally — don’t jump ahead

---

## 🏁 Final Goal

By the end, your CLI should support:

```bash
my-ls
my-ls folder/
my-ls -a
my-ls -r
my-ls -t
my-ls -l
my-ls -R
```

---
