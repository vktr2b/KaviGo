# KaviGo
---
**KaviGo** is a simple Go-based CLI tool that automatically renames files to match the naming conventions required by [Kavita](https://github.com/Kareadita/Kavita), helping you keep your manga library better organized.

![Example](./misc/usage_example.svg)

Kavita is easily one of my favorite open-source projects, but one downside I encountered was the need to manually copy files to the correct directory on the server. Additionally, if I wanted to use Volumes, I had to manually rename hundreds of files.
So it seamed like a really good excuse to learn Go, I'm not much of a programmer (my favorite language is YAML) but I'm really proud of the end result even though the code is kind of a mess.

**KaviGo** renames manga chapters to a standardized naming convention. This ensures Kavita can properly parse the volume and chapter numbers, and even handle special chapters (see "Known Limitations" for more details).

I mostly get my manga from the [Haku Neko](https://github.com/manga-download/hakuneko) Mangasee123 connector. If you're using a different source with a different naming convention, this tool may not work as expected. However, if your manga chapters are organized in a directory with the manga's name and the chapter filenames include a number representing the chapter, it should work fine.


## Usage

Download the binary and run it using the following command:

```bash
./kavigo -d /path/to/source -o /path/to/destination -r /path/to/ranges.file -v -p
```

### Flags

| Flag | Description                                                                                             | Required |
| ---- | ------------------------------------------------------------------------------------------------------- | -------- |
| `-d` | Path to the source directory containing manga chapters (directory name should be the name of the manga) | Yes      |
| `-o` | Path to the destination directory. If not provided, files will remain in the source directory           | Optional |
| `-r` | Path to the Volume ranges file (comma-delimited). [[#Required volume ranges file\|Why?]] , example      | Yes      |
| `-v` | Verbose output                                                                                          | Optional |
| `-p` | Preserve original files                                                                                 | Optional |
| `-s` | Mark special chapters as special. See [[#Special Chapters]] for more info.                              | Optional |


## Limitations / Known Issues

### Required Volume Ranges File

Since I couldn't find an easily accessible API that provides the volume number for each chapter, I decided on a manual approach. You must provide a comma-delimited file that specifies the first and last chapters of each volume, along with the corresponding volume number. Use the `-r` (ranges) flag to specify the location of this file.

### Special Chapters

Kavita's handling of special chapters is a bit unconventional, or I might not fully understand how they should be used. My expectation was that if I named the file like `manga_name_v4_chp24.5_SP1.cbz`, it would be added to both the first volume and the "Special" tab, showing the volume and chapter number along with the special episode number. Unfortunately, this wasn't the case.

The workaround I implemented is to set both the chapter (`chp`) and special episode (`SP`) numbers to the same value. Special episodes will only be generated if you use the `-s` flag.