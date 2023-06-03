document.addEventListener('keydown', function (event) {
    if (!event.ctrlKey || event.key != 'Enter') {
        return
    }
    console.log("ctrl+enter")
    document.getElementById('run-button').click()
})
