const button = document.getElementById('button')

function setButton() {
    let type = selectorMode.value
    let content
    removeChildNodes(button)
    switch (type) {
        case "create":
            content =
                '<button type="button" class="btn btn-outline-primary" onclick="handleCreateHash()">' +
                    'Получить хеш-сумму' +
                '</button>'
            break
        case "compare":
            content =
                '    <div>' +
                '        <textarea type="text" class="form-control" id="compareInput" rows="1"' +
                '               placeholder="Введите хеш значение или загрузите его из файла..."></textarea>' +
                '    </div>' +
                '            <input class="form-control" type="file" id="compareFile" onChange="readFilePlainText()">' +
                '    <div>' +
                '        <button type="button" class="btn btn-outline-primary" onclick="handleCompareHash()">' +
                '            Сравнить хеш-суммы' +
                '        </button>' +
                '    </div>'
            break
        default:
            return
    }
    button.insertAdjacentHTML('afterbegin', content)
}
