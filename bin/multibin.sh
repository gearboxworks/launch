#!/bin/bash

. ./bin/_Functions.sh

CMD="$(basename "$0" | sed 's/\.sh//')"

gb_${CMD} $@

EXIT="$?"

exit ${EXIT}
