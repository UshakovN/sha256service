const emptyClaimHashFile = "Не задан файл содержащий контрольную хеш-сумму"

async function readFilePlainText() {
    let compareInput = document.getElementById('compareInput')
    let compareFile = document.getElementById('compareFile')
    let reader = new FileReader()
    reader.onload = function () {
        compareInput.value = reader.result
    }
    if (compareFile.files.length === 0) {
        compareInput.value = emptyPayloadFile
        return
    }
    reader.readAsText(compareFile.files[0])
}
