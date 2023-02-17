# niudour

`niudour` is impletioned with go, js and wasm, you can draw with codes through it, for kids and coders~

You can take a look at the [wiki](https://github.com/zrcoder/niudour/wiki) for detail.

![niudor](https://raw.githubusercontent.com/wiki/zrcoder/niudour/images/main.png)

## develop

If you want to run `niudour` on your local machine：

1. Install `git` and `Go`
2. Clone the repo

    ```shell
    git clone https://github.com/zrcoder/niudour
    ```

3. Run

   ```shell
   cd niudour
   make local
   make run
   ```

   > delete or comment the line "Version: version" in main.go for every dev build
   > this is to break the cache mechanism, just for developing
   >
   > for product, modify the const "version" variable and then push and make pr

4. Open `http://localhost:9999`, start your painting journey～
