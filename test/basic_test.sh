#!/bin/bash

set -e
dir=$(mktemp -d)
cat > $dir/input <<DOC
im catbug! whats your name?
DOC

pass='spooky action at a distance'
ct='./cryptoad -v 6'

$ct -pass "$pass" $dir/input $dir/toadbug
pushd $dir
./toadbug -pass "$pass" -out decrypted.txt
diff input decrypted.txt

echo "blahblah" >> input
tar -cjvf input.tar.bz2 input
popd
$ct -pass "$pass" $dir/input.tar.bz2 $dir/toadbug2
pushd $dir

./toadbug2 -pass "$pass" -out dec.tar.bz2
mv input og_input
tar -xjvf dec.tar.bz2
diff input og_input
popd

rm -rf $dir/*
rmdir $dir
