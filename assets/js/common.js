const saveButton = document.getElementById('save')
saveButton.addEventListener('click', saveAsFile)

function clipboardContent() {
    let text = document.getElementById("outputHashSum")
    text.select()
    navigator.clipboard.writeText(text.value).then()
}

function removeChildNodes(el) {
    while (el.firstChild) {
        el.removeChild(el.firstChild)
    }
}

function saveAsFile() {
    let output = document.getElementById('outputHashSum')
    let data = output.value
    let ta = document.createElement("a");
    let taBlob = new Blob([data], {type: 'text/plain'});
    ta.setAttribute('href', URL.createObjectURL(taBlob));
    ta.setAttribute('download', 'hash.txt');
    ta.click();
    URL.revokeObjectURL(ta.href);
}

function onlyNumberKey(evt) {
    let ASCIICode = (evt.which) ? evt.which : evt.keyCode
    return !(ASCIICode > 31 && (ASCIICode < 48 || ASCIICode > 57));
}