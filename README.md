# Barcode Code49

Code49 barcode encode/decode implementation.

> Reference: [Barcode-Symbology-Specification-Code-49](https://www.expresscorp.com/wp-content/uploads/2023/02/USS-49.pdf)

## Methods

### Encode

* `Encode(text string) (patterns []string, encodationPatterns [][]int, err error)`

Encode a given text into Barcode patterns in Code 49.

#### Example

* params

```text
overcontact binary
```

* returns

```text
patterns []string
[
    "11143121314115211131114321124131314",
    "11221611211411251111225122311314214",
    "11123232212411212332131231332321114",
    "11251311211242114112215212413213114",
    "11123121511212521211113243422213114",
    "11224211311211313421211153141112154"
]

encodationPatterns [][]int
[
    [ 1220, 1563, 730, 1355 ],
    [ 2168, 2180, 2179, 2195 ],
    [ 1465, 534, 632, 1437 ],
    [ 1906, 583, 926, 1153 ],
    [ 2166, 2183, 2190, 2358 ],
    [ 2400, 73, 835, 1643 ]
]
```

### DecodeRaw

* `DecodeRaw(patterns []string) string`

Decodes barcode `patterns` which contains multiple rows with scanned bar/space amounts, start with prefix `11` and end with suffix `4`, returns raw text.

#### Example

* params

```text
patterns []string
[
    "11143121314115211131114321124131314",
    "11221611211411251111225122311314214",
    "11123232212411212332131231332321114",
    "11251311211242114112215212413213114",
    "11123121511212521211113243422213114",
    "11224211311211313421211153141112154"
]
```

* returns

```text
overcontact binary
```

### Decode

* `Decode(patterns []string, skipChecksum bool) string`

Decodes barcode `patterns` which contains multiple rows with scanned bar/space amounts, start with prefix `11` and end with suffix `4`, returns basic text.
The basic text means that text decoded by this method contains checksum-mix-in characters.

It treads the Code49 as Mode-0, ignores other non-data characters.

If `skipChecksum` presents to true, Decode will ignore the last line of `patterns`.

#### Example

* params

```text
patterns []string
[
    "11143121314115211131114321124131314",
    "11221611211411251111225122311314214",
    "11123232212411212332131231332321114",
    "11251311211242114112215212413213114",
    "11123121511212521211113243422213114",
    "11224211311211313421211153141112154"
]
```

* returns

```text
result: 
  - skip checksum: oveRWCON$tacTGbiNQARY6
  - with checksum: oveRWCON$tacTGbiNQARY6\1OH2XQ
```
