## 画正方形

```c
context 400, 400
color "blue"
linew 3

penup
forward -100
right 90
forward -100
left 90
pendown

for i <- :4 {
    forward 200
    right 90
}
```

## 画等边三角形

```c
context 400, 400
color "green"
linew 3

penup
forward -80
right 30
pendown

for i <- :3 {
    forward 160
    right 120
}
```

## 画螺旋线

```c
context 600, 600
color "purple"
linew 2

for i <- :100 {
    forward float(i) * 3
    right 91
}
```

## 抬笔移动

使用 `penup` 抬笔后，海龟移动时不绘图；使用 `pendown` 落笔后再绘图。`jump x, y` 可以直接瞬移到指定坐标而不留下痕迹。
