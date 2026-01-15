#!/bin/zsh
# MacroFlow Shell Wrapper for Zsh
# This wrapper enables cd macros to change your actual terminal directory
#
# Installation:
#   Add this to your ~/.zshrc:
#   source /usr/local/bin/macro-shell-wrapper.zsh
#   (or paste the function directly into your ~/.zshrc)

macro() {
    # Run the macro binary and capture output
    local output
    output=$(/usr/local/bin/macro "$@" 2>&1)
    local exit_code=$?
    
    # Check if this was a macro execution (look for % command marker)
    if echo "$output" | grep -q "^% "; then
        # Extract the command to execute
        local cmd=$(echo "$output" | grep "^% " | sed 's/^% //')
        
        # Show what we're executing
        echo "$output"
        
        # Execute the command in the current shell
        eval "$cmd"
        return $?
    fi
    
    # For management commands (init, add, list, etc.), just show output
    echo "$output"
    return $exit_code
}
