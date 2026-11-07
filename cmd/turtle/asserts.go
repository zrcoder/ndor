package main

import "embed"

//go:embed resources/teacher.png
var teacherImgData []byte

//go:embed resources/teacher_labeled.png
var teacherLabeledImgData []byte

//go:embed resources/gopher_labeled.png
var gopherLabeledImgData []byte

//go:embed resources/icons/close.svg
var closeIconData []byte

//go:embed resources/wiki
var wikiFS embed.FS
