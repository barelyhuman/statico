#! /bin/bash 

file=commitlog-linux-amd64.tar.gz
curl -L https://github.com/barelyhuman/commitlog/releases/latest/download/$file --output $file
tar -xvf $file
./commitlog > CHANGELOG.txt
