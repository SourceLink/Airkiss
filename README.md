# Airkiss

[Airkiss](https://iot.weixin.qq.com/wiki/doc/wifi/AirKissDoc.pdf)是微信提出的一种无线应用层协议,主要用于给无法交互的硬件设备进行网络配置, 比如小爱同学等.

## About:

这是一个用`go`语言编写的`SmartConfig App`的命令行版本, 这样就不需要在微信小程序上进行发送测试;

你现在可以根据你不同的环境使用go编译工具进行不同平台的版本编译;

- Linux
- Windows
- Arm Linux
- Mac Os


airkiss目录下是用于数据封包的实现, main.go文件中主要是调用airkiss进行数据封包, 再通过udp进行发送;

## Build

我的电脑环境是linux, 所以直接编译:

```go
go build
```



## Run

**PS**: 因为现在只是个简单的实现, 并没有指定网卡发送, 所以在发送前要断开你的有线网络;


查看帮助:

```
./airkiss_send -help
```

帮助:

```

Usage of ./airkiss_send:
  -e string
    	essid
  -p string
    	passwd
  -t int
    	timeout (default 1000)

```

测试:


```
./airkiss_send -e hardwareLab -p dsj88888 -t 10
```


测试结果图:

![](/image/test_result.png)

**PS:** 我是用的ESP8266做的设备配置.

