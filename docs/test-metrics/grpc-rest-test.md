# Test gRPC vs REST

## 1 Replica each - gRPC - valid Token

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books/a
Document Length:        23 bytes

Concurrency Level:      1000
Time taken for tests:   42.626 seconds
Complete requests:      50000
Failed requests:        0
Non-2xx responses:      50000
Total transferred:      10050000 bytes
HTML transferred:       1150000 bytes
Requests per second:    1172.99 [#/sec] (mean)
Time per request:       852.526 [ms] (mean)
Time per request:       0.853 [ms] (mean, across all concurrent requests)
Transfer rate:          230.24 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   5.8      1      57
Processing:     6  845 1591.9    359   28605
Waiting:        2  845 1591.9    359   28604
Total:          6  847 1592.1    360   28605

Percentage of the requests served within a certain time (ms)
  50%    360
  66%    604
  75%    868
  80%   1082
  90%   1893
  95%   3024
  98%   5640
  99%   8407
 100%  28605 (longest request)
```

## 1 Replica each - gRPC - invalid Token

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books/a
Document Length:        45 bytes

Concurrency Level:      1000
Time taken for tests:   30.626 seconds
Complete requests:      100000
Failed requests:        0
Non-2xx responses:      100000
Total transferred:      22400000 bytes
HTML transferred:       4500000 bytes
Requests per second:    3265.23 [#/sec] (mean)
Time per request:       306.257 [ms] (mean)
Time per request:       0.306 [ms] (mean, across all concurrent requests)
Transfer rate:          714.27 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   5.8      0    1002
Processing:     1  303 213.4    289     862
Waiting:        1  303 213.4    289     862
Total:          1  305 213.4    292    1337

Percentage of the requests served within a certain time (ms)
  50%    292
  66%    437
  75%    491
  80%    526
  90%    591
  95%    635
  98%    695
  99%    727
 100%   1337 (longest request)
```

## 2 Replica each - gRPC - valid Token

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books/a
Document Length:        23 bytes

Concurrency Level:      1000
Time taken for tests:   42.349 seconds
Complete requests:      50000
Failed requests:        0
Non-2xx responses:      50000
Total transferred:      10050000 bytes
HTML transferred:       1150000 bytes
Requests per second:    1180.67 [#/sec] (mean)
Time per request:       846.978 [ms] (mean)
Time per request:       0.847 [ms] (mean, across all concurrent requests)
Transfer rate:          231.75 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   8.2      1      78
Processing:     4  838 1677.9    329   35214
Waiting:        3  838 1677.9    329   35214
Total:          5  841 1678.2    331   35215

Percentage of the requests served within a certain time (ms)
  50%    331
  66%    601
  75%    872
  80%   1083
  90%   1894
  95%   2975
  98%   5416
  99%   8439
 100%  35215 (longest request)
```

## 2 Replica each - gRPC - invalid Token

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books/a
Document Length:        45 bytes

Concurrency Level:      1000
Time taken for tests:   26.874 seconds
Complete requests:      100000
Failed requests:        0
Non-2xx responses:      100000
Total transferred:      22400000 bytes
HTML transferred:       4500000 bytes
Requests per second:    3721.11 [#/sec] (mean)
Time per request:       268.737 [ms] (mean)
Time per request:       0.269 [ms] (mean, across all concurrent requests)
Transfer rate:          813.99 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    4   9.4      1    1024
Processing:     1  264 184.5    254     830
Waiting:        1  263 184.4    253     829
Total:          1  267 184.5    256    1564

Percentage of the requests served within a certain time (ms)
  50%    256
  66%    357
  75%    414
  80%    448
  90%    527
  95%    576
  98%    642
  99%    678
 100%   1564 (longest request)
```

## 1 Replica each - REST - valid Token

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books/a
Document Length:        23 bytes

Concurrency Level:      1000
Time taken for tests:   52.028 seconds
Complete requests:      50000
Failed requests:        0
Non-2xx responses:      50000
Total transferred:      10050000 bytes
HTML transferred:       1150000 bytes
Requests per second:    961.02 [#/sec] (mean)
Time per request:       1040.558 [ms] (mean)
Time per request:       1.041 [ms] (mean, across all concurrent requests)
Transfer rate:          188.64 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   9.9      1    1018
Processing:     7 1031 1900.0    457   30498
Waiting:        4 1030 1900.0    457   30498
Total:          8 1033 1900.2    466   30498

Percentage of the requests served within a certain time (ms)
  50%    466
  66%    789
  75%   1090
  80%   1330
  90%   2306
  95%   3584
  98%   6319
  99%   9679
 100%  30498 (longest request)
```

## 2 Replica each - REST - valid Token

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /api/v1/books/a
Document Length:        23 bytes

Concurrency Level:      1000
Time taken for tests:   33.868 seconds
Complete requests:      50000
Failed requests:        0
Non-2xx responses:      50000
Total transferred:      10050000 bytes
HTML transferred:       1150000 bytes
Requests per second:    1476.31 [#/sec] (mean)
Time per request:       677.363 [ms] (mean)
Time per request:       0.677 [ms] (mean, across all concurrent requests)
Transfer rate:          289.78 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    3   7.5      1      71
Processing:     4  667 1445.3    219   31031
Waiting:        3  666 1445.3    219   31031
Total:          5  670 1445.5    222   31032

Percentage of the requests served within a certain time (ms)
  50%    222
  66%    415
  75%    616
  80%    800
  90%   1534
  95%   2616
  98%   4918
  99%   7485
 100%  31032 (longest request)
```
