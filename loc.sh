#!/bin/bash

wc -l $(fd  "(go|css|templ)$" -E "**/*_test.go")
