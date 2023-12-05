split("\n")
| map(select(length > 0)
  | capture("Game (?<id>[0-9]+): (?<entries>.*)")
  | .id |= tonumber
  | .entries |= (.
    | split("; ")
    | map(.
      | split(", ")
      | map(capture("(?<number>[0-9]+) (?<color>.*)") | {"\(.color)":.number | tonumber})
      | reduce .[] as $x ({}; . + $x)
    )
  )
) 
| map(select(.entries | any(.red > 12 or .green > 13 or .blue > 14) | not))
| map(.id)
| reduce .[] as $val (0; . + $val)