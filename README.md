[toc]

# Goditor

> - A notebook and markdown editor written with go
> - markdown编辑器的go语言实现


## 安装
- 一些依赖：
  - windows:在win下面打开的是IE的渲染引擎，如果IE被取消使用的话需要重新启用
  - linux:默认使用的webviewgtk+,在fedora/centos下需要安装这些依赖项，其他发行版类似：`webkit2gtk3-devel、mesa-utils、libgl1-mesa-glx、glibc-devel`
- 下载：`go get -u github.com/blinderjay/Goditor`
- 运行：确保添加了go的环境变量`Goditor`

> 对开发者：如果您修改了app文件夹下的任何文件，必须先重新将app文件夹打包成为一个staik包，为了方便，您可以直接运行脚本
> - 进入安装目录：`cd $GOPATH/src/github.com/blinderjay/Goditor/buildscript`
> - 然后执行脚本：`./build.sh`或`./build.bat`

## 目标

功能：

- [ ] 实现基本的markdown功能
- [ ] 实现谷歌云盘同步
- [ ] html导出
- [ ] pdf导出
- [ ] png导出
- [ ] 多文件视图
- [ ] 界面美观


实现：

- [ ] 体积尽量精简
- [ ] 高性能
- [ ] 跨平台
- [ ] 使用有活力，通用的技术手段
- [ ] 方便扩展功能
- [ ] 尽量使用语言编程语言核心包，减少外来包
- [ ] 外来扩展包需要处于长期维护状态，或者自身有能力维护

## 技术规划

下面是一些分析

- markdown有两种实现手段
  - 双界面
    - 编辑页面提供语法高亮支持
    - 预览页面提供预览(一般是html)
  - 单界面
    - 一边编辑一般即时显示预览界面
    - (并不喜欢，个人感觉会严重影响写作效率)
- 所以最后的前端界面必然需要html显示:
  - 最简单的方法其实是完全采用前端的技术
    - 很多node模块可以处理这些
    - 基于前端打包本地应用的方法很多：
      - electron
      - nwjs
      - 单页web应用
    - web不方便提供本地的支持
    - 而另外的打包往往需要同时提供chromium和node的运行环境
      - 体积太大
      - 资源占用过多
- 考虑后续的扩展和一些其他的需求，后段必须使用go语言
- 那么麻烦的事情来了，前端怎么写
  - go ui库：最后证明使用纯正的 go UI库来写组建并不合适(不支持webview)
  - 用web来写
- 最后是前后端通信的问题

## 前端

- 其实最大的纠结在于html预览的部分，go本来就没有定位于客户端，一直都是用于服务端，ui库本来就及其少，而提供完善的 html+css+js 的更是极少。
- 最后我采用的方案是调用系统自带的web渲染引擎，省去了打包浏览器环境的体积和资源
  - 之所以可以调用底层的渲染引擎，要得益于go的 **cgo** 特性，绝大多数引擎都提供了c的接口
  - 单单调用web页面的渲染引擎不是什么难事，但同时又要实现跨平台，如果自己针对三大操作系统分别调试则会麻烦不少，感谢 [zserge/webview项目](https://github.com/zserge/webview) ，针对 golang/c/c++/py/ruby/haskal 等分别统一包装了webview接口
  - 目前绝大多数浏览器都是采用的 **webkit** web引擎，同时我主要的开发环境是fedora(gnome)和macos，分别有webkitgtk+ 和苹果正宗的webkit实现，主要调试环境还是在unix like环境

寻找ui库的过程中发现了不少优秀的项目，让我对计算机图形，计算机语言，计算机通信和web增加了不少的了解，所以我把那些年踩过的坑也拿出来总结一下
 
### fyne : 调用opengl底层编写 

glfw框架

### qt移殖

### gtk移殖

### web打包

#### gowd : 基于nwjs

#### go-astilectron : electron支持

### go-flutter : flutter桌面开发

### webview : 调用本地的"微型浏览器"

## 通信

### 数据绑定

### socket

### rpc

### websocket


## 感谢

感谢下面这些项目的支持，他们在**Goditor**中都起到了不可或缺的作用

- [russross/blackfriday](https://github.com/russross/blackfriday)：gomarkdown/markdown的fork源，提供了markdown最核心的功能，是golang中对markdown支持最好的一个
- [omarkdown/markdown](https://github.com/gomarkdown/markdown)：fork自russross/blackfriday，提供了go语言对markdown的解析和一些功能的扩展
- [sindresorhus/github-markdown-css](https://github.com/sindresorhus/github-markdown-css)：提供markdown的html显示
- [gorilla/websocket](https://github.com/gorilla/websocket#gorilla-websocket-compared-with-other-packages)：提供前后端通信的支持
- [rakyll/statik](https://github.com/rakyll/statik)：对静态资源打包，使goditor更像一个本地应用而非浏览器
- [zserge/webview](https://github.com/zserge/webview)：打开不同平台的默认浏览器渲染引擎，对golang提供操控窗口的接口
- [flowchart.min.js](https://github.com/adrai/flowchart.js)：提供流程图的绘制
- [fontawesome](https://fontawesome.com/how-to-use/on-the-web/setup/hosting-font-awesome-yourself)：一些图标的显示
- [editor.md](https://github.com/pandao/editor.md)：一份成功的参考，不过很多地方我不是很喜欢
