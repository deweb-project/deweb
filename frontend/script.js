function wrapper_openChat() {
    let uid = document.getElementById("userid").value;
    getUser(uid).then((v) => {
        console.log(v)
    })
}