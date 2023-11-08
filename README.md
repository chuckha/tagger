# Tagger

Use tagger to read or modify id3 tags either one file at a time or in bulk.

Configuration can make the process easier.

## Supported versions

- `id3v2.3.0`

## Unsupported versions

- `id3v2.2.0`
- `id3v2.4.0`

## Configuration

- What to do when duplicate tags (according to the spec) are found?

Duplicate frame action is the action to take when a duplicate frame is encountered. A duplicate frame is one according to the id3v2 specification. Defaults to "do nothing".

Compress will minimize the padding of the id3 tag. This will make the filesize as small as possible, but it will cause an entire file-rewrite on tag edits in the future. Defaults to "false".

{
    "duplicate_frame_action": "keep first"|"do nothing"|"ask", (default: "do nothing")
    "compress": "true"|"false", (default: false)
}

## Tag editing

tagger --set-text-information --frame-id TRCK --value 1/2
tagger --set-attached-picture --description "abc" --picture @pic.jpg
tagger apic set --description "abc" --picture @pic.jpg
tagger trck set --value 2/3 my.mp3
tagger trck set --value 1/2 --

tagger fields trck get
tagger fields trck set 2/3
tagger fields apic set "description" @pic.jpeg
