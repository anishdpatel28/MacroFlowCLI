#!/bin/zsh
# MacroFlow Shell Wrapper for Zsh
# This wrapper enables cd macros to change your actual terminal directory
#
# Installation:
#   Add this to your ~/.zshrc:
#   source /usr/local/bin/macro-shell-wrapper.zsh
#   (or paste the function directly into your ~/.zshrc)

macro() {
    # Check if this is a management command (not a macro execution)
    case "$1" in
        init|add|list|delete|delete-name|export|import|help|--help|-h|"")
            # Run management commands directly (preserves stdin/stdout/stderr)
            /usr/local/bin/macro "$@"
            return $?
            ;;
        *)
            # For macro execution, capture and eval
            local output
            output=$(/usr/local/bin/macro "$@" 2>&1)
            local exit_code=$?
            
            if echo "$output" | grep -q "^% "; then
                # Extract and execute the command
                local cmd=$(echo "$output" | grep "^% " | sed 's/^% //')
                echo "$output"
                eval "$cmd"
                return $?
            fi
            
            # If no % marker, just show output (error case)
            echo "$output"
            return $exit_code
            ;;
    esac
}
