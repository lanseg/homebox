function mkcd() {
    mkdir -p $@
    cd $@
}
function color_my_prompt {
    local __user_and_host="\[\033[01;32m\]\u@\h"
    local __cur_location="\[\033[01;34m\]\w"
    local __git_branch_color="\[\033[31m\]"
    local __git_branch='`git branch 2> /dev/null | grep -e "^*" | sed -E  s/^\\\\\*\ \(.+\)$/\(\\\\\1\)\ /`'
    local __prompt_tail="\[\033[35m\]$"
    local __last_color="\[\033[00m\]"
    export PS1="[\D{%T}] $__user_and_host $__cur_location $__git_branch_color$__git_branch\n$__prompt_tail$__last_color "
}
color_my_prompt

EDITOR=vim

alias ls='ls --color=auto --group-directories-first'
alias pwgen='pwgen -cnys '
alias note="printf '%(%Y-%m-%d %H:%M)T \t%s\t%s\n' -1 "

if [[ -n "$PS1" ]] && [[ -z "$TMUX" ]] && [[ -n "$SSH_CONNECTION" ]]; then
          tmux attach || tmux
fi

