#! /usr/bin/env bash

if [ -z "$1" ]; then
  echo 'You should give the sources for this script, like:'
  echo './deploy.sh ../../../gitee/rdor/niudour/static/*'
  exit 1
fi

src=$@

rsync -r ${src} .

git add -A
git commit -m "deploy"
git push origin master