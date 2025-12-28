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

EDIT5OR=vim
alias ytdlp='yt-dlp --no-overwrite -v --retries 3 --playlist-random --continue --write-description --write-info-json --write-subs --sub-langs all --output "%(upload_date)s %(channel)s - %(fulltitle)s.%(ext)s" '
alias ls='ls --color=auto --group-directories-first'
alias pwgen='pwgen -cnys '
alias note="printf '%(%Y-%m-%d %H:%M)T \t%s\t%s\n' -1 "
alias gl='git log --pretty="format:%C(blue)[%ad] %C(green)(%ar) %C(yellow)%h %C(reset)%s %C(dim white)- %an%C(reset)" --date="format:%Y-%m-%d %H:%M" '
alias shredr='shred -n 10 -z -u -v '

if [[ -n "$PS1" ]] && [[ -z "$TMUX" ]] && [[ -n "$SSH_CONNECTION" ]]; then
          tmux attach || tmux
fi

