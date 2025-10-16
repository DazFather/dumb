# 🧠 DUMB

### *Delightfully Unintelligent Multilanguage Beautifier*

> Because sometimes all you need is a beautifier that *tries its best*.


**DUMB** is a formatting tool that pretends to understand your code, but really just wants to make it look **nice enough** to fool your coworkers.
It takes your files, guesses where the brackets go, panics a little, and then outputs something that *might* compile.

It’s not smart — it’s **DUMB**.

> If it messes up your indentation, just remember: you asked for it.

---

## ✨ Features (allegedly)

* Automatically indents your code _(most of the time)_
* Detects mismatched brackets and yells at you
* Works with *any* language (as long as it has brackets)
* Doesn't use ASTs, just vibes
* Supports custom line endings (`LF`, `CRLF`, or “whatever”)
* Can overwrite files or dump them somewhere else with `-o`
* Multithreaded, because chaos should be fast

---

## ⚙️ Installation

```bash
go install github.com/YourUsername/dumb@latest
```

## ⚙️ Usage

```bash
dumb [flags] [paths...]
```

### Usage example

Beautify everything in the current directory:

```bash
dumb
```

Use spaces instead of tabs (we forgive you):

```bash
dumb -s "    "
# shorthand for --spacer
```

Force Windows-style line endings (we don’t forgive you):

```bash
dumb --eol CRLF
```

Write output to a different folder:

```bash
dumb -o beautified/
# shorthand for --output
```

Print output to stdout:

```bash
dumb -echo
# equivalent of --output -

```

Remove joy for terminal minimalists:

```bash
dumb -nc
# shorthand for --no-color
```

---

## 🧩 Supported Languages

Yes.

*(DUMB doesn’t discriminate — if it has `{}`, `[]`, or `()`, it’ll try.)*

---

## 🧪 Example

Input:

```go
type block struct {
blocktype
Body[]string
}
```

Output:

```go
type block struct {
	blocktype
	Body []string
}
```

*(Wow. Stunning. Breathtaking. Absolutely gorgeous. 10/10. Would commit without reading )*

---

## 🐛 Bugs

Yes.

---

## ❤️ Contributing

Feel free to open issues or PRs.
Just don’t try to make DUMB *too* smart — that would defeat the purpose.


