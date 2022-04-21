# ------------------------------------------------------------------------------
export CLICOLOR=1
export LSCOLORS=GxFxCxDxBxegedabagaced
export PATH=$HOME/go/bin:$PATH
NO_COLOUR="\[\033[0m\]"
RED="\[\033[00;31m\]"
GREEN="\[\033[00;32m\]"
YELLOW="\[\033[00;33m\]"
BLUE="\[\033[00;34m\]"
MAGENTA="\[\033[00;35m\]"
PURPLE="\[\033[00;35m\]"
CYAN="\[\033[00;36m\]"
LIGHTGRAY="\[\033[00;37m\]"
LRED="\[\033[01;31m\]"
LGREEN="\[\033[01;32m\]"
LYELLOW="\[\033[01;33m\]"
LBLUE="\[\033[01;34m\]"
LMAGENTA="\[\033[01;35m\]"
LPURPLE="\[\033[01;35m\]"
LCYAN="\[\033[01;36m\]"
LWHITE="\[\033[01;37m\]"
WHITE="\e[97m"
GSCRIPT_VERSION="master"
# ------------------------------------------------------------------------------
PS1="$WHITE[$RED""gscript$NO_COLOUR/$YELLOW""docker $WHITE"" $YELLOW\w$WHITE]$NO_COLOUR\\$ "
# ------------------------------------------------------------------------------
