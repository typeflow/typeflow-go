# Benchmarks

A few notes on performance for this tiny package are shared here

Results
-------

Currently a matrix-based implementation for the levenshtein computation is provided.
You can read more about it [here](https://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_full_matrix).

The branch **feature/2vecs** is there to support the [2 matrix rows approach](https://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_two_matrix_rows) that is expected to improve performance a lot :) .

  Times | ns/op  | B/op   | allocs/op |       CPU              | Needs Update 
--------|--------|--------|-----------|------------------------|--------------
 100000 |  12418 | 3776   |    54     |  2.5 GHz Intel Core i7 |     NO

Usage
------

Just issue the following command to run benchmarks on your machine.

```shell
go test -bench=".*" -benchmem
```

Contribute
----------

Please, help me keep this up-to-date and complete as well. If you happen to run the same benchmarks on a different machine, please, fill in the table and report your results.
