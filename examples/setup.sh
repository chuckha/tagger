#!/usr/bin/env bash

# get the example mp3 files
curl --location --output aesop.zip https://www.archive.org/download/aesop_fables_volume_one_librivox/aesop_fables_volume_one_librivox_64kb_mp3.zip

# extract them
unzip aesop.zip
mkdir examples/output
# run tagger in dry-run mode

# tagger template-tag -template-config ./examples/template-config.json .
# run tagger in normal mode
