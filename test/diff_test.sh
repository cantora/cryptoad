#!/bin/bash

set -e
dir=$(mktemp -d)
cat > $dir/input <<DOC
the piper must be paid... in piperton pipe bucks.
DOC

pass="pipertonians are a proud people"
ct=./cryptoad

$ct -pass "$pass" $dir/input $dir/pipetoad
if [ ! -e $dir/pipetoad ]; then
	false #fail out (-e)
fi
if diff $dir/input $dir/pipetoad >/dev/null; then
	#they shouldnt be the same
	false #fail out
fi

if grep -i 'piperton pipe' $dir/pipetoad; then
	false #shouldnt find any message plain text in output
fi

result=0
