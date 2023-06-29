# ndor

`ndor` is an app with which you can draw by codes, for kids and coders~

![ndor](https://raw.githubusercontent.com/wiki/zrcoder/ndor/images/main.png)

You can take a look at the [wiki](https://github.com/zrcoder/ndor/wiki) for detail.

## web app

Vist <https://ndor.netlify.app> or run on your local machine as described bellow.

## cli app

Install with go:

```shell
go install github.com/zrcoder/ndor/cmd/ndor@latest
```

> TODO: Download from github releases

Run:

```shell
ndor {source file}
```

> The argument is the code source file, the cli will generate a image as xxx.png

## develop

`ndor` is impletioned with go, js and wasm.

If you want to run `ndor` on your local machine：

1. Install `git`, `Go` and [`task`](https://taskfile.dev)
2. Clone the repo

    ```shell
    git clone https://github.com/zrcoder/ndor
    ```

3. Run
    - web app:

        ```shell
        task run
        ```

        > delete or comment the line "Version: version" in cmd/web/main.go for every dev build
        > this is to break the cache mechanism, just for developing
        >
        > for product, modify the const "version" variable and then push and make pr
        >
        >  Open `http://localhost:9999`, start your painting journey～

    - cli app

        ```shell
        task cli -- examples/rainbow.gop
        ```
