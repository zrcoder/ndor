#! /usr/bin/env bash

src='../../../gitee/rdor/niudour/static/*'

rsync -r ${src} .

git add -A
git commit -m "deploy"
git push origin master