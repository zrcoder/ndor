package main

import "embed"

//go:embed resources/paint.png
var paintImgData []byte

//go:embed resources/teacher.png
var teacherImgData []byte

// //go:embed resources/gopher.png
// var gopherImgData []byte

//go:embed resources/code.png
var codeImgData []byte

//go:embed resources/teacher_labeled.png
var teacherLabeledImgData []byte

//go:embed resources/gopher_labeled.png
var gopherLabeledImgData []byte

//go:embed resources/wiki
var wikiFS embed.FS

//go:embed resources/icons/copy.svg
var copyIconData []byte

//go:embed resources/icons/check.svg
var checkIconData []byte

//go:embed resources/icons/save.svg
var saveIconData []byte

//go:embed resources/icons/close.svg
var closeIconData []byte

//go:embed resources/icons/sun.svg
var sunIconData []byte

//go:embed resources/icons/moon.svg
var moonIconData []byte
