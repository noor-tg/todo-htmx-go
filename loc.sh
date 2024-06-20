#!/bin/bash

wc -l $(fd -e go -E "**/*_test.go")
wc -l $(fd -e templ)
wc -l $(fd -e css)

