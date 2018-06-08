# 简介

这是一个使得MWeb支持直接上传文件至腾讯云的服务器。

![](https://media-1256569450.cos.ap-chengdu.myqcloud.com/blog/15284246994720.jpg)

# 背景

由于平时使用MWeb来写博客，发现使用腾讯云来上传文件没有相应的方式，并且官方提供的API文档似乎不能直接应用到MWeb上。

于是自己动手搞了一下这个专门用来上传图片至腾讯云的本地中转服务器。

# 使用方法

- 先按照 https://github.com/tencentyun/coscmd 安装好`coscmd`
- 通过`coscmd config -a <secret_id> -s <secret_key> -b <bucket> -r <region>` 写好配置文件 `~/.cos.conf`
- 检查配置文件`~/.cos.conf`是否已配置OK
- 克隆此项目至本地，然后启动`qcloud-cos-upload`

* 帮助文档

```txt
./qcloud-cos-upload -help
Usage of ./qcloud-cos-upload:
  -cosdir string
        COS目录 (default "/blog/")
  -help
        显示帮助
  -http string
        监听地址 (default "127.0.0.1:8016")
  -tmp string
        临时目录 (default "/tmp/qcloud-cos-tmpdir")
```

> 提示：若您使用了`anaconda3`安装多个Python环境，`run.sh`或许对你有帮助~

# 开发者

* 源码：有且只有一个文件`main.go`，一看便知~

* 编译：`go build -o qcloud-cos-upload .`

# 程序截图

![](https://media-1256569450.cos.ap-chengdu.myqcloud.com/blog/15284259864269.jpg)

# 后台运行？

由于刚刚从Node.js转入Golang，我这边直接使用`pm2`来管理后台程序，还蛮好用

* 启动：`pm2 start --name qcloud-cos-mweb ./run.sh`
* 停止：`pm2 stop qcloud-cos-mweb`
* 删除：`pm2 delete qcloud-cos-mweb`
* 状态：`pm2 status`

![](https://media-1256569450.cos.ap-chengdu.myqcloud.com/blog/15284262471222.jpg)


相比于Node.js开发的`SinaWeiboPictureBed`，可以清楚看到，Golang开发的程序似乎内存占用好小~