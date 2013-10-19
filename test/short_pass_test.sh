#!/bin/bash

dir=$(mktemp -d)
cat > $dir/input <<DOC
a chicken run banana! cmon, lets go!
DOC

ct=./cryptoad
result=1

function finish {
	rm -rf $dir/*
	rmdir $dir
	
	exit $result
}

echo | $ct -pass "" $dir/input $dir/toadman
if [ $? -eq 0 -o -e $dir/toadman ]; then
	finish
fi

$ct -pass "a" $dir/input $dir/toadman
if [ $? -eq 0 -o -e $dir/toadman ]; then
	finish
fi

$ct -pass "abcdefg" $dir/input $dir/toadman
if [ $? -eq 0 -o -e $dir/toadman ]; then
	finish
fi

$ct -pass "abcdefgh" $dir/input $dir/toadman
if [ $? -ne 0 -o ! -e $dir/toadman ]; then
	finish
fi

result=0
finish
