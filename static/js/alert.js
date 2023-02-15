function alertTeacherHelp() {
    Swal.fire({
        imageUrl: 'images/teacher.png',
        imageHeight: 100,
        title: 'Need help?',
        text: 'You can learn by reading the document.',
        confirmButtonText: 'Read',
        showCancelButton: true,
    }).then((result) => {
        if (result.isConfirmed) {
            window.location.href = 'https://github.com/zrcoder/niudour/wiki/%E7%89%9B%E8%B1%86%E5%84%BF%E7%94%BB%E5%9B%BE'
        }
    })
}

let paintToast = null

function getToast() {
    if (paintToast == null) {
        paintToast = Swal.mixin({
            toast: true,
            position: 'center',
            showConfirmButton: false,
            didOpen: () => {
                Swal.showLoading()
            },
        })
    }
    return paintToast
}

function toastPainting() {
    let img = 'images/paint.png'
    if (Math.random() < 0.5) {
        img = 'images/code.png'
    }
    getToast().fire({
        imageUrl: img,
        imageHeight: 150,
    })
}

function closePaintToast() {
    if (paintToast != null) {
        paintToast.close()
        paintToast = null
    }
}

function alertError(number, errInfo) {
    if (number != -1) {
        errInfo = number + ': ' + errInfo
    }
    Swal.fire({
        icon: 'warning',
        text: errInfo,
        showConfirmButton: false,
    }).then((result) => {
        if (number != -1) {
            markErrorLine(number)
        }
    })
}

function alertEmptyInputWith(exampleCode) {
    Swal.fire({
        imageUrl: 'images/teacher.png',
        imageHeight: 100,
        showConfirmButton: false,
        titleText: 'no code to run',
        text: "I'll give you an example soon",
        timer: 666
    }).then((reslut) => {
        setCode(exampleCode)
    })
}
