split("\n")
| map(select(length > 0)
  | capture("Game (?<id>[0-9]+): (?<entries>.*)")
  | .entries
  | split("; ")
  | map(.
    | split(", ")
    | map(capture("(?<number>[0-9]+) (?<color>.*)") | .number |= tonumber)
  )
  | flatten
  | group_by(.color)
  | map({"\(.[0].color)": max_by(.number).number})
  | reduce .[] as $x ({}; . + $x)
) 
| map(.red * .green * .blue)
| reduce .[] as $val (0; . + $val)