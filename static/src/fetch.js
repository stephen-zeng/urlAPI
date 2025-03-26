import {snackbar} from "mdui";
import {url} from "./main.js";

export function Post(data) {
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