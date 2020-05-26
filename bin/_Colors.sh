#!/bin/sh

TERM=xterm-256color
export TERM

GB_TERM_HB="{{ .Json.bold.on }}{{ .Json.fg.cyan }}Gearbox[{{ .Json.fg.white }}"
GB_TERM_HE="{{ .Json.fg.cyan }}]: "

GB_TERM_CMD="${GB_BINFILE} -json ${GB_BINDIR}/_Colors.json -template-string"

p_echo() {
	local _name=${1:-$0}; shift
	local _fg=$1; shift
	local _txt="$@"

	${GB_TERM_CMD} "${GB_TERM_HB}${_name}${GB_TERM_HE}{{ .Json.fg.${_fg} }}${_txt}{{ .Json.reset }}\n"
}

p_n_echo() {
	local _name=${1:-$0}; shift
	local _fg=$1; shift
	local _txt="$@"

	${GB_TERM_CMD} "${GB_TERM_HB}${_name}${GB_TERM_HE}{{ .Json.fg.${_fg} }}${_txt}{{ .Json.reset }}"
}

bp_echo() {
	local _name=${1:-$0}; shift
	local _bg=$1; shift
	local _fg=$1; shift
	local _txt="$@"

	${GB_TERM_CMD} "${GB_TERM_HB}${_name}${GB_TERM_HE}{{ .Json.bg.${_bg} }}{{ .Json.fg.${_fg} }}${_txt}{{ .Json.reset }}\n"
}

bp_n_echo() {
	local _name=${1:-$0}; shift
	local _bg=$1; shift
	local _fg=$1; shift
	local _txt="$@"

	${GB_TERM_CMD} "${GB_TERM_HB}${_name}${GB_TERM_HE}{{ .Json.bg.${_bg} }}{{ .Json.fg.${_fg} }}${_txt}{{ .Json.reset }}"
}


c_echo()    { p_echo "" $@; }
c_n_echo()  { p_n_echo "" $@; }
bc_echo()   { bp_echo "" $@; }
bc_n_echo() { bp_n_echo "" $@; }

blinkon()  { ${GB_TERM_CMD} "{{ .Json.blink.slow }}"; }
blinkoff() { ${GB_TERM_CMD} "{{ .Json.blink.off }}"; }


p_err()    { local _name="$1"; shift; blinkon; bp_echo "${_name}" red white $@; blinkoff; }
p_warn()   { local _name="$1"; shift; p_echo "${_name}" yellow $@; }
p_info()   { local _name="$1"; shift; p_echo "${_name}" white $@; }
p_ok()     { local _name="$1"; shift; p_echo "${_name}" green $@; }

p_black()  { local _name="$1"; shift; p_echo "${_name}" black $@; }
p_red()    { local _name="$1"; shift; p_echo "${_name}" red $@; }
p_green()  { local _name="$1"; shift; p_echo "${_name}" green $@; }
p_yellow() { local _name="$1"; shift; p_echo "${_name}" yellow $@; }
p_blue()   { local _name="$1"; shift; p_echo "${_name}" blue $@; }
p_pink()   { local _name="$1"; shift; p_echo "${_name}" pink $@; }
p_cyan()   { local _name="$1"; shift; p_echo "${_name}" cyan $@; }
p_white()  { local _name="$1"; shift; p_echo "${_name}" white $@; }


c_err()    { blinkon; bp_echo "" red white "$@"; blinkoff; }
c_warn()   { p_echo "" yellow "$@"; }
c_info()   { p_echo "" white "$@"; }
c_ok()     { p_echo "" green "$@"; }

c_black()  { p_echo "" black "$@"; }
c_red()    { p_echo "" red "$@"; }
c_green()  { p_echo "" green "$@"; }
c_yellow() { p_echo "" yellow "$@"; }
c_blue()   { p_echo "" blue "$@"; }
c_pink()   { p_echo "" pink "$@"; }
c_cyan()   { p_echo "" cyan "$@"; }
c_white()  { p_echo "" white "$@"; }


checkExit()
{
	EXITCODE="$?"
	if [ "$EXITCODE" != "0" ]
	then
		c_err "Exited with $EXITCODE - \"$@\""
		exit $EXITCODE
	fi
}


# Usage
# c_echo $green "success!"
