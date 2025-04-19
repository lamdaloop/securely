const $  = (id)       => document.getElementById(id);
const $$ = (selector) => document.querySelectorAll(selector);

let currentSecretId = "";

const login = () => {
    window.location.href = "/auth/login";
};

const logout = () => {
    window.location.href = "/auth/logout";
};

const formatEmail = (email) => {
    return email.includes("@") ? email.split("@")[0] : email;
};

const checkAuth = async () => {
    try {
        const res = await fetch("/auth/me", {
            credentials: "include",
        });
        const data = await res.json();

        if (res.ok && data.email) {
            $("auth-user-text").innerText = `Logged in as ${formatEmail(data.email)}`;
            toggleAuthButtons(true);
        } else {
            showLoggedOut();
        }
    } catch {
        showLoggedOut();
    }
};

const toggleAuthButtons = (loggedIn) => {
    $("loginBtn")?.classList.toggle("hidden", loggedIn);
    $("logoutBtn")?.classList.toggle("hidden", !loggedIn);
};

const showLoggedOut = () => {
    $("auth-user-text").innerText = "Not logged in";
    toggleAuthButtons(false);
};

const parseUrlAndLoadSecret = () => {
    const match = window.location.pathname.match(/^\/secret\/(\w+)/);
    if (match) {
        currentSecretId = match[1];
        tryGetSecretWithoutPassword();
    }
};

const tryGetSecretWithoutPassword = async () => {
    const res = await fetch("/api/secret/" + currentSecretId, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({}),
    });

    const content = $("secretContent");
    const form    = $("passwordForm");
    const output  = $("retrieveResponse");

    try {
        const data = await res.json();

        if (res.ok) {
            content.innerHTML = `🔓 <b>Decrypted:</b><br>${data.message}<br><small>🧑‍💼 Created by: ${data.created_by}</small>`;
            form?.classList.add("hidden");
            output?.classList.add("hidden");
        } else if (res.status === 401) {
            console.log("🔐 Password required, showing form now...");
            form?.classList.remove("hidden");
            content.innerHTML = `🔐 This secret is password-protected.`;
        } else {
            content.innerHTML = `❌ ${data.error || "Something went wrong."}`;
        }
    } catch {
        content.innerHTML = "❌ Failed to load secret.";
    }
};

const submitPassword = async (event) => {
    event.preventDefault();

    const password = $("passwordInput").value;
    const content  = $("secretContent");
    const form     = $("passwordForm");
    const output   = $("retrieveResponse");

    content.innerHTML = "⏳ Trying to unlock...";

    const res = await fetch("/api/secret/" + currentSecretId, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ password }),
    });

    try {
        const data = await res.json();
        if (res.ok) {
            content.innerHTML = `🔓 <b>Decrypted:</b><br>${data.message}<br><small>🧑‍💼 Created by: ${data.created_by}</small>`;
            form?.classList.add("hidden");
            output?.classList.add("hidden");
        } else {
            output.classList.remove("hidden");
            output.innerText = "❌ Incorrect password.";
            content.innerHTML = "🔐 This secret is password-protected.";
        }
    } catch {
        content.innerHTML = "❌ Failed to unlock secret.";
    }
};

const setupExpiryButtons = () => {
    $$(".expiry-btn").forEach((btn) => {
        btn.addEventListener("click", () => {
            const minutes = btn.getAttribute("data-minutes");
            $("expiryInput").value = minutes;

            $$(".expiry-btn").forEach((b) => b.classList.remove("active"));
            btn.classList.add("active");
        });
    });
};

const submitSecret = async () => {
    const message  = $("secretInput")?.value.trim();
    const password = $("passwordInput")?.value;
    const expiry   = parseInt($("expiryInput")?.value);
    const oneTime  = $("oneTimeCheckbox")?.checked;
    const output   = $("createResponse");

    if (!message) {
        output.classList.remove("hidden");
        output.innerText = "❌ Please enter a secret.";
        return;
    }

    output.classList.remove("hidden");
    output.innerText = "⏳ Creating secret...";

    try {
        const res = await fetch("/api/secret", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({
                message,
                password,
                expire_in_minutes: expiry,
                one_time: oneTime,
            }),
        });

        const data = await res.json();
        if (res.ok) {
            const secretUrl  = `${window.location.origin}/secret/${data.id}`;
            output.innerHTML = `✅ Secret created!<br><a href="${secretUrl}" target="_blank">${secretUrl}</a>`;
        } else {
            output.innerText = `❌ ${data.error || "Something went wrong"}`;
        }
    } catch (err) {
        output.innerText = "❌ Failed to create secret.";
    }
};

document.addEventListener("DOMContentLoaded", () => {
    checkAuth();
    setupExpiryButtons?.();
    parseUrlAndLoadSecret();
});
