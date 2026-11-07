# 海龟作图 (Turtle Graphics)

这是一个用来引导孩子编程兴趣的小应用。你可以用编码的形式来控制一只小海龟在画布上爬行并绘图。

## 使用说明

我们将在一个二维坐标系里用`编码`的方式来画海龟图。

`turtle`支持 [XGo](https://xgo.dev) 编程语言。

## Hello Turtle

```c
// 600×600 的画布
context 600, 600

// 设定当前颜色为红色，线宽为 4
color "red"
linew 4

// 抬笔移动到起点
penup
forward -150
right 90
forward -150
left 90
pendown

// 画一个正方形
for i <- :4 {
    forward 300
    right 90
}
```

> `"//"`开头的行或被`"/*"`与`"*/"`包围的块是注释，会被忽略，并不参与画图。
>
> `x` 坐标从左到右，`y` 坐标从下到上（画布中心为原点）。
