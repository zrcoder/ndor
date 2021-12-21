#! /usr/bin/env bash

# bash deploy.sh ../../niudour/static/

srcDir=$1

rsync -r ${srcDir} .

git add -A
git commit -m "deploy"

git push origin master