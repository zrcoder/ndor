class NiudourAlert {
    constructor() {
        this.paintToast = null
    }

    showHelp(version) {
        Swal.fire({
            imageUrl: 'images/teacher.png',
            imageHeight: 100,
            title: 'Need help?',
            html: 'You can learn by reading the document.<br>' +
                'v' + version,
            confirmButtonText: 'Read',
            showCancelButton: true,
        }).then((result) => {
            if (result.isConfirmed) {
                window.location.href = 'https://github.com/zrcoder/niudour/wiki'
            }
        })
    }

    alertError(number, errInfo) {
        if (number != -1) {
            errInfo = number + ': ' + errInfo
        }
        Swal.fire({
            icon: 'warning',
            title: errInfo,
        }).then((result) => {
            if (number != -1) {
                MarkErrorLine(number)
            }
        })
    }

    alertEmptyInputWith(exampleCode) {
        Swal.fire({
            imageUrl: 'images/teacher.png',
            imageHeight: 100,
            showConfirmButton: false,
            titleText: 'no code to run',
            text: "I'll give you an example soon",
            timer: 1999
        }).then((reslut) => {
            SetCode(exampleCode)
        })
    }

    toastPainting() {
        if (this.paintToast == null) {
            this.paintToast = Swal.mixin({
                toast: true,
                position: 'center',
                showConfirmButton: false,
                didOpen: () => {
                    Swal.showLoading()
                },
            })
        }
        let img = 'images/paint.png'
        if (Math.random() < 0.5) {
            img = 'images/code.png'
        }
        this.paintToast.fire({
            imageUrl: img,
            imageHeight: 150,
        })
    }

    closePaintToast() {
        if (this.paintToast != null) {
            this.paintToast.close()
            this.paintToast = null
        }
    }
}

const niudourAlert = new NiudourAlert()

// for go
function getNiudourAlert() {
    return niudourAlert
}