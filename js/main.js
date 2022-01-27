// polyfill
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

function loadWasm(path) {
  const go = new Go()

  return new Promise((resolve, reject) => {
    WebAssembly.instantiateStreaming(fetch(path), go.importObject)
    .then(result => {
      go.run(result.instance)
      resolve(result.instance)
    })
    .catch(error => {
      reject(error)
    })
  })
}

loadWasm("niudour.wasm").then(wasm => {
  // console.log("wasm is loaded 👋")
})