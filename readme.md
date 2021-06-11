
##### 进程间stream通信(共享内存)

```
cd example
```

1. 写进程 
```
go run main.go -mode=read
```

2. 读进程
```
go run main.go -mode=read
```

x. 调试
```
go run main.go -mode=debug
```


##### Intel(R) Core(TM) i3-10100 CPU @ 3.60GHz
```
$ go run main.go -mode=write -size=128000
2021/06/11 12:02:45 main.go:33: mode: write size: 128000
2021/06/11 12:02:46 main.go:88: send:    1226.81 MB     1210.42 MB/s [per]    1210.42 MB/s [total]    1286400000 [total bytes]
2021/06/11 12:02:47 main.go:88: send:    2471.92 MB     1231.87 MB/s [per]    1221.13 MB/s [total]    2592000000 [total bytes]
2021/06/11 12:02:48 main.go:88: send:    3698.73 MB     1213.41 MB/s [per]    1218.56 MB/s [total]    3878400000 [total bytes]
2021/06/11 12:02:49 main.go:88: send:    4962.16 MB     1249.99 MB/s [per]    1226.41 MB/s [total]    5203200000 [total bytes]
2021/06/11 12:02:50 main.go:88: send:    6243.90 MB     1278.59 MB/s [per]    1236.77 MB/s [total]    6547200000 [total bytes]
2021/06/11 12:02:51 main.go:88: send:    7489.01 MB     1231.01 MB/s [per]    1235.81 MB/s [total]    7852800000 [total bytes]
2021/06/11 12:02:52 main.go:88: send:    8697.51 MB     1207.22 MB/s [per]    1231.76 MB/s [total]    9120000000 [total bytes]
2021/06/11 12:02:53 main.go:88: send:    9924.32 MB     1225.99 MB/s [per]    1231.04 MB/s [total]   10406400000 [total bytes]
2021/06/11 12:02:54 main.go:88: send:   11151.12 MB     1212.21 MB/s [per]    1228.94 MB/s [total]   11692800000 [total bytes]
2021/06/11 12:02:55 main.go:88: send:   12469.48 MB     1309.75 MB/s [per]    1237.01 MB/s [total]   13075200000 [total bytes]
2021/06/11 12:02:56 main.go:88: send:   13751.22 MB     1265.53 MB/s [per]    1239.61 MB/s [total]   14419200000 [total bytes]
2021/06/11 12:02:57 main.go:88: send:   14996.34 MB     1228.60 MB/s [per]    1238.69 MB/s [total]   15724800000 [total bytes]
2021/06/11 12:02:58 main.go:88: send:   16223.14 MB     1216.03 MB/s [per]    1236.95 MB/s [total]   17011200000 [total bytes]
2021/06/11 12:02:59 main.go:88: send:   17468.26 MB     1227.85 MB/s [per]    1236.30 MB/s [total]   18316800000 [total bytes]
2021/06/11 12:03:00 main.go:88: send:   18695.07 MB     1213.97 MB/s [per]    1234.81 MB/s [total]   19603200000 [total bytes]
```
```
$ go run main.go -mode=write
2021/06/11 12:06:36 main.go:33: mode: write size: 8192
2021/06/11 12:06:37 main.go:88: send:      79.69 MB       78.63 MB/s [per]      78.63 MB/s [total]      83558400 [total bytes]
2021/06/11 12:06:38 main.go:88: send:     158.20 MB       77.38 MB/s [per]      78.00 MB/s [total]     165888000 [total bytes]
2021/06/11 12:06:39 main.go:88: send:     234.38 MB       75.38 MB/s [per]      77.13 MB/s [total]     245760000 [total bytes]
2021/06/11 12:06:40 main.go:88: send:     310.55 MB       75.32 MB/s [per]      76.68 MB/s [total]     325632000 [total bytes]
2021/06/11 12:06:41 main.go:88: send:     386.72 MB       75.25 MB/s [per]      76.39 MB/s [total]     405504000 [total bytes]
2021/06/11 12:06:42 main.go:88: send:     462.89 MB       75.44 MB/s [per]      76.23 MB/s [total]     485376000 [total bytes]
2021/06/11 12:06:43 main.go:88: send:     539.06 MB       75.16 MB/s [per]      76.08 MB/s [total]     565248000 [total bytes]
2021/06/11 12:06:44 main.go:88: send:     615.23 MB       75.08 MB/s [per]      75.96 MB/s [total]     645120000 [total bytes]
2021/06/11 12:06:45 main.go:88: send:     691.41 MB       75.03 MB/s [per]      75.85 MB/s [total]     724992000 [total bytes]
2021/06/11 12:06:46 main.go:88: send:     767.58 MB       75.45 MB/s [per]      75.81 MB/s [total]     804864000 [total bytes]
2021/06/11 12:06:47 main.go:88: send:     843.75 MB       75.78 MB/s [per]      75.81 MB/s [total]     884736000 [total bytes]
2021/06/11 12:06:48 main.go:88: send:     919.92 MB       75.18 MB/s [per]      75.76 MB/s [total]     964608000 [total bytes]
```
