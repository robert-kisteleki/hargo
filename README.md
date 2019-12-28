# HAR extractor

This simple tool extracts files from saved HAR (HTTP Archive) files. It can
come handy if you want to save multiple components from a site like pictures
and don't want to do that manually.

You can create a HAR file with virtually any browser nowadays (Firefox, Chrome,
Safari, probably more) by inspecting the network traffic (right click on the
page, "inspect", go to network tab in the developer tools, reload your page to
start recording then right click on the list and select "export HAR" or "save
all as HAR" or selecting the HAR button and "save all as HAR" in Firefox).

Then you can provide the resulting HAR to this tool (`hargo < my.har`).
There are several options for your convenience:
* `-l` lists the URLs contained in the file
* `-i` and `-e` lets you filter which URLs should be extracted by providing
  simple strings that must appear in the URL ("include") or must not
  ("exclude")
* `-n` saves files with a simple sequence number intead of a name resembling
  the original URL 
* `-p` and `-s` provide a prefix and suffix that are added when constructing
  file names
* `-m` lets you filter for a method such as GET or POST
* `-v` turns on some logging 
* `-V` gives you version info 
* `-h` lists all the options


## Caveats

The tool is not built to be resilient against malformed input and will likely
blow up if you give it anything but a well formatted HAR. Use it at your own
risk.

It's most likely possible to change the code to use some kind of JSON streaming
thereby reducing the memory footprint if that becomes an issue.


## Why is this written in Go? It would be much simpler in $LANGUAGE!

Because it was a good exercise. Also, once compiled it can run in standalone
mode, no dependencies are needed. 
