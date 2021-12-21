#! /usr/bin/env bash

if [ $1 == "" ]; then
  echo 'usage:'
  echo 'bash deploy.sh ../niudour/static/'
  exit 1
fi

srcDir=$1

rsync -r ${srcDir} .

git add -A
git commit -m "deploy"

git push origin master