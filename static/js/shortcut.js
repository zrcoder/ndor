document.addEventListener('keydown', function (event) {
    if (!event.ctrlKey || event.code != 'KeyR') {
        return
    }
    console.log("ctrl+r")
    document.getElementById('run-button').click()
})
