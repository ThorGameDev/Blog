import { err } from "./lib/utils.ts"

const password: HTMLInputElement | null = document.getElementById("password") as HTMLInputElement | null;
const confirmPass: HTMLInputElement | null = document.getElementById("confirmPass") as HTMLInputElement | null;
const form: HTMLFormElement | null = document.getElementById("signupForm") as HTMLFormElement | null;

async function sendSignup(): Promise<void> {
    if (!confirmPass || !password || !form) {
        err("Missing a critical element!");
    }

    if (confirmPass.value !== password.value) {
        alert("Passwords do not match!");
        return;
    }

    const formData: FormData = new FormData(form);
    formData.set("hasjs", "1");
    const dataURL = new URLSearchParams(formData as any);
    try {
        const response = await fetch("api/security/signup", {
            method: "POST",
            body: dataURL
        }) ;
        const feedback = await response.text();

        if (feedback == "Success!"){
            const params = new URLSearchParams(window.location.search)
            const from =  params.get("from") ?? "/index.html"
            window.location.replace(from)
        } else {
            const errorDisplay = document.getElementById("errorDisplay") ?? err(feedback);
            errorDisplay.innerHTML = feedback;
        }

    } catch (e) {
        console.error(e);
    }

}


form?.addEventListener("submit", function(event){
    event.preventDefault();
    sendSignup()
})
