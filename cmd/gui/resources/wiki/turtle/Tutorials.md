## 正方形

重复四次：前进一步，右转 90 度。

```c
color "lightgreen"
width 3
for range 4 {
    step 10
    right 90
}
```

## 五角星

每次前进后右转 144 度，五次后恰好回到起点。

```c
color "orange"
for range 5 {
    step 10
    right 144
}
```

## 后退

`step` 传负数即可后退，落笔时会画出反向的线。先转向再后退，效果更明显：

```c
step 8
right 90
step -8
```

## 彩色正方形

用循环改变颜色，连续画三个首尾相接、每次转向 30° 的正方形（是一条折线，不是绕某点排成一圈）。

```c
func squar() {
    step 5
    right 90
    step 5
    right 90
    step 5
    right 90
    step 5
}

colors := ["red", "green", "blue"]

for c <- colors {
    color c
    squar
    right 30
}
```
