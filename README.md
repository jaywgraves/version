version
=======

_version_ is a command line tool for quickly saving copies of files with a different name.  This is a [Go] port of a similar tool I wrote in [Python] mostly as a exercise.  Instead of just appending a new suffix, it strips the extension and puts the new identifier before the file extension.  This way the file will be recognized as the same file type and can be opened by the right programs.

Many times when I'm writing scripts that have file output I want to save the previous versions of a file so I can test that any refactoring I do still produes the same output.  To verify the data I often I use a companion program [hashchk].
Usage
-----

    >version
    No file specifications given as arguments.
    Usage:  version -c|-m [-d -t] files... 

    This utility will copy all the files that are given as arguments and give each
    new file a name that has a datetime string (either current or file modification
    time) preceding the file extension.

    If neither the -m or the -c flag is set, then -m will be defaulted.
    Any combination of the -d and -t flags can be set. If neither are set,
    then both will be defaulted.

      -c    current: use current date in version string
      -d    date: add YYMMDD to version string
      -m    modification: use file modification time in version string
      -s    silent: suppress output
      -t    time: add HHMMSS to version string
      files...: 1 or more file name specifications


Example
-------

    >version sample.output 
    sample.output -> sample.20201002.084532.output

Version
----
1.0

License
----
MIT

[go]:http://golang.org/
[python]:http://python.org
[hashchk]:https://github.com/jaywgraves/hashchk