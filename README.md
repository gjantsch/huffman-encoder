# A Huffman Encoder

The result of the enconding is near optimal.

A improvement added was the canonical codes, to reduce the header. A lot of work to sava a hundred bytes on a large file but is ok.

There are room for improvement on the code and code structure, but will not do since is just a weekend work.

## How Canonical Huffman Codes Work

### Step 1 -- Build the Huffman tree from frequencies

```
Symbol  Freq
--------------
  ' '    15
  'e'    10
  't'     6
  'a'     4
  'x'     1
```

The Huffman tree (lower frequency = deeper = longer code):

```
           [36]
          /     \
       [21]     ' '(15)
      /    \
   [11]    'e'(10)
   /   \
 't'(6) [5]
        /  \
      'a'(4) 'x'(1)
```

Tree-traversal gives these codes (left=0, right=1):

```
Symbol  Freq  Tree code
-------------------------
  ' '    15   1          <- depth 1
  'e'    10   01         <- depth 2
  't'     6   000        <- depth 3
  'a'     4   0010       <- depth 4
  'x'     1   0011       <- depth 4
```

The actual bit patterns depend on which branch was traversed -- they are arbitrary.

### Step 2 -- Extract only the lengths

```
Symbol  Length
--------------
  ' '     1
  'e'     2
  't'     3
  'a'     4
  'x'     4
```

This is all the header needs to store. The bit patterns are thrown away.

### Step 3 -- Sort by (length, then symbol value)

```
Rank  Symbol  Length
---------------------
  1     ' '     1
  2     'e'     2
  3     't'     3
  4     'a'     4
  5     'x'     4
```

### Step 4 -- Assign canonical codes

Start with `code = 0`. At each step: if the length grew, shift `code` left by
the difference. Then format it as a zero-padded binary string of that length.

```
Symbol  Prev len  Cur len  Shift       code (binary)  Assigned
---------------------------------------------------------------
  ' '      0        1      0 << 1 = 0   "0"           ' ' = "0"    code++ -> 1
  'e'      1        2      1 << 1 = 2   "10"          'e' = "10"   code++ -> 3
  't'      2        3      3 << 1 = 6   "110"         't' = "110"  code++ -> 7
  'a'      3        4      7 << 1 = 14  "1110"        'a' = "1110" code++ -> 15
  'x'      4        4      15 << 0 = 15 "1111"        'x' = "1111" code++ -> 16
```

Final result:

```
Symbol  Tree code  Canonical code
-----------------------------------
  ' '      1           0
  'e'      01          10
  't'      000         110
  'a'      0010        1110
  'x'      0011        1111
```

The lengths are identical, so the compression ratio is the same. But the bit
patterns are now deterministic and reconstructible from lengths alone.

### Why the left-shift works

When moving from length N to length N+1, every existing code gets a 0
appended -- that is the left-shift by 1. This keeps the prefix-free property:
no code is a prefix of another.

```
0  -> ' '
1  -+- 10  -> 'e'
   +- 11 -+- 110  -> 't'
          +- 111 -+- 1110 -> 'a'
                  +- 1111 -> 'x'
```

### What the header stores

```
5 entries:
  [0x20] [1]    <- ' ', length 1
  [0x65] [2]    <- 'e', length 2
  [0x74] [3]    <- 't', length 3
  [0x61] [4]    <- 'a', length 4
  [0x78] [4]    <- 'x', length 4
```

10 bytes total instead of 10 bytes plus the actual code bytes. For 92 distinct
symbols (like in a typical English text) that saves roughly 92 bytes from the
header.
