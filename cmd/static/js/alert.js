class NdorAlert {
  constructor() {
    this.paintToast = null;
  }

  showHelp(version) {
    Swal.fire({
      imageUrl: "/static/images/teacher.png",
      imageHeight: 100,
      title: "Ndor v" + version,
      showCancelButton: true,
      confirmButtonText: "Document",
      cancelButtonText: "Examples",
      reverseButtons: true,
    }).then((result) => {
      if (result.isConfirmed) {
        window.location.href = "https://github.com/zrcoder/ndor/wiki";
      } else {
        document.getElementById("example-button").click();
      }
    });
  }

  alertError(number, errInfo) {
    if (number != -1) {
      errInfo = number + ": " + errInfo;
    }
    Swal.fire({
      icon: "warning",
      title: errInfo,
    }).then((result) => {
      if (number != -1) {
        MarkErrorLine(number);
      }
    });
  }

  alertEmptyInputWith(exampleCode) {
    Swal.fire({
      imageUrl: "/static/images/teacher.png",
      imageHeight: 100,
      showConfirmButton: false,
      titleText: "no code to run",
      text: "I'll give you an example soon",
      timer: 1999,
    }).then((reslut) => {
      SetCode(exampleCode);
    });
  }

  toastPainting() {
    if (this.paintToast == null) {
      this.paintToast = Swal.mixin({
        toast: true,
        position: "center",
        showConfirmButton: false,
        didOpen: () => {
          Swal.showLoading();
        },
      });
    }
    let img = "/static/images/paint.png";
    if (Math.random() < 0.5) {
      img = "/static/images/code.png";
    }
    this.paintToast.fire({
      imageUrl: img,
      imageHeight: 150,
    });
  }

  closePaintToast() {
    if (this.paintToast != null) {
      this.paintToast.close();
      this.paintToast = null;
    }
  }
}

const ndorAlert = new NdorAlert();

// for go

function getndorAlert() {
  return ndorAlert;
}
