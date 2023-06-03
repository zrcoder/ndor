# niudour

`niudour` is an app with wich you can draw by codes, for kids and coders~

![niudor](https://raw.githubusercontent.com/wiki/zrcoder/niudour/images/main.png)

- web app

    Vist <https://niudour.netlify.app> or run on your local machine as described bellow.

- cli app

    ```shell
    niudour xxx.gop
    ```

You can take a look at the [wiki](https://github.com/zrcoder/niudour/wiki) for detail.

## develop

`niudour` is impletioned with go, js and wasm.

If you want to run `niudour` on your local machine：

1. Install `git`, `Go` and [`task`](https://taskfile.dev)
2. Clone the repo

    ```shell
    git clone https://github.com/zrcoder/niudour
    ```

3. Run
    - web app:
        
        ```shell
        task run
        ```
        
        > delete or comment the line "Version: version" in main.go for every dev build
        > this is to break the cache mechanism, just for developing
        >
        > for product, modify the const "version" variable and then push and make pr
        >
        >  Open `http://localhost:9999`, start your painting journey～

    - cli app

        ```shell
        go run ./cmd/cli {source file}
        ```

