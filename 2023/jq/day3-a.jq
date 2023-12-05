".\n\(.)\n."
| split("\n")
| map([match("[^.0-9]"; "g")] | map(.offset)) as $s
| map([match("[0-9]+"; "g")]) as $n
| keys | map({
    numbers: $n[.],
    symbols: ($s[. - 1]+$s[.]+$s[. + 1])
})
| map(
  foreach .numbers[] as $n (.symbols; .;
    select(any(. >= $n.offset - 1 and . <= $n.offset + $n.length))
    | $n.string | tonumber
  )
)
| reduce .[] as $val (0; . + $val)