import { err } from "./utils.ts"

let storedAccount: [string, string] | null = null;
export function getAccountDetails(): [string, string] | null {
    // If an account already exists, return it without doing anything
    if (storedAccount != null) {
        return storedAccount;
    }

    // Check if there is a "No account" indicator, and if so return nothing
    const accountlessIcon = document.getElementById("accountlessIcon");
    if (accountlessIcon != null) {
        return null
    }

    // Otherwise, get details about the existing account
    const accountIcon = document.getElementById("accountIcon") ?? err("Could not find account icon");
    const accountName = accountIcon.getAttribute("title") ?? err("No title found on account icon");
    const accountImg = accountIcon.getElementsByTagName("img")[0].getAttribute("src") ?? err("No image attached");

    // Finally, store the account details, and return them
    storedAccount = [accountName, accountImg]
    return storedAccount
}
