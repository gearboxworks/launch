#!/bin/bash

go run MakeDMG.go \
	-assets assets/ \
	-bin launch \
	-icon assets/Gearbox.png \
	-identifier com.gearboxworks.launch \
	-name Launch \
	-dmg "Template.dmg" \
	-o out

