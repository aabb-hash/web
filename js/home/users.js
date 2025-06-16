const users = new Map()

function createUser(id, username) {
    users.set(id, username)

    if (window.location.pathname == "/") {
        const clone = circle.cloneNode(true);
        clone.id = "user-" + id;

        const avatar = clone.querySelector("#avatar");
        if (avatar) {
            avatar.src = "/api/avatar?id=" + encodeURIComponent(id);
        }

        document.body.appendChild(clone);
        return clone
    }
}

function deleteUser(id) {
    users.delete(id)

    document.getElementById("user-" + id).remove()
}