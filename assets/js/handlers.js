const payloadTypeString = "string", payloadTypeEncodedFile = "encoded-file"


const compareStatusEqual = "Значения хеш-сумм совпадают"
const compareStatusNotEqual = "Значения хеш-сумм различны"

const compareStatusFailed = "Не удалось сравнить значения хеш-сумм"
const createStatusFailed = "Не удалось вычислить значение хеш-суммы"

const emptyClaimHashStatus = "Не задано контрольное значение хеш-суммы"
const emptyPayloadFile = "Не задан файл для хеширования"

const invalidMinimalSecretLen = "Минимальная длина парольной фразы должна быть положительным целым числом"
const invalidSecretLen = "Длина парольной фразы меньше допустимой"

const mandatoryCharNotFound = "В парольной фразе содержатся не все обязательные символы"
const uppercaseCharNotFound = "В парольной фразе отсутствуют символы верхнего регистра"
const lowercaseCharNotFound = "В парольной фразе отсутствуют символы нижнего регистра"
const digitsCharNotFound = "В парольной фразе отсутствуют цифры"



function reqCreateHash(secret, hashPayload, hashPayloadType) {
    return fetch("/create", {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(
            {
                "secret": secret,
                "payload": hashPayload,
                "payload_type": hashPayloadType,
            })
    }).then(resp => resp.json())
}

function reqCompareHash(claimHash, secret, hashPayload, hashPayloadType) {
    return fetch("/compare", {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(
            {
                "claim_hash": claimHash,
                "secret": secret,
                "payload": hashPayload,
                "payload_type": hashPayloadType,
            })
    }).then(resp => resp.json())
}

function handleCompareHash() {
    let output = document.getElementById('outputHashSum')
    let selector = document.querySelector('select')
    let input = document.getElementById('inputControl')
    let inputSecret = document.getElementById('inputSecret')

    let inputCompare = document.getElementById('compareInput')

    let hashPayload = input.value
    let hashPayloadType = payloadTypeString

    let secret = inputSecret.value
    if (!validateSecret(inputSecret, output)) {
        return
    }
    if (inputCompare.value === "") {
        output.innerText = emptyClaimHashStatus
        return
    }
    let claimHash = inputCompare.value

    switch (selector.value) {
        case "plain-text":
            let resp = reqCompareHash(claimHash, secret, hashPayload, hashPayloadType)
            setOutputCompareStatus(resp, output)
            break
        case "file":
            hashPayloadType = payloadTypeEncodedFile
            let reader = new FileReader()
            reader.onload = function () {
                hashPayload = reader.result.split(',')[1];
                let resp = reqCompareHash(claimHash, secret, hashPayload, hashPayloadType)
                setOutputCompareStatus(resp, output)
            }
            if (input.files.length === 0) {
                output.innerText = emptyPayloadFile
                return
            }
            reader.readAsDataURL(input.files[0])
            break
    }
}

function setOutputCompareStatus(promise, output) {
    promise.then(data => {
        let equal = data["equal"]
        if (equal === true) {
            output.innerText = compareStatusEqual
        }
        if (equal === false) {
            output.innerText = compareStatusNotEqual
        }
    }).catch(() => {
        output.innerText = compareStatusFailed
    })
}

function handleCreateHash() {
    let output = document.getElementById('outputHashSum')
    let selector = document.querySelector('select')
    let input = document.getElementById('inputControl')
    let inputSecret = document.getElementById('inputSecret')

    let hashPayload = input.value
    let hashPayloadType = payloadTypeString

    let secret = inputSecret.value
    if (!validateSecret(inputSecret, output)) {
        return
    }
    switch (selector.value) {
        case "plain-text":
            let resp = reqCreateHash(secret, hashPayload, hashPayloadType)
            setOutputHashSum(resp, output)
            break
        case "file":
            hashPayloadType = payloadTypeEncodedFile
            let reader = new FileReader()
            reader.onload = function () {
                hashPayload = reader.result.split(',')[1]
                let resp = reqCreateHash(secret, hashPayload, hashPayloadType)
                setOutputHashSum(resp, output)
            }
            if (input.files.length === 0) {
                output.innerText = emptyPayloadFile
                return
            }
            reader.readAsDataURL(input.files[0])
            break
    }
}

function setOutputHashSum(promise, output) {
    promise.then(data => {
        output.innerText = data["hash_sum"]
    }).catch(() => {
        output.innerText = createStatusFailed
    })
}

function validateSecret(secret, output) {
    let switchConstraints = document.getElementById('switchConstraints')
    if (!switchConstraints.checked) {
        return true
    }
    let input = document.getElementById('secretMinLen')
    let secretMinLen = input.value
    if (secretMinLen <= 0 || notIntegerValue(secretMinLen)) {
        output.innerText = invalidMinimalSecretLen
        return false
    }
    let secretText = secret.value
    console.log('Парольная фраза: ', secretText)
    if (secretText.length < secretMinLen) {
        output.innerText = invalidSecretLen
        return false
    }
    let inputMandatoryChars = document.getElementById('secretMandatoryChars')
    let secretMandatoryChars = inputMandatoryChars.value
    if (secretMandatoryChars !== "") {
        if (!checkContains(secretText, secretMandatoryChars)) {
            output.innerText = mandatoryCharNotFound
            return false
        }
    }
    let checkboxLowercase = document.getElementById('secretCheckboxLowercase')
    if (checkboxLowercase.checked) {
        if (!hasLowercase(secretText)) {
            output.innerText = lowercaseCharNotFound
            return false
        }
    }
    let checkboxUppercase = document.getElementById('secretCheckboxUppercase')
    if (checkboxUppercase.checked) {
        if (!hasUppercase(secretText)) {
            output.innerText = uppercaseCharNotFound
            return false
        }
    }
    let checkboxDigits = document.getElementById('secretCheckboxDigits')
    if (checkboxDigits.checked) {
        if (!hasDigits(secretText)) {
            output.innerText = digitsCharNotFound
            return false
        }
    }
    return true
}

function checkContains(secret, set) {
    for (let i = 0; i < set.length; i++) {
        let ch = set.charAt(i)
        if (!secret.includes(ch)) {
           return false
        }
    }
    return true
}

function hasUppercase(secret) {
    return /[A-ZА-Я]/.test(secret)
}

function hasLowercase(secret) {
    return /[a-zа-я]/.test(secret)
}

function hasDigits(secret) {
    return /[0-9]/.test(secret)
}

function notIntegerValue(len) {
    return /[^0-9]/.test(len)
}