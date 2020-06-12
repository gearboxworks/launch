#!/bin/bash

# for i in v1.3 v1.4 v1.4.1 v1.4.3 v1.6.0 v1.6.1 v1.6.4 v1.6.5 v1.6.6 v1.6.7
for i in v1.4 v1.4.1 v1.4.3 v1.6.0 v1.6.1 v1.6.4 v1.6.5 v1.6.6 v1.6.7
do
	rm dist/*
	buildtool sync $i
done

