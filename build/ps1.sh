# ------------------------------------------------------------------------------
export CLICOLOR=1
export LSCOLORS=GxFxCxDxBxegedabagaced
export PATH=$HOME/go/bin:$PATH
# ------------------------------------------------------------------------------
function parse_git_branch() {
	git branch 2>/dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/'
}
# ------------------------------------------------------------------------------
function get_gscript_version() {
	local init_cwd
	local gscript_version
	init_cwd=$(pwd)
	cd "$GOPATH/src/github.com/gen0cide/gscript"
	gscript_version="$(parse_git_branch)"
	cd "$init_cwd"
	echo "$gscript_version"
}
# ------------------------------------------------------------------------------
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
GSCRIPT_VERSION="$(get_gscript_version)"
# ------------------------------------------------------------------------------
PS1="$WHITE[$RED""gscript$NO_COLOUR/$YELLOW""docker $WHITE""version$NO_COLOUR:$CYAN\$GSCRIPT_VERSION$NO_COLOUR $YELLOW\w$WHITE]$NO_COLOUR\\$ "
# ------------------------------------------------------------------------------
