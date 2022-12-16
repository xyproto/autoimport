#!/bin/sh -x
echo JAR 1
find /opt/hostedtoolcache/Java_Zulu_jdk/11.0.17-8/x64 -name "*.jar" -type f
echo JAR 2
find /opt/hostedtoolcache/Java_Zulu_jdk/ -name "*.jar" -type f
echo ENVIRONMENT
cat /etc/environment
echo JAVA READLINK
readlink -f java
echo JAVA WHICH
which java
echo JAVA VERSION
java -version
