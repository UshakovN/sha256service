const button = document.getElementById('button')
selector.addEventListener('change', setButton)

function setButton() {
    removeChildNodes(button)
    if (selector.value.trim() !== '') {
        let content =
            '<button type="button" class="btn btn-outline-primary" onclick="handleCreateHash()">' +
                'Получить хеш-сумму' +
            '</button>' +
                '<span class="tab"> </span>' +
            '<button type="button" class="btn btn-outline-primary">' +
                'Сравнить хеш-суммы' +
            '</button>'
        button.insertAdjacentHTML('afterbegin', content)
    }
}