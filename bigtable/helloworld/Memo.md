## Preparataion

- Created a service account `bitable-test@booming-opus-188014.iam.gserviceaccount.com`.
- Created a bigtable instance `south37-test`.

## Trial

```console
$ ./script/run
+ export PROJECT_ID=booming-opus-188014
+ PROJECT_ID=booming-opus-188014
+ export INSTANCE_ID=south37-test
+ INSTANCE_ID=south37-test
++ pwd
+ export GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ go run main.go -project booming-opus-188014 -instance south37-test
2020/10/24 22:40:25 Creating table Hello-Bigtable
2020/10/24 22:40:29 Writing greeting rows to table
2020/10/24 22:40:29 Getting a single greeting by row key:
2020/10/24 22:40:29     greeting0 = Hello World!
2020/10/24 22:40:29 Reading all greeting rows:
2020/10/24 22:40:30     greeting0 = Hello World!
2020/10/24 22:40:30     greeting1 = Hello Cloud Bigtable!
2020/10/24 22:40:30     greeting2 = Hello golang!
2020/10/24 22:40:30 Deleting the table
```

## After adding duration (in 1 multiplier)

```console
$ ./script/run
+ export PROJECT_ID=booming-opus-188014
+ PROJECT_ID=booming-opus-188014
+ export INSTANCE_ID=south37-test
+ INSTANCE_ID=south37-test
++ pwd
+ export GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ go run main.go -project booming-opus-188014 -instance south37-test
2020/10/24 23:24:25 Writing greeting rows to table
duration of applyBulk: 116.137215ms

2020/10/24 23:24:25 Getting a single greeting by row key:
duration of read all: 280.484608ms
duration of read all in average: 93ms

2020/10/24 23:24:25 Getting a single greeting by row key:
duration of read all: 64.279655ms
duration of read all in average: 21ms

2020/10/24 23:24:25 Reading all greeting rows:
duration of scan all: 18.025233ms
duration of scan all in average: 6ms
```

## After adding duration (in 10 multiplier)

```console
$ ./script/run
+ export PROJECT_ID=booming-opus-188014
+ PROJECT_ID=booming-opus-188014
+ export INSTANCE_ID=south37-test
+ INSTANCE_ID=south37-test
++ pwd
+ export GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ go run main.go -project booming-opus-188014 -instance south37-test
2020/10/24 23:23:22 Writing greeting rows to table
duration of applyBulk: 113.238113ms

2020/10/24 23:23:23 Getting a single greeting by row key:
duration of read all: 559.520785ms
duration of read all in average: 18ms

2020/10/24 23:23:23 Getting a single greeting by row key:
duration of read all: 219.645239ms
duration of read all in average: 7ms

2020/10/24 23:23:23 Reading all greeting rows:
duration of scan all: 15.608731ms
duration of scan all in average: 0s
```

## After adding duration (in 100 multiplier)

```console
$ ./script/run
+ export PROJECT_ID=booming-opus-188014
+ PROJECT_ID=booming-opus-188014
+ export INSTANCE_ID=south37-test
+ INSTANCE_ID=south37-test
++ pwd
+ export GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ go run main.go -project booming-opus-188014 -instance south37-test
2020/10/24 23:20:57 Writing greeting rows to table
duration of applyBulk: 113.682519ms

2020/10/24 23:20:57 Getting a single greeting by row key:
duration of read all: 2.622521056s
duration of read all in average: 8ms

2020/10/24 23:21:00 Getting a single greeting by row key:
duration of read all: 2.2326723s
duration of read all in average: 7ms

2020/10/24 23:21:02 Reading all greeting rows:
duration of scan all: 15.819643ms
duration of scan all in average: 0s
```

## After adding duration (100kB, 1 multiplier)

```console
$ ./script/run
+ export PROJECT_ID=booming-opus-188014
+ PROJECT_ID=booming-opus-188014
+ export INSTANCE_ID=south37-test
+ INSTANCE_ID=south37-test
++ pwd
+ export GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ go run main.go -project booming-opus-188014 -instance south37-test
2020/10/24 23:32:59 Writing greeting rows to table
duration of applyBulk: 141.090014ms

2020/10/24 23:32:59 Getting a single greeting by row key:
duration of read all: 554.310419ms
duration of read all in average: 184ms

2020/10/24 23:33:00 Getting a single greeting by row key:
duration of read all: 238.456533ms
duration of read all in average: 79ms

2020/10/24 23:33:00 Reading all greeting rows:
duration of scan all: 126.503222ms
duration of scan all in average: 42ms
```

## After adding duration (30GB, 100k multiplier)
GCP の BigTable 上では、storage size は 770MB になっていた。一方、program や network 転送量としては 30GB 送ってる。かなり圧縮されてるらしい。

ただ、1 key あたりの read の latency は 70-80ms に増えた。key が増えたことで binary search にかかる時間が増えたのかも。

```console
$ ./script/run
+ export PROJECT_ID=booming-opus-188014
+ PROJECT_ID=booming-opus-188014
+ export INSTANCE_ID=south37-test
+ INSTANCE_ID=south37-test
++ pwd
+ export GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ GOOGLE_APPLICATION_CREDENTIALS=/Users/minami/.go/src/github.com/GoogleCloudPlatform/golang-samples/bigtable/helloworld/bigtable-test-credential.json
+ go run main.go -project booming-opus-188014 -instance south37-test

2020/10/25 00:35:15 Writing greeting rows to table
duration of applyBulk: 18m32.679418706s

2020/10/25 00:53:47 Getting a single greeting by row key:
duration of read all: 7.310446608s
duration of read all in average: 73ms

2020/10/25 00:53:55 Getting a single greeting by row key:
duration of read all: 7.9617953s
duration of read all in average: 79ms

./script/run  120.22s user 282.94s system 35% cpu 18:52.14 total
```
