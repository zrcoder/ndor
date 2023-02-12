function teacherAction() {
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

function alertResult(imgSrc, msg) {
    const Toast = Swal.mixin({
        toast: true,
        position: 'center',
        showConfirmButton: false,
        timer: 666,
        didOpen: () => {
            Swal.showLoading()
        },
    })
    let img = 'images/paint.png'
    if (Math.random() < 0.5) {
        img = 'images/code.png'
    }
    Toast.fire({
        imageUrl: img,
        imageHeight: 150,
    }).then((reslut) => {
        if (imgSrc !== '') {
            document.getElementById('pictureArea').src = imgSrc
            return
        }
        if (msg === 'empty input') {
            Swal.fire({
                imageUrl: 'images/teacher.png',
                imageHeight: 100,
                showConfirmButton: false,
                titleText: 'no code to run',
                text: "I'll give you an example soon",
            }).then((reslut) => {
                const exampleCode = 'context 800 800\ncolor 0 255 0 255\ncircle 400 400 300\nfill\n\n// click the gopher on bottom right to draw!'
                setCode(exampleCode)
            })
        } else {
            Swal.fire({
                icon: 'warning',
                text: msg,
                showConfirmButton: false,
            })
        }
    })
}

function alertError(msg) {
    alertResult('', msg)
}

function toastSuccess(src) {
    alertResult(src, '')
}

