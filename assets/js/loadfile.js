async function readFile(event) {
    let file = document.getElementById('compareFile').files[0]
    document.getElementById('compareInput').value = await file.text()
}