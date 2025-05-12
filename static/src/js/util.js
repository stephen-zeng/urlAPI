import {Post} from "@/js/fetch.js";
import Cookies from "js-cookie";
import {snackbar} from "mdui";

export function Notification(data) {
    snackbar({
        message: data,
        placement: "top-end",
    })
}

export async function Login(token, term) {
    const session = await Post({
        "Token": token,
        "Send": {
            "operation": "login",
            "login_term": term,
        }
    })
    if (session.error) {
        Notification(session.error)
        return false
    } else {
        Notification("Login successful");
        Cookies.set("token", session.session_token, {expires: 7});
        return true
    }
}

export async function Logout(token) {
    const session = await Post({
        "Token": token,
        "Send": {
            "operation": "logout",
            "login_term": false,
        }
    })
    if (session.error) {
        Notification(session.error)
        return false
    } else {
        Notification("Logout successful");
        return true
    }
}

export async function Repo(operation, repoUUID = "", repoAPI = "", repoInfo = "") {
    const session = await Post({
        "Token": Cookies.get("token"),
        "Send": {
            "operation": operation,
            "repo_uuid": repoUUID,
            "repo_api": repoAPI,
            "repo_info": repoInfo,
        }
    })
    if (session.error) {
        Notification(session.error)
        return null
    } else {
        switch (operation) {
            case "fetchRepo":
                return session.repo_data;
                break;
            case "refreshRepo" || "delRepo":
                Notification("Successful");
                break;
        }
    }
}

export async function Setting(operation, settingPart = "", settingEdit = null) {
    const session = await Post({
        "Token": Cookies.get("token"),
        "Send": {
            "operation": operation,
            "setting_part": settingPart,
            "setting_edit": settingEdit,
        }
    })
    if (session.error) {
        Notification(session.error)
        return null
    } else {
        switch (operation) {
            case "editSetting":
                Notification("Saved");
                break;
            case "fetchSetting":
                return session.setting_data;
        }
    }
}

export async function Task(operation, taskCatagory, taskBy, taskPage = -1) {
    const session = await Post({
        "Token": Cookies.get("token"),
        "Send": {
            "operation": operation,
            "task_catagory": taskCatagory,
            "task_by": taskBy,
            "task_page": taskPage,
        }
    })
    if (session.error) {
        Notification(session.error)
        return null
    } else {
        return session
    }
}