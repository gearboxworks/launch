#!/bin/bash

FILE="$1"

ffmpeg -i ${FILE} -vf mpdecimate ${FILE}-interim.mp4


ffmpeg -y -i ${FILE}-interim.mp4 \
	-vsync 0 \
	-r 5 \
	-filter_complex "[0:v] fifo,fps=5,scale=w=700:h=-1,split [a][b];[a] palettegen [p];[b][p] paletteuse" \
	${FILE}.gif

