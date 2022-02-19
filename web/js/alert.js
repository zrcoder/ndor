function teacherAction(img) {
    Swal.fire({
        imageUrl: img,
        imageHeight: 100,
        title: 'Need help?',
        text: 'You can learn by reading the document.',
        confirmButtonText: 'Read',
        showCancelButton: true,

    }).then((result) => {
          if (result.isConfirmed) {
            window.location.href='https://gitee.com/rdor/niudour/blob/master/README.md'
          }
    })
}

function alertResult(imgSrc, msg, painterImg, teacherImg) {
    const Toast = Swal.mixin({
        toast: true,
        position: 'center',
        showConfirmButton: false,
        timer: 666,
        didOpen: () => {
            Swal.showLoading()
        },
    })
    Toast.fire({
        imageUrl: painterImg,
        imageHeight: 150,
    }).then((reslut) => {
        if(imgSrc !== '') {
            document.getElementById('pictureArea').src = imgSrc
            return
        }
        if(msg === 'no operations') {
            Swal.fire({
                imageUrl: teacherImg,
                imageHeight: 100,
                showConfirmButton: false,
                titleText: 'no code to run',
                text: "I'll give you an example soon",
            }).then((reslut) => {
                const exampleCode = 'context 800 800\ncolor 0 255 0 255\ncircle 400 400 300\nfill\n'
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

function alertError(msg, painterImg, teacherImg) {
    alertResult('', msg, painterImg, teacherImg)
}

function toastSuccess(src, painterImg, teacherImg) {
    alertResult(src, '', painterImg, teacherImg)
}

