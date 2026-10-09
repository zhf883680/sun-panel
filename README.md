[[ 简体中文 ]](https://doc.sun-panel.top/zh_cn/) |
[[ English ]](https://doc.sun-panel.top/)

<div align=center>

<img src="./doc/images/logo.png" width="100" height="100" />

# Sun-Panel

[![Github](https://img.shields.io/badge/Github-123456?logo=github&labelColor=242424)](https://github.com/zhf883680/sun-panel)
[![Gitee](https://img.shields.io/badge/Gitee-123456?logo=gitee&labelColor=c71d23)](https://gitee.com/hslr/sun-panel)
[![docker](https://img.shields.io/badge/docker-123456?logo=docker&logoColor=fff&labelColor=1c7aed)](https://hub.docker.com/r/zhf883680/sun-panel) 
[![Bilibili](https://img.shields.io/badge/Bilibili-123456?logo=bilibili&logoColor=fff&labelColor=fb7299)](https://space.bilibili.com/27407696/channel/collectiondetail?sid=2023810)
[![YouTube](https://img.shields.io/badge/YouTube-123456?logo=youtube&labelColor=ff0000)](https://www.youtube.com/channel/UCKwbFmKU25R602z6P2fgPYg)
<br>
[![GitHub User's stars](https://img.shields.io/github/stars/zhf883680%2Fsun-panel?style=flat&logo=github)](https://github.com/zhf883680/sun-panel)
[![github downloads](https://img.shields.io/github/downloads/zhf883680/sun-panel/total.svg?logo=github)](https://github.com/zhf883680/sun-panel/releases)
[![docker pulls](https://img.shields.io/docker/pulls/zhf883680/sun-panel.svg?logo=docker)](https://hub.docker.com/r/zhf883680/sun-panel)

[[ 中文文档 ]](https://sun-panel-doc.enianteam.com/zh_cn) |
[[ Document ]](https://sun-panel-doc.enianteam.com) |
[[ Demo ]](http://sunpaneldemo.enianteam.com) 

A server, NAS navigation panel, Homepage, Browser homepage.
<br>
一个服务器、NAS导航面板、Homepage、浏览器首页。

</div>

## 📌 关于本仓库（二次开发版）

本仓库是基于原版 [Sun-Panel](https://github.com/hslr-s/sun-panel) 二次开发的版本（原版自 `v1.8.1` 起闭源）。
**原版仓库：<https://github.com/hslr-s/sun-panel>** ｜ 原版文档：<https://doc.sun-panel.top/zh_cn/> ｜ 原版最后一个开源版本：`v1.3.0`

与原版对比：

- ✅ **同步的功能**（对齐原版 `v1.8.1`）
  图库（滚动分页、图标/壁纸分类）、公共图库、全局设置（站点设置 + 自定义 JS/CSS 在线编辑）、
  卡片「更多选项」（卡片背景色 / 分组 / 卡片类型 / 唯一标识）、关于页「检查新版本 + 运行信息」、
  多账号切换、搜索栏自定义搜索引擎、登录页自定义背景图、项目卡片缓存开关
- ➕ **本仓库新增**
  默认集成 [HD-Icons](https://github.com/xushier/HD-Icons) 在线图标库（后端代理转发，解决跨域与访问速度问题），
  图库弹窗提供「本地图库 / 在线图库」双 Tab；添加卡片时「获取图标」会顺带回填网站标题
- ➖ **本仓库移除**
  PRO 授权（功能全部免费、无任何限制）、Docker 管理、系统监控（CPU / 内存 / 磁盘挂载）、OpenAPI
- 🔧 **其他调整**
  配置路径统一到 `conf/` 下（`conf/uploads`、`conf/database`、`conf/custom`）；
  前端 `naive-ui` 锁定 `2.44.x`（`2.45+` 的 vapor 重写与 `Vue 3.2` 不兼容）

> 上游更新日志（含各版本功能变动）见 [UPDATELOG.md](./UPDATELOG.md) 与原版 Releases。

![](./doc/images/main-dark.png)

> [!IMPORTANT]
> In order to maintain the livelihood, the author added some [`PRO`] (https://pro.sun-panel.top) function, so the project temporarily entered a closed source state.; At present, the latest version of the open source is `v1.3.0`, [Please see the latest version of closed source](https://github.com/hslr-s/sun-panel/releases).; When the modular technology is developed, the separation of the PRO and the programs will be opened again, and the closed source will have no effect on ordinary users.; Let's look forward to open source again, and at the same time, we are welcome to supervise and review the security of the program.
> 
> 作者为了维持生计，增加了一些 [`PRO`](https://pro.sun-panel.top) 功能，所以项目暂时进入闭源状态。目前开源最新版本为`v1.3.0`，[闭源最新版本请查看](https://github.com/hslr-s/sun-panel/releases)。待开发出模块化技术，然后对PRO和主程序进行分离会再次开源，闭源对普通用户没有任何影响。我们一起期待再次开源吧，同时也欢迎各位大佬对程序的安全性进行监督和审查。

## 😎 Features

- 🍉 Clean interface, powerful functionality, low resource consumption
- 🍊 Easy to use, visual operation, zero-code usage
- 🍠 One-click switch between internal and external network modes
- 🍵 Supports Docker deployment (compatible with Arm systems)
- 🎪 Supports multi-account isolation
- 🎏 Supports viewing system status
- 🫙 Supports custom JS, CSS
- 🍻 Simple usage without the need to connect to an external database
- 🍾 Rich icon styles for free combination, supports [Iconify icon library](https://icon-sets.iconify.design/)
- 🚁 Supports opening small windows in the webpage (some third-party websites may block this feature)

## 🖼️ Preview Screenshots

**Various styles, freely combined**

![](./doc/images/icon-small-new.png)
![](./doc/images/transparent-info.png)
![](./doc/images/transparent-small.png)
![](./doc/images/solid-color-info.png)
![](./doc/images/full-color-small.jpg)

**Built-in small windows**

![](./doc/images/window-ssh.png)
![](./doc/images/window-xunlei.png)



## 🐳 Deployment tutorial
[Deployment Tutorial](https://sun-panel-doc.enianteam.com/usage/quick_deploy.html)

## 🍵 Donate

> Open-source development is not easy. If you feel that my project has helped you, you are welcome to [donate](./doc/donate.md) or buy me a cup of tea☕ (please leave your nickname or name in the note if possible). Your support is my motivation, thank you.


<a href="https://www.paypal.me/hslrs">
<img height="60" src="./doc/images/donate/paypal.png" target="_blank"></img> 
</a>


|   |   |
| ------------ | ------------ |
| <img height="300" src="./doc/images/donate/weixin.png"/> |  <img height="300" src="./doc/images/donate/alipay.png" /> |

## 🏖️ Communication group & community

Author：**[红烧猎人](https://blog.enianteam.com/u/sun/content/11)**

[Github Discussions](https://github.com/zhf883680/sun-panel/discussions)

QQ交流群，进不去可以点上方连接联系作者

<img src="./doc/images/qq_group_qr2.png"  height="350" />

## ❤️ Thanks

- [Roc](https://github.com/RocCheng)
- [jackloves111](https://github.com/jackloves111)
- [Rock.L](https://github.com/gitlyp)


---

[![Star History Chart](https://api.star-history.com/svg?repos=zhf883680/sun-panel&type=Date)](https://star-history.com/#zhf883680/sun-panel&Date)
