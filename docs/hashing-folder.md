# hash

hash is a **short mathematical fingerprint that represents the file**.

---

### How to check weather folder changed or not using hash

we can find out hash of command output - ex- echo "hello" | sha256sum , ls temp | sha256sum

folder hash find : find list of file in directory > then will get hash of each file > then will print sorted files list with hash in output and > will find hash of final sorted list output, which is our folder hash  

from folder hash we can easily find weather the folder is same or not. 

```
cd <folder-to-hash-check> 
```

run below command you'll get hash , going into parent directory is the must otherwise, if you ls inside directory and ls dir-name both output is different thats going in parent directory must needed otherwise you will get different output and hash will change.

```bash
find . -type f -exec sha256sum {} + | sort
```

produces **plain text** like:

```text
185f8db3...  ./file1.txt
4824c2f3...  ./file2.txt
a91bc4d1...  ./file3.txt
```

Now this command:

```bash
find . -type f -exec sha256sum {} + | sort | sha256sum
```

takes **that text output** as its input and hashes it.

So internally it's equivalent to:

```text
Input to sha256sum:

185f8db3...  ./file1.txt
4824c2f3...  ./file2.txt
a91bc4d1...  ./file3.txt

        │
        ▼
     SHA-256
        │
        ▼
Final Folder Hash
```

The important point is:

- `sha256sum file.txt` → hashes the **contents of a file**.
- `echo "Hello" | sha256sum` → hashes the **text "Hello"**.
- `find ... | sort | sha256sum` → hashes the **text output produced by the previous commands**. This works because when `sha256sum` is invoked without a filename, it reads from **standard input (stdin)**. ([man7.org](https://man7.org/linux/man-pages/man1/sha256sum.1.html?utm_source=chatgpt.com))

So your one-line understanding is:

> **We generate a text representation of the folder (file hashes + filenames), and then we calculate the SHA-256 hash of that text.**

That is the core logic behind a simple folder hash. Once you understand this, commands like:

```bash
echo "Hello" | sha256sum
printf "ABC\nDEF\n" | sha256sum
find . -type f -exec sha256sum {} + | sort | sha256sum
```

all follow the **same principle**: **whatever text (or bytes) comes through the pipe is what gets hashed.** ([man7.org](https://man7.org/linux/man-pages/man1/sha256sum.1.html?utm_source=chatgpt.com))