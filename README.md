# trie
A trie data structure create for the sole purpose of easily reducing a large list of strings into a small list of strings.

More specifically, this library was created with only one use case in mind. A certain CDN only allows a small 
number of concurrent purge requests. If you have a thousands of purgeable files with frequent purges of various
combinations of these files things get a little more complicated. You don't want to purge the contents of the entire
folder every time you need to purge something but you also don't want to wait a long time for files to be purged. Fortunately, 
the CDN supports file globbering or the use of the '*' as a wildcard at the end of filename. So it is possible to purge these files:

```Bash
somefile.txt
someotherfile.txt
some.txt
```

with thise purge request:

```Bash
some*
```

This library makes it easy to generate a handful of purge strings from large lists of resources to be purged. Try it out by
running the test file or experimenting on the list of files given in `testset.txt` in this repo.
