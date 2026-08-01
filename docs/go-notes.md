# Go Notes (from DockLens code)

Short notes on Go concepts learned while reading the DockLens CLI (`internal/cli`).

## Contents

- [`func init()`](#func-init)
- [Functions are values](#functions-are-values)
- [Structs can hold any type](#structs-can-hold-any-type)
- [Variable vs Field](#variable-vs-field)
- [Assignment syntax: `=`, `:=`, `var`, `Field:`](#assignment-syntax---var-field)
- [Reading files: types vs return types](#reading-files-types-vs-return-types)

Related doc: [Deep Dive: Docker](./deep-dive-docker.md)

## `func init()`

- Special Go function. Runs **automatically**, you never call it yourself.
- Order: package vars initialized → `init()` → `main()`.
- A package/file can have **multiple** `init()`s; imported packages init first.
- Takes no args, returns nothing.
- In Cobra: each command file uses `init()` to self-register with its parent, e.g. `imageCmd.AddCommand(analyzeCmd)`. Keeps `main()` tiny.

## Functions are values

- In Go, a function is a **first-class value** — assign it to a variable/field, pass it around.
- An **anonymous function** (function literal) has no name and is written inline:

```go
Run: func(cmd *cobra.Command, args []string) {
	imageName := args[0]
	analyzer.AnalyzeImage(imageName)
}
```

- `Run` is a **Cobra field** whose type is a function. The `func(...) {}` is a **Go** value matching Cobra's required signature. Cobra **calls** it later and supplies `cmd` + `args`.

## Structs can hold any type

- A struct groups named **fields**. A field's type can be almost anything: basic types, slices, maps, pointers, other structs, interfaces, channels, and **functions**.
- You store a **value of a type**, not "a type". The value's type must match the field's declared type (checked at compile time).
- Useful for callbacks/hooks (Cobra `Run`), swappable strategies, config with custom logic.
- Caveat: structs containing funcs/slices/maps are **not comparable** with `==`.

## Variable vs Field

- **Field** = a property that *belongs to a struct/thing*  — it models the entity (like columns of a DB row / labeled lines on a form). basically it is similar to vars but it attached to particular struct.
- **Variable** = a name *your code* uses to hold/refer to a thing while it works (scratch space for the process).
- Example: `Customer` has fields `Name`, `Email`, `Balance`. `alice := Customer{...}` — `alice` is a variable holding one customer.
- In Cobra: `Use`/`Short`/`Args`/`Run` are **fields** (shape of a command); `analyzeCmd` is a **variable** (your handle to one command).

> **Important — "field" is a Go term used *specifically for structs*.** Other composite types hold members too, but Go calls them by different names. So don't say "field of a slice/map"; use the right word below.


| Type        | Members are called              | Accessed by  | Named?             |
| ----------- | ------------------------------- | ------------ | ------------------ |
| **struct**  | **fields**                      | `value.Name` | yes (names)        |
| array/slice | **elements**                    | `value[0]`   | no (index)         |
| map         | **keys** & **values** (entries) | `value["k"]` | keys               |
| interface   | **methods** (method set)        | `value.Do()` | yes (method names) |


- A field is variable-like storage, but it is **owned by a struct value** — the struct *type* defines it once, every struct *value* gets its own copy (`alice.Balance` vs `bob.Balance` are separate).
- Outside Go the *word* "field" is general (DB columns, form fields, JSON properties, class fields in other languages) — but in Go code, reserve it for **structs**.

## Assignment syntax: `=`, `:=`, `var`, `Field:`


| Syntax       | Where             | What it does                          | Colon? |
| ------------ | ----------------- | ------------------------------------- | ------ |
| `b := X{}`   | inside a function | declare **new** variable + infer type | no     |
| `var b X`    | anywhere          | declare variable (zero value)         | no     |
| `a = X{}`    | anywhere          | assign to **existing** variable       | no     |
| `Field: val` | **inside** `X{ }` | set a **struct field** / map key      | yes    |


- `=` / `:=` / `var` operate on **variables** (name the "box").
- `Field:` only appears **inside a composite literal** `SomeStruct{ ... }` and sets **fields** (label the compartments inside the box).
- There is no standalone `var : value` statement — the colon form is always inside `{ }`.

```go
b := Button{        // ':=' declares variable b
	Label: "OK",    // 'Label:' sets a field
	Size:  10,
}
b = Button{Label: "Cancel", Size: 12}  // '=' reassigns existing b
```

## Reading files: types vs return types

A **struct/type has no return type**. Only **functions and methods** return values.

### Trail: `os.Open` → `*os.File` → read bytes

1. Start at [`os.Open`](https://pkg.go.dev/os#Open):

   ```go
   func Open(name string) (*File, error)
   ```

   Returns `*File` (handle) + `error` — **not** the file contents.

2. `File` is a type — [`os.File`](https://pkg.go.dev/os#File):

   ```go
   type File struct { ... }
   ```

   Go-to-definition here only shows the shape. Next: look at **methods** on `*File`.

3. In docs / Go-to-definition near `File`, you’ll see things like:

   | Method | Signature idea | Docs |
   | ------ | -------------- | ---- |
   | `Read` | `(n int, err error)` — bytes into a buffer, not whole file as text | [`(*File).Read`](https://pkg.go.dev/os#File.Read) |
   | `Close` | `error` — release the handle | [`(*File).Close`](https://pkg.go.dev/os#File.Close) |
   | `Name` | `string` — path used to open | [`(*File).Name`](https://pkg.go.dev/os#File.Name) |

   Full method list: [os.File methods](https://pkg.go.dev/os#File) (scroll to “Methods”).

4. `*os.File` implements [`io.Reader`](https://pkg.go.dev/io#Reader) (has `Read`). So helpers in package `io` that take a `Reader` work on it — e.g. [`io.ReadAll`](https://pkg.go.dev/io#ReadAll):

   ```go
   func ReadAll(r Reader) ([]byte, error)
   ```

5. Make bytes readable: `string(data)`, or parse with [`json.Unmarshal`](https://pkg.go.dev/encoding/json#Unmarshal) / [`json.NewDecoder`](https://pkg.go.dev/encoding/json#NewDecoder).

### Shortcut (no Open/Close yourself)

- [`os.ReadFile`](https://pkg.go.dev/os#ReadFile) → `([]byte, error)` in one call.

### How to find this without AI

1. What type do I have? → Go-to-definition / [pkg.go.dev](https://pkg.go.dev/).
2. Is it a `type`? → look for **methods**.
3. Does it look like `io.Reader`? → search package [`io`](https://pkg.go.dev/io).
4. Match function **signature** (inputs/outputs) to what you want.
5. Local CLI: `go doc os.Open`, `go doc os.File.Read`, `go doc io.ReadAll`.

