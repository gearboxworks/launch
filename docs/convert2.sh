#!/bin/bash

FILE="$1"

ffmpeg -i ${FILE} -vf mpdecimate ${FILE}-interim.mp4


ffmpeg -y -i ${FILE}-interim.mp4 \
	-vsync 0 \
	-r 5 \
	-filter_complex "[0:v] fifo,fps=5,crop=700:650:0:0,scale=w=700:h=-1,split [a][b];[a] palettegen [p];[b][p] paletteuse" \
	${FILE}.gif

# -filter_complex "[0:v] fifo,fps=12,crop=700:650:0:0,scale=700:-1,split [a][b];[a] palettegen [p];[b][p] paletteuse" \
# cycle=6,setpts=N/25/TB
exit

# Other options:

ffmpeg -i ${FILE} -vf "fps=10,scale=320:-1:flags=lanczos" -c:v pam -f image2pipe - | convert -delay 10 - -loop 0 -layers optimize ${FILE}.gif


