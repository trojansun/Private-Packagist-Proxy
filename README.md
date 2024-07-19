# Private Packagist Proxy

加速Packagist的访问，资源本地化。

## 使用方法

### 配置

**所有的路径都不要带最后的斜杠**

#### service

把`scripts/ppp.service`复制到`/lib/systemd/system/ppp.service`

```shell
# 启动
systemctl start ppp
# 停止
systemctl stop ppp
# 开机启动
systemctl enable ppp
```

### 编译

#### Windows

```shell
.\scripts\build.bat
```

#### Linux

```shell
chmod +x ./scripts/build.sh
./scripts/build.sh
```

### 运行

`ppp`是生成的可执行文件，根据你的系统选择对应的文件

```shell
chmod +x ppp
./ppp
```

### 测试方法

#### 设置方法

##### 方法1-全局

设置本地的composer仓库地址

```shell
composer config -g repo.packagist composer https://pkg.trojancdn.com
```

##### 方法2-项目

设置本地的composer仓库地址

```shell
composer config repo.packagist composer https://pkg.trojancdn.com
```

上面的设置完毕之后，随便找一个项目测试即可。
如使用`laravel`测试

```shell
composer create-project laravel/laravel example-app
```

如果你使用的是`http`协议，那么你还需要做如下设置

```shell
composer config -g secure-http false
```

#### 还原设置

```shell
composer config -g --unset repos.packagist
```

## 开发计划

- [ ] 只缓存操作前最新的版本
- [ ] 缓存指定版本
- [ ] 定时更新
- [ ] 全量更新同步
