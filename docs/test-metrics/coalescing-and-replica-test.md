# Test Coalescing and Replicas

## 1 Replica without Coalescing

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books
Document Length:        36896 bytes

Concurrency Level:      1000
Time taken for tests:   58.771 seconds
Complete requests:      50000
Failed requests:        0
Total transferred:      1850150000 bytes
HTML transferred:       1844800000 bytes
Requests per second:    850.76 [#/sec] (mean)
Time per request:       1175.424 [ms] (mean)
Time per request:       1.175 [ms] (mean, across all concurrent requests)
Transfer rate:          30742.72 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    3   9.9      2    1025
Processing:     8 1161 2720.2    407   49310
Waiting:        3 1152 2720.2    396   49303
Total:          9 1164 2720.3    410   49311

Percentage of the requests served within a certain time (ms)
  50%    410
  66%    736
  75%   1087
  80%   1371
  90%   2500
  95%   4166
  98%   8024
  99%  13862
 100%  49311 (longest request)
```

## 1 Replica with Coalescing

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books
Document Length:        36896 bytes

Concurrency Level:      1000
Time taken for tests:   30.334 seconds
Complete requests:      50000
Failed requests:        4
   (Connect: 0, Receive: 0, Length: 0, Exceptions: 4)
Total transferred:      1850150000 bytes
HTML transferred:       1844800000 bytes
Requests per second:    1648.34 [#/sec] (mean)
Time per request:       606.671 [ms] (mean)
Time per request:       0.607 [ms] (mean, across all concurrent requests)
Transfer rate:          59563.98 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   57 196.2     19    3068
Processing:    30  526 197.0    508    3315
Waiting:        2  423 159.7    416    3047
Total:        106  583 277.3    535    5358

Percentage of the requests served within a certain time (ms)
  50%    535
  66%    611
  75%    662
  80%    699
  90%    828
  95%   1061
  98%   1536
  99%   1647
 100%   5358 (longest request)
```

## 2 Replicas without Coalescing

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books
Document Length:        36896 bytes

Concurrency Level:      1000
Time taken for tests:   39.419 seconds
Complete requests:      50000
Failed requests:        0
Total transferred:      1850150000 bytes
HTML transferred:       1844800000 bytes
Requests per second:    1268.42 [#/sec] (mean)
Time per request:       788.384 [ms] (mean)
Time per request:       0.788 [ms] (mean, across all concurrent requests)
Transfer rate:          45835.18 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    6  16.8      3    1032
Processing:     5  771 1710.2    201   29335
Waiting:        3  757 1709.6    187   29330
Total:          6  776 1710.5    207   29337

Percentage of the requests served within a certain time (ms)
  50%    207
  66%    394
  75%    675
  80%    927
  90%   1927
  95%   3326
  98%   5794
  99%   8307
 100%  29337 (longest request)
```

## 2 Replicas with Coalescing

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books
Document Length:        36896 bytes

Concurrency Level:      1000
Time taken for tests:   26.582 seconds
Complete requests:      50000
Failed requests:        0
Total transferred:      1850150000 bytes
HTML transferred:       1844800000 bytes
Requests per second:    1880.95 [#/sec] (mean)
Time per request:       531.646 [ms] (mean)
Time per request:       0.532 [ms] (mean, across all concurrent requests)
Transfer rate:          67969.55 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   62 212.6     22    3070
Processing:     5  440 222.5    405    3925
Waiting:        2  356 195.1    336    3010
Total:          6  502 315.2    433    4262

Percentage of the requests served within a certain time (ms)
  50%    433
  66%    530
  75%    597
  80%    644
  90%    800
  95%   1010
  98%   1469
  99%   1605
 100%   4262 (longest request)
```

## 3 Replicas without Coalescing

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books
Document Length:        36896 bytes

Concurrency Level:      1000
Time taken for tests:   39.505 seconds
Complete requests:      50000
Failed requests:        0
Total transferred:      1850150000 bytes
HTML transferred:       1844800000 bytes
Requests per second:    1265.67 [#/sec] (mean)
Time per request:       790.093 [ms] (mean)
Time per request:       0.790 [ms] (mean, across all concurrent requests)
Transfer rate:          45736.08 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   39 134.6     20    3068
Processing:     9  743 520.0    655    7495
Waiting:        3  673 512.8    583    7494
Total:         13  782 538.3    685    7520

Percentage of the requests served within a certain time (ms)
  50%    685
  66%    840
  75%    949
  80%   1026
  90%   1265
  95%   1597
  98%   2247
  99%   3050
 100%   7520 (longest request)
```

## 3 Replicas with Coalescing

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books
Document Length:        36896 bytes

Concurrency Level:      1000
Time taken for tests:   30.748 seconds
Complete requests:      50000
Failed requests:        0
Total transferred:      1850150000 bytes
HTML transferred:       1844800000 bytes
Requests per second:    1626.12 [#/sec] (mean)
Time per request:       614.961 [ms] (mean)
Time per request:       0.615 [ms] (mean, across all concurrent requests)
Transfer rate:          58761.02 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0  118 338.1     33    7116
Processing:    35  451 292.3    399   13466
Waiting:        2  298 205.8    267    6990
Total:        106  570 456.1    448   14523

Percentage of the requests served within a certain time (ms)
  50%    448
  66%    556
  75%    624
  80%    679
  90%   1024
  95%   1429
  98%   1705
  99%   2288
 100%  14523 (longest request)
```

## 10 Replicas with Coalescing

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books
Document Length:        36896 bytes

Concurrency Level:      1000
Time taken for tests:   39.548 seconds
Complete requests:      50000
Failed requests:        0
Total transferred:      1850150000 bytes
HTML transferred:       1844800000 bytes
Requests per second:    1264.29 [#/sec] (mean)
Time per request:       790.959 [ms] (mean)
Time per request:       0.791 [ms] (mean, across all concurrent requests)
Transfer rate:          45685.99 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   97 285.2     30    7304
Processing:    25  545 360.5    493   13610
Waiting:        2  401 302.7    378   12900
Total:         93  642 469.5    538   14675

Percentage of the requests served within a certain time (ms)
  50%    538
  66%    664
  75%    754
  80%    828
  90%   1042
  95%   1478
  98%   1790
  99%   2058
 100%  14675 (longest request)
```
