#!/bin/sh
go build -mod=vendor
echo '--- without globs, for Example1.kt ---'
./autoimport -n -f Example1.kt
echo '--- without globs, for Example2.kt ---'
./autoimport -n -f Example2.kt
echo '--- verbose, glob imports are possible, for Example1.kt ---'
./autoimport -V -f Example1.kt
echo '--- verbose, glob imports are possible, for Example2.kt ---'
./autoimport -V -f Example2.kt
