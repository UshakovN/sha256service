const mode = document.getElementById('mode')
selectorContentType.addEventListener('change', setMode)
let selectorMode

function setMode() {
    removeChildNodes(mode)
    if (selectorContentType.value.trim() !== '') {
        let content =
            '<label for="selectorMode" class="col-sm col-form-label">Режим работы️</label>' +
            '<div class="col-sm-10">' +
            '    <select class="form-select" id="selectorMode" name="selector">' +
            '        <option value="">Нажмите для выбора</option>' +
            '        <option value="create">Вычисление хеш-суммы</option>' +
            '        <option value="compare">Сравнение хеш-сумм</option>' +
            '    </select>' +
            '</div>'
        mode.insertAdjacentHTML('afterbegin', content)
        selectorMode = document.getElementById('selectorMode')
        selectorMode.addEventListener('change', setButton)
    }
}

