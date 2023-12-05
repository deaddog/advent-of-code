"one|two|three|four|five|six|seven|eight|nine" as $numbers
| gsub("(?<num>" + $numbers + ")"; "\(.num)\(.num[-1:])")
| split("\n")
| map(
  select(length > 0) |
  [scan("[1-9]|"+$numbers)] | map(
    if . == "one" then 1
    elif . == "two" then 2
    elif . == "three" then 3
    elif . == "four" then 4
    elif . == "five" then 5
    elif . == "six" then 6
    elif . == "seven" then 7
    elif . == "eight" then 8
    elif . == "nine" then 9
    else . end  
  )
  | "\(first)\(last)"
  | tonumber
)
| reduce .[] as $val (0; . + $val)