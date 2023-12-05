split("\n") |
map(
  select(length > 0) |
  [scan("[0-9]")] |
  "\(first)\(last)" |
  tonumber
) |
reduce .[] as $val (0; . + $val)