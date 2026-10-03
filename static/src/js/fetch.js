import {url} from "../main.js";

// Errors the backend returns for missing, invalid or expired sessions.
const AUTH_ERRORS = ["Authentication failed", "Expired token"];

const sessionExpiredListeners = new Set();

// Registers a callback invoked when the server rejects the current session.
export function onSessionExpired(listener) {
    sessionExpiredListeners.add(listener);
    return () => sessionExpiredListeners.delete(listener);
}

function isAuthError(message) {
    return typeof message === "string" && AUTH_ERRORS.some((text) => message.includes(text));
}

// Post sends a session request and always resolves to an object. Failures
// (network errors, non-2xx responses, invalid JSON) resolve to
// {error: "..."} so callers never see a rejected promise.
export async function Post(data) {
    let res;
    try {
        res = await fetch(url, {
            method: "POST",
            body: JSON.stringify(data.Send),
            headers: {
                "Content-Type": "application/json",
                "Authorization": data.Token || "",
            }
        });
    } catch (err) {
        return {error: "Network error: " + (err && err.message ? err.message : "request failed")};
    }

    let body = null;
    try {
        const text = await res.text();
        body = text ? JSON.parse(text) : null;
    } catch {
        body = null;
    }

    if (body === null || typeof body !== "object" || Array.isArray(body)) {
        return {error: `Invalid response from server (HTTP ${res.status})`};
    }
    if (!res.ok && !body.error) {
        body.error = `Request failed (HTTP ${res.status})`;
    }
    if (body.error && data.Send && data.Send.operation !== "login" && isAuthError(body.error)) {
        sessionExpiredListeners.forEach((listener) => {
            try {
                listener(body.error);
            } catch (err) {
                console.error(err);
            }
        });
    }
    return body;
}
