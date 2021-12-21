#! /usr/bin/env bash

if [ "$@" = "" ]; then
  echo 'usage:'
  echo 'bash deploy.sh ../niudour/static/*'
  exit 1
fi

src=$@

rsync -r ${src} .

git add -A
git commit -m "deploy"

git push origin master