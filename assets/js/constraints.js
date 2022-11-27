const selectConstraints = document.getElementById('switchConstraints')
const outConstraints = document.getElementById('constraints')

selectConstraints.addEventListener('change', function () {
    let content
    removeChildNodes(outConstraints)
    if (this.checked) {
        content =
            '<div class="input-group mb-3">\n' +
            '                    <input type="number" min="0" max="40" onKeyDown="return onlyNumberKey(event)"\n' +
            '                           class="form-control" placeholder="Минимальная длина" id="secretMinLen">\n' +
            '                    <input type="text" maxlength="20" class="form-control" placeholder="Обязательные символы" id="secretMandatoryChars">\n' +
            '                </div>\n' +
            '                <div class="input-group mb-3">\n' +
            '                    <div class="form-check">\n' +
            '                        <input class="form-check-input" type="checkbox" value="" id="secretCheckboxLowercase">\n' +
            '                        <label class="form-check-label" for="secretCheckboxLowercase">\n' +
            '                            Символы нижнего регистра\n' +
            '                        </label>\n' +
            '                    </div>\n' +
            '                    <pre>    </pre>\n' +
            '                    <div class="form-check tab-content">\n' +
            '                        <input class="form-check-input" type="checkbox" value="" id="secretCheckboxUppercase">\n' +
            '                        <label class="form-check-label" for="secretCheckboxUppercase">\n' +
            '                            Символы верхнего регистра\n' +
            '                        </label>\n' +
            '                    </div>\n' +
            '                    <pre>    </pre>\n' +
            '                    <div class="form-check tab-content">\n' +
            '                        <input class="form-check-input" type="checkbox" value="" id="secretCheckboxDigits">\n' +
            '                        <label class="form-check-label" for="secretCheckboxDigits">\n' +
            '                            Наличие цифр\n' +
            '                        </label>\n' +
            '                    </div>\n' +
            '                </div>'
        outConstraints.insertAdjacentHTML('afterbegin', content)
    }
})