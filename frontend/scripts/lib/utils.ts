export function err(errorMessage: string): never {
    throw new Error(errorMessage);
}

async function baseFormBehavior(form: HTMLFormElement, after: () => void) {
    try {
        const formData: FormData = new FormData(form);
        const dataURL = new URLSearchParams(formData as any);
        const response = await fetch(form.action, {
            method: form.method,
            body: dataURL
        });

        if (response.redirected) {
            if (window.location.href != response.url) {
                window.location.replace(response.url);
            } else {
                after();
            }
        } else {
            err(await response.text())
        }

    } catch (e) {
        console.error(e);
    }
}

export function formJsEnhancement(form: HTMLFormElement, after: () => void) {
    form.addEventListener('submit', function (event) {
        event.preventDefault();
        baseFormBehavior(form, after);
    })
}
