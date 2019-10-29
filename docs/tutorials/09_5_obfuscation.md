# Obfuscation

## Built in GSCRIPT Obfuscation

GSCRIPT does two types of obfuscation:

- **Pre Compile Obfuscation** - This step is done _after_ GSCRIPT generates it's Golang intermediate representation, but _before_ it compiles a native binary with Go's compiler.
- **Post Compile Obfuscation** - This step is done _after_ GSCRIPT uses the Go compiler to build a native binary and operates atomically on the produced binary.

Each have strengths and weaknesses and knowing more about them will help you decide in what cases you will find them most useful.

**WARNING - Before reading about them, it's encouraged that you have an understanding of the GSCRIPT compiler and how it works at a high level.**

### Pre-Compiler Obfuscation

The **Pre Compile Obfuscator** is known as the **Stylist**. The basic premise of the Stylist is to remove string literals out of the intermediate representation Go source code.

Here is a hypothetical example of this:

```js
// Original Source
var foo = "bar"
```

And after the Stylist is run over this Go source:

```js
var foo = s(31337, generatedVarA)

var generatedVarA = []rune{
  0x0001, 0x0002, 0x0003, //... this is a rune slice containing encrypted bytes of the string
}

func s(key int, data []rune) string {
  // this function contains the logic to decrypt, decode, and translate the generated var into a string
}
```

Note that the Stylist only does this for Golang source code existing within the build directory and will not tangle strings outside of the GSCRIPT intermediate representation's `main` package.

If you want to see this in action, you can use the compiler flag `--keep-build-dir` and examine the differences between source with the Stylist turned on vs. off.

### Post-Compiler Obfuscation

**WARNING: There be dragons here. Proceed with caution.**

The **Post Compile Obfuscation** is known as the **Mordorifier**. The basic premise of the Mordorifier is that once GSCRIPT uses the Go compiler to build a native binary, regardless of platform, plaintext strings might exist which would be undesirable in a production build. This includes references to GSCRIPT, it's Engine, your `$HOME` directory or username, etc.

The **Mordorifier** follows the following basic steps:

1.  As the compiler runs, various pieces of it (including each VM, the Stylist, and the Compiler) retain a cache of "known strings".
2.  When the compiler initializes the Mordorifier, the cache of strings is supplied.
3.  The newly instantiated `Mordor` will enumerate this list of strings, creating 0-N number of regular expression matchers (known as `Orc` objects)
4.  Note if the regular expression cannot be built, the `Mordorifier` will discard string.
5.  Also `Orc` objects within the "`Horde`" (the map of `Orc` types within the `Mordor` object) are uniquely keyed. You will not get two `Orc` objects with the same string base.
6.  Besides the compiler provided cache, the `Mordorifier` also attempts to intelligently add `Orc` objects to it's `Horde` based on some presets we've defined within the code. This includes references to `gscript` as well as a number of other Golang specific strings we'd rather not be present in plain text.
7.  Finally, the compiler runs an **Assault** on the compiled executable - this enumerates each `Orc` within the `Horde`, performing a binary search and replace on any byte sequences which match the regular expression of the `Orc`. When a match is found, it is substituted simply with random data.

A quick glance at a fully obfuscated GSCRIPT binary in a disassembler will show you the magic of this. (Or the pain, depending on how you're looking at it).

Note that you cannot use the Mordorifier with either `--enable-logging` or `--enable-debugging` because you will end up with a _very_ messed up STDOUT buffer.

// NOTE: The **Mordorifier** has been known to cause some execution time problems. Debugging them is an incredibly difficult and slow process and if you experience binary failure during testing, it's always advised to simply disable the **Mordorifier** for the same build and check for any differences in execution.

### JavaScript Obfuscation

Currently, the compiler will attempt to "minify" your GSCRIPTs. While this is not "obfuscation" it does add another layer of stuff to reverse engineer. It is also written in a way where a javascript obfuscation engine could theoretically be easily bolted in.
