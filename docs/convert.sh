#!/bin/bash

FILE="$1"

ffmpeg -y -itsscale 0.4 -i ${FILE} -vf mpdecimate ${FILE}-interim.mp4

ffmpeg -i ${FILE}-interim.mp4 \
	-vf "fps=5,crop=700:650:0:0,scale=700:-1" \
	-c:v pam \
	-f image2pipe - | convert -delay 10 - -loop 0 -layers optimize ${FILE}.gif

