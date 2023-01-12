# joplin2obsidian

## What is joplin2obsidian?
joplin2obsidian is a conversion tool to help you making the move from Joplin to Obsidian

## How to use
1. Firstly, Open Joplin, export all data (*notes and corresponding resources*) as RAW to a directory

![export](exportnotes.png)

2. Run joplin2obsidian to convert the "RAW - Joplin Export Directory" to Obsidian-Compatible vault directory 
~~~bash
Usage of joplin2obsidian:
  -s string
        Specify the source directory where Joplin exported the RAW data
  -d string
        The destination directory of Obsidian vault
~~~

3. Open the destination directory as vault in Obsidian

Done!

## Build from source
~~~bash
$ git clone https://github.com/luxi78/joplin2obsidian.git
$ cd  joplin2obsidian
$ make
$ cd dist
~~~

## Notes on additions (paulhagstrom)

I added something that would specifically on macOS (with XCode installed)
change the creation dates of the files, along with the modification date.
I also added a -nodates option to suppress the addition of creation date
and modification date (and Joplin ID) in the YAML frontmatter.
More interestingly, I added tag migration as well.  That should be platform-
agnostic.  I don't know how to write in go, so the status display as it converts
seems to be a bit glitchy but it does the conversion.  Allowed me to get from
Evernote to Obsidian via Joplin.  By the way: feel free to not persist this
comment if this code is ever adopted by a PR.  Maybe I'll do a PR someday if
nobody else does, but feel free to see how I did it or just lift the code,
I do not care.  If someone else finds it useful in doing a real PR against
the main repo, wonderful.