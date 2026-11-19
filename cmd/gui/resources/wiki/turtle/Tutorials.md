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

每次前进后右转 144 度，五次回到起点附近。

```c
color "orange"
for range 5 {
    step 10
    right 144
}
```

## 后退

`step` 传负数即可后退，画笔保持落笔状态时会画出反向的线。

```c
step 8
step -8
```

## 彩色正方形

用循环改变颜色，画出一圈旋转的正方形。

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
