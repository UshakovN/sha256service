const plainText = "plain-text", httpContent = "http-content", file = "file"



function handleCreateHash() {
    let selector = document.querySelector('select')
    let input = document.getElementById('inputControl')
    console.log(input.value)
    let data = input.value
    console.log(data)
    if (data !== null) {
        fetch("/create-hash", {
            method: 'POST',
            body: data,
        })
            .then(res => res.json())
            .then(data =>
                console.log(JSON.stringify(data))
            )
    }
}
