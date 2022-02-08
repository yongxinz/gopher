**原文链接：** [读 Go 源码，可以试试这个工具](https://mp.weixin.qq.com/s/E2TL_kcbVcRJ0CnxwbXWLw)

编程发展至今，从面向过程到面向对象，再到现在的面向框架。写代码变成了一件越来越容易的事情。

学习基础语法，看看框架文档，几天时间搞出一个小项目并不是一件很难的事情。

但时间长了就会发现，一直这样飘在表面是不行的，技术永远得不到提升。

想要技术水平有一个质的飞跃，有一个很好的方法，就是读源码。

但读源码真的是一件很有挑战的事情。

想想当年自己读 Django 源码，从启动流程开始看，没走几步就放弃了，而且还放弃了很多次。

这么说吧，我对 Django 启动部分的代码，就像对英文单词 abandon 那么熟悉。

后来总结经验，发现是方法不对。

主要原因是一上来就深入细节了，事无巨细，每个函数都不想错过。结果就导致对整体没有概念，抓不住重点，又深陷无关紧要的代码。最后就是看不进去，只能放弃。

最近看了一点 Go 源码，慢慢也摸索出了一些心得。有一个方法我觉得挺好，可以带着问题去读源码，比如：

- [Go Error 嵌套到底是怎么实现的？](https://mp.weixin.qq.com/s/nWb-0RTDG1Pg5ZmJZfbEPA)
- [为什么要避免在 Go 中使用 ioutil.ReadAll？](https://mp.weixin.qq.com/s/e2A3ME4vhOK2S3hLEJtPsw)
- [如何在 Go 中将 []byte 转换为 io.Reader？](https://mp.weixin.qq.com/s/nFkob92GOs6Gp75pxA5wCQ)

在解决问题的过程中也就对源码更熟悉了。

还有一点要注意的就是，先看整体，再看细节。

在这里推荐给大家一个工具，这个工具可以帮我们梳理出代码的整体结构，我觉得还是挺有用的。是一个开源项目：

**项目地址：** https://github.com/jfeliu007/goplantuml

这个项目可以分析一个 Go 项目，然后生成接口和结构体的 UML 图。有了这个图之后，基本上也就对项目整体关系有了一个基本概念，再读代码的话，相对来说会容易一些。

项目具体怎么用我倒是没仔细研究，因为老哥非常贴心的写了一个 WEB 页面：

**网站链接：** https://www.dumels.com/

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/dumels-1.png)

使用起来很方便，首先在页面最上方输入框输入项目地址，然后在左侧输入要分析的代码目录就可以了。默认生成的图中会包括 Fields 和 Methods。

填写好信息之后就可以生成 UML 图了。比如我输入的 `src/sync`，就得到了下面这张图，有了这张图，对代码结构之间的关系就更清晰了。

![](https://cdn.jsdelivr.net/gh/yongxinz/picb@main/data/dumels-2.png)

还可以一次分析多个目录，多个目录用英文逗号分割。

如果不填写要分析的目录，则会分析整个项目，也可以选择是否要忽略某个目录。

友情提示一点，不要试图分析整个 Go 项目，可能是项目太大了，页面是不会给你返回的。

好了，本文就到这里了。你有什么好用的工具吗？欢迎给我留言交流。

---

**往期推荐：**

- [开始读 Go 源码了](https://mp.weixin.qq.com/s/iPM-mPOepRuDqkBtcnG1ww)
- [推荐三个实用的 Go 开发工具](https://mp.weixin.qq.com/s/3GLMLhegB3wF5_62mpmePA)
