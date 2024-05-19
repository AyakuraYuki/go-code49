# Barcode Code49

Code49 barcode encode/decode implementation.

## Completed

* `Decode(bsLines []string, skipChecksum bool) string`

### Example

```text
11143121314115211131114321124131314
11221611211411251111225122311314214
11123232212411212332131231332321114
11251311211242114112215212413213114
11123121511212521211113243422213114
11224211311211313421211153141112154
```

```text
result: 
  - skip checksum: oveRWCON$tacTGbiNQARY6
  - accept checksum: oveRWCON$tacTGbiNQARY61OH2XQ
```

## Todo

* `Encode(text string) (bsLines []string)`
