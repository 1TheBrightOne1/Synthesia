function updateCanvas() {
    let canvas = document.getElementById("canvas");

    let img = document.getElementById("frame");
    let ctx = canvas.getContext("2d");

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    ctx.drawImage(img, 0, 0);

    let keyCount = Number.parseInt(document.getElementById("keyCount").value);
    let spacing = canvas.width / keyCount;

    let keyList = document.getElementById("keys");
    if (keyList.value.length === 0) {
        let keys = [];
        for (let i = 0; i < keyCount; i++) {
            keys.push(Math.round(spacing * i));
        }
        keyList.value = keys.join(",");
    }

    for (let key of keyList.value.split(",")) {
        drawKeyBoundary(canvas, Number.parseInt(key));
    }

    document.getElementById("keyBorders").value = keyList.value;
}

function drawKeyBoundary(canvas, x) {
    let ctx = canvas.getContext("2d");
    ctx.beginPath();
    ctx.strokeStyle = "#fff";
    ctx.moveTo(x, 0);
    ctx.lineTo(x, canvas.height);
    ctx.stroke();
}