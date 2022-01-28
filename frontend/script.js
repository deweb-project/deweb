function wrapper_openChat() {
    let id = document.getElementById("userid").value;
    getUser(id).then((v) => {
        console.log(v)
    })
}