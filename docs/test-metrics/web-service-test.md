# Web-Service Test

## 1 Replica

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /
Document Length:        450 bytes

Concurrency Level:      1000
Time taken for tests:   25.076 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      65400000 bytes
HTML transferred:       45000000 bytes
Requests per second:    3987.82 [#/sec] (mean)
Time per request:       250.764 [ms] (mean)
Time per request:       0.251 [ms] (mean, across all concurrent requests)
Transfer rate:          2546.91 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    5  24.2      2    1035
Processing:     1  244  96.8    237    1297
Waiting:        1  243  96.8    236    1296
Total:          2  249  98.7    241    1310

Percentage of the requests served within a certain time (ms)
  50%    241
  66%    280
  75%    305
  80%    320
  90%    366
  95%    413
  98%    491
  99%    537
 100%   1310 (longest request)
```

## 2 Replicas

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /
Document Length:        450 bytes

Concurrency Level:      1000
Time taken for tests:   23.638 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      65400000 bytes
HTML transferred:       45000000 bytes
Requests per second:    4230.51 [#/sec] (mean)
Time per request:       236.378 [ms] (mean)
Time per request:       0.236 [ms] (mean, across all concurrent requests)
Transfer rate:          2701.91 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    8  47.4      3    1044
Processing:     1  227 138.8    213    1439
Waiting:        1  226 138.8    212    1439
Total:          1  235 146.6    220    1457

Percentage of the requests served within a certain time (ms)
  50%    220
  66%    300
  75%    341
  80%    367
  90%    423
  95%    470
  98%    515
  99%    563
 100%   1457 (longest request)
```

## 3 Replicas

```
Server Software:
Server Hostname:        vv.hsfl.de
Server Port:            32131

Document Path:          /
Document Length:        450 bytes

Concurrency Level:      1000
Time taken for tests:   25.083 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      65400000 bytes
HTML transferred:       45000000 bytes
Requests per second:    3986.74 [#/sec] (mean)
Time per request:       250.831 [ms] (mean)
Time per request:       0.251 [ms] (mean, across all concurrent requests)
Transfer rate:          2546.22 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    6  38.5      2    1036
Processing:     1  243 163.3    220     724
Waiting:        0  242 163.3    219     724
Total:          1  249 167.3    228    1513

Percentage of the requests served within a certain time (ms)
  50%    228
  66%    320
  75%    378
  80%    409
  90%    487
  95%    532
  98%    583
  99%    617
 100%   1513 (longest request)
```
