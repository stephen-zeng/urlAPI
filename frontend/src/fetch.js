import {snackbar} from "mdui";

export function Post(url, data) {
    return fetch(url, {
        method: "POST",
        body: JSON.stringify(data.Send),
        headers: {
            "Content-Type": "application/json",
            "Authorization": data.Token,
        }
    }).then(res => res.json())
}

export function Notification(data) {
    snackbar({
        message: data,
        placement: "top-end",
    })
}