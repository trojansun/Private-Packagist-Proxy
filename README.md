# Private Packagist Proxy

加速Packagist的访问，资源本地化。

## 使用方法

### 配置

**所有的路径都不要带最后的斜杠**

### 编译

windows
```shell
.\scripts\build.bat
```

linux 
```shell
chmod +x ./scripts/build.sh
./scripts/build.sh
```

### 测试方法

#### 设置方法

##### 方法1-全局

设置本地的composer仓库地址

```shell
composer config -g repo.packagist composer http://127.0.0.1:8000
```

##### 方法2-项目

设置本地的composer仓库地址

```shell
composer config repo.packagist composer http://127.0.0.1:8000
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

### 开发计划

- [ ] 只缓存操作前最新的版本
- [ ] 缓存指定版本
- [ ] 定时更新
- [ ] 全量更新同步