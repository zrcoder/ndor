#! /usr/bin/env bash

if [ "$1" = "" ]; then
  echo 'usage:'
  echo 'bash deploy.sh ../niudour/static/*'
  exit 1
fi

src=$1

rsync -r ${src} .

git add -A
git commit -m "deploy"

git push origin master