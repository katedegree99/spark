#!/bin/sh
input=$(cat)

used=$(echo "$input" | jq -r '.context_window.used_percentage // empty')
remaining=$(echo "$input" | jq -r '.context_window.remaining_percentage // empty')
total_input=$(echo "$input" | jq -r '.context_window.total_input_tokens // empty')
window_size=$(echo "$input" | jq -r '.context_window.context_window_size // empty')
model=$(echo "$input" | jq -r '.model.display_name // empty')
cwd=$(echo "$input" | jq -r '.workspace.current_dir // empty')

dir=$(basename "$cwd")

if [ -n "$used" ] && [ -n "$remaining" ]; then
  used_int=$(printf '%.0f' "$used")
  remaining_int=$(printf '%.0f' "$remaining")
  if [ -n "$total_input" ] && [ -n "$window_size" ]; then
    printf '%s  %s  ctx: %s/%s tokens (%s%% used, %s%% left)' \
      "$model" "$dir" "$total_input" "$window_size" "$used_int" "$remaining_int"
  else
    printf '%s  %s  ctx: %s%% used / %s%% left' \
      "$model" "$dir" "$used_int" "$remaining_int"
  fi
else
  printf '%s  %s  ctx: --' "$model" "$dir"
fi
