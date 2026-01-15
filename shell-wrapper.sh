#!/bin/bash
# MacroFlow Shell Wrapper for Bash
# This wrapper enables cd macros to change your actual terminal directory
#
# Installation:
#   Add this to your ~/.bashrc:
#   source /usr/local/bin/macro-shell-wrapper.sh
#   (or paste the function directly into your ~/.bashrc)

macro() {
    # Capture output and execute the macro
    local output
    output=$(/usr/local/bin/macro "$@" 2>&1)
    local exit_code=$?
    
    # Check if this was a macro execution (not a management command)
    if echo "$output" | grep -q "^% "; then
        # Extract the command that was executed
        local cmd=$(echo "$output" | grep "^% " | sed 's/^% //')
        
        # Check if it's a cd command
        if [[ "$cmd" =~ ^cd[[:space:]]+(.*) ]]; then
            local target="${BASH_REMATCH[1]}"
            # Execute cd in the current shell
            echo "$output"
            builtin cd "$target" 2>/dev/null || builtin cd
            return $?
        fi
    fi
    
    # For non-cd commands or management commands, just show output
    echo "$output"
    return $exit_code
}
