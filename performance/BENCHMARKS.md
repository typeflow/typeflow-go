# Benchmarks

A few notes on performance for this tiny package are shared here

Results
-------

[2 matrix rows approach](https://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_two_matrix_rows) is now used!

Operation             |  Times | ns/op  | B/op   | allocs/op |       CPU              | Revision
----------------------|--------|--------|--------|-----------|------------------------|--------------
 2 Rows Lev. Distance | 200000 |  7537 | 760   |    16     |  2.3 GHz Intel Core i5 |     master
 Find match (WordSource) | 2000 | 856892 | n/a | n/a |  2.3 GHz Intel Core i5  | master

Usage
------

Just issue the following command to run benchmarks on your machine.

```shell
go test -bench=".*" -benchmem
```

Contribute
----------

Please, help me keep this up-to-date and complete as well. If you happen to run the same benchmarks on a different machine, please, fill in the table and report your results.
