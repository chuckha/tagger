# Tagger

Use tagger to read or modify id3 tags either one file at a time or in bulk.

Configuration can make the process easier.

## Supported versions

- `id3v2.3.0`

## Unsupported versions

- `id3v2.2.0`
- `id3v2.4.0`

## Configuration

A single file can be configured via a configuration file that looks like this:

```.json
{
    "Frames": {
        "TRCK": {"Information": "1/17"},
        "TEXT": {"Information": "J.K. Rowling"},
        "TPE1": {"Information": "Stephen Fry"},
        "TIT2": {"Information": "Chapter 01 - The Other Minister"},
        "TALB": {"Information": "Harry Potter and the Half-Blood Prince"}
    }
}
```

This configuration ensures that proper id3v2.3.0 specification is followed for all frames in the tag, as well as either modifying existing frames to match the configuration or else adding frames to the tag.

## Templated configuration

`tagger` offers a way to manage this configuration across a set of files. This is called templated configuration. A templated configuration looks like this:

```.json
{
    "FilePattern": "Harry Potter and the Half Blood Prince Disk $disk$/%reader% - Track $part$.mp3",
    "Overrides": {
        "reader": "Stephen Fry",
        "totalDiscs": 17
    },
    "OutputFilePattern": "hbp/{{.disk}}-{{.part}}.mp3",
    "FramesTemplate": "./templates/all/config.json.tmpl",
    "UserData": {
        "chapters": [
            "Chapter 01 - The Other Minister",
            "Chapter 02 - Spinner's End",
            "Chapter 03 - Will and Won't",
            "Chapter 04 - Horace Slughorn",
            "Chapter 05 - An Excess of Phlegm",
            "Chapter 06 - Draco's Detour",
            "Chapter 07 - The Slug Club",
            "Chapter 08 - Snape Victorious",
            "Chapter 09 - The Half-Blood Prince",
            "Chapter 10 - The House of Gaunt",
            "Chapter 11 - Hermione's Helping Hand",
            "Chapter 12 - Silver and Opals",
            "Chapter 13 - The Secret Riddle",
            "Chapter 14 - Felix Felicis",
            "Chapter 15 - The Unbreakable Vow",
            "Chapter 16 - A Very Frosty Christmas",
            "Chapter 17 - A Sluggish Memory",
            "Chapter 18 - Birthday Surprises",
            "Chapter 19 - Elf Tails",
            "Chapter 20 - Lord Voldemort's Request",
            "Chapter 21 - The Unknowable Room",
            "Chapter 22 - After the Burial",
            "Chapter 23 - Horcruxes",
            "Chapter 24 - Sectumsempra",
            "Chapter 25 - The Seer Overheard",
            "Chapter 26 - The Cave",
            "Chapter 27 - The Lightning-Struck Tower",
            "Chapter 28 - Flight of the Prince",
            "Chapter 29 - The Phoenix Lament",
            "Chapter 30 - The White Tomb"
        ]
    }
}
```
Below describes each section and what it does.

### `FilePattern`

`tagger` will recursively walk whatever directory is passed into it. If it encounters a file that matches this pattern, the generated configuration will be applied to it.

#### Special matchers

`tagger` provides several special matchers to be used within the FilePattern. These special patterns allow you to extract values from the file names and use them within the config template itself. For example, if your MP3 only has its track number in the file name and not in the frame data, this would allow you to extract the track number and put it in a TRCK frame.

`tagger` can extract digits or words from a title using this syntax: `$whatever$` for digits and `%something%` for words. The `whatever` and `something` are whatever you like and will be available to the template file.

The variable you choose (e.g. `whatever`, `something` in the sentence above) must be unique.

##### Example

If `template-config.json` declares a `FilePattern` of `hello_vol_$volume$_track $track$ (%author%).mp3`, then the template will have access to `{{.volume}}`, `{{.track}}` and `{{.author}}`.

### `Overrides`

Overrides are good for when you can extract something, but it might not be consistent spelling or it's maybe just wrong. These overrides allow you to override the extracted data from the file path. You can also add custom data here as well.

### `OutputFilePattern`

This is a go template that generates where the modified file will get written to. Note: this always requires a full file rewrite.

### `FramesTemplate`

A FramesTemplate looks liket this:

```
{
    "Frames": {
        "TPOS": {"Information": "{{.disk}}/{{.totalDiscs}}"},
        "TRCK": {"Information": "{{.special.count}}/{{.special.total}}"},
        "TCOM": {"Information": "Stephen Fry"},
        "TCON": {"Information": "Audiobook"},
        "TPE1": {"Information": "J. K. Rowling, Stephen Fry"},
        "TIT1": {"Information": "Harry Potter"},
        "TALB": {"Information": "Harry Potter and the Half-Blood Prince"}
    }
}
```

`tagger` will generate a config file for every mp3 it encounters. It has lots of data available to it. For example, anything that is extracted from the path name becomes available here. It can be added to a frame for future use.

#### Special template variables

`tagger` exposes a few special template variables to make data frame tagging more consistent. Here are a list of special variables exposed by `tagger` and what they do.

| name | template key | what |
| --- | --- | --- |
| count | `{{.special.count}}` | This is an ongoing count of every file processed regardless of where in the directory hierarchy it has been found |
| total | `{{.special.total}}` | This finds and counts all matching files before processing begins in order to keep a good consistent count and file order.

## Developing
