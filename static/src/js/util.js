import {Post} from "@/js/fetch.js";
import Cookies from "js-cookie";
import {snackbar} from "mdui";

export function Notification(data) {
    snackbar({
        message: data,
        placement: "bottom-end",
        action: "Copy",
        autoCloseDelay: 3000,
        onActionClick: () => navigator.clipboard.writeText(data)
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
                return session.repo_data || [];
            case "refreshRepo":
            case "delRepo":
                Notification("Successful");
                break;
        }
    }
}

const settingPartMap = {
    openai: "provider.openai",
    deepseek: "provider.deepseek",
    alibaba: "provider.alibaba",
    otherapi: "provider.otherapi",
    txt: "feature.text",
    img: "feature.image",
    web: "feature.web",
    rand: "feature.random",
    security: "security.dashboard",
    contxt: "prompt",
    taskBehavior: "security.task_behavior",
    txtSecurity: "security.text_prompt",
    imgSecurity: "security.image_prompt",
}

export async function Setting(operation, settingPart = "", settingEdit = null) {
    const usesNewSettingsAPI = operation === "fetchSettings" || operation === "editSettings"
    const actualPart = usesNewSettingsAPI ? (settingPartMap[settingPart] || settingPart) : settingPart
    const session = await Post({
        "Token": Cookies.get("token"),
        "Send": {
            "operation": operation,
            "setting_part": actualPart,
            "setting_edit": usesNewSettingsAPI ? null : settingEdit,
            "setting_body": usesNewSettingsAPI ? settingEdit : null,
        }
    })
    if (session.error) {
        Notification(session.error)
        return null
    } else {
        switch (operation) {
            case "editSetting":
            case "editSettings":
                Notification("Saved");
                break;
            case "fetchSetting":
                return session.setting_data;
            case "fetchSettings":
                return session.setting_body;
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
