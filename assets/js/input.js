const selectorContentType = document.getElementById('selectorContentType')
const input = document.getElementById('input')

function setInput() {
    let type = selectorContentType.value
    let content
    switch (type) {
        case "plain-text":
            content =
                '<label for="textareaPlainText" class="col-sm col-form-label">' +
                    'Открытый текст' +
                '</label>' +
                '<div class="col-sm-10">' +
                    '<textarea class="form-control plaintext" id="inputControl" rows="1" name="input" placeholder="Введите текст...">' +
                    '</textarea>' +
                '</div>'
            break
        case "file":
            content =
                '<div class="mb-3">' +
                    '<input class="form-control" type="file" id="inputControl" name="input">' +
                '</div>'
            break
        default:
            removeChildNodes(input)
            return
    }
    removeChildNodes(input)
    input.insertAdjacentHTML('afterbegin', content)
}

selectorContentType.addEventListener('change', setInput)