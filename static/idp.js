async function loadURLs(){
    let rawResponse = await fetch( "/urls")
    return await rawResponse.json()
}

async function send(event){
    // stop the form from submitting
    event.preventDefault()

    // get configured urls
    const urls = await loadURLs();

    // get the needed elements
    let form = document.getElementById("login-form")
    let identifier = form.elements["email"].value;
    let csrfToken = document.getElementsByName("gorilla.csrf.Token")[0].value

    // build the request data
    var urlParams = window.location.search
    let data = {
        identifier: identifier,
        login_challenge: urlParams.split("=")[1]
    }

    let url = "/preflight";
    let rawResponse = await fetch(url, {
        method: "POST",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            "X-CSRF-Token": csrfToken
        },
        body: JSON.stringify(data)
    })

    let content = await rawResponse.json();

    if (rawResponse.status !== 200) {
        form.elements["email"].classList.add("is-invalid")
        form.elements["password"].classList.add("is-invalid")

        // this message is vague (and yes possibly wrong) for security reasons
        document.getElementById("passwordHelp").innerHTML = "Username or password wrong"
        return
    }

    if (content.needs_redirect){
        // FEATURE: handle kratos flow init here
        form.action = urls.kratos_url;
    } else {
        form.action = "/"
    }

    form.submit()
}

document.getElementById("login-form").addEventListener("submit", send)