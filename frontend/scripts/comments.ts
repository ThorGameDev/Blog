import { err, formJsEnhancement } from "./lib/utils.ts"
import { getAccountDetails } from "./lib/account.ts"

const langCode = document.getElementsByTagName("html")[0].getAttribute("lang")

async function appendReplies(targetElement: HTMLElement, commentId: string): Promise<void> {
    // unhide hidden comments if they exist
    const preloaded = targetElement.querySelectorAll(":scope > .comment");
    if (preloaded.length > 0) {
        for (let i = 0; i < preloaded.length; i++) {
            preloaded[i].setAttribute('class', 'comment');
        }
        return;
    }

    // Get replies from backend
    const response = await fetch(`/api/blog/getReplies?lang=${langCode}&commentId=${commentId}`);
    if (!response.ok) {
        err("Could not get replies from backend");
    }
    const data = await response.text();
    targetElement.insertAdjacentHTML('beforeend', data);

    const comments = targetElement.querySelectorAll(":scope > .comment");
    for (let i = 0; i < comments.length; i++) {
        reRouteComment(comments[i] as HTMLElement);
    }
}

function hideReplies(targetElement: HTMLElement): void {
    const comments = targetElement.querySelectorAll(":scope > .comment");
    for (let i = 0; i < comments.length; i++) {
        comments[i].setAttribute('class', 'comment hidden');
    }
}


function postReply(targetElement: HTMLElement, commentForm: HTMLFormElement){
    const accountDetails = getAccountDetails() ?? err("No account, could not post reply");

    //create the comment
    const newComment = document.createElement("article");
    newComment.setAttribute("class", "comment");

    // Create the header
    const header = document.createElement("header");

    const profilePic = document.createElement("img");
    profilePic.setAttribute("src", accountDetails[1]);
    header.append(profilePic);

    const profileName = document.createElement("h3");
    profileName.innerText = accountDetails[0];
    header.append(profileName);

    newComment.append(header);

    // Create the body
    const body = document.createElement("p");
    const formData = new FormData(commentForm);
    const replyText = formData.get("replyData") as string;
    body.innerText = replyText;
    newComment.append(body);

    // TODO: Create the footer.
    // This is tricky, because as of now, the created comment ID is not shown.
    // Although, it would be a half-decent cheat to make it impossible to reply to your own comments in general
    // Then the footer would become un-necessary
    // Another option is to force a comments expand on reply, and to add a "check if changed" feature

    targetElement.append(newComment);
}

function toggleCommentReplyForm(targetElement: HTMLElement, commentId: string): void {
    // if the form already exists, toggle visibility
    const oldReplyForm = targetElement.querySelector(":scope > .replyForm");
    if (oldReplyForm != null) {
        if (oldReplyForm.getAttribute("class") == "replyForm hidden") {
            oldReplyForm.setAttribute("class", "replyForm");
        } else {
            oldReplyForm.setAttribute("class", "replyForm hidden");
        }
        return;
    }

    // Create form
    const replyForm = document.createElement("form");
    replyForm.setAttribute("class", "replyForm");
    replyForm.setAttribute("action", `/api/blog/reply?commentId=${commentId}&lang=${langCode}`);
    replyForm.setAttribute("method", "post");
    formJsEnhancement(replyForm, function() {
        postReply(targetElement, replyForm);
    })

    // Add input field
    const replyData = document.createElement("textarea");
    replyData.setAttribute("name", "replyData");
    replyForm.append(replyData);

    // Add submit button
    const submitButton = document.createElement("button");
    submitButton.setAttribute("type", "submit");
    submitButton.innerHTML = "Reply"; // TODO: Translate!!!
    replyForm.append(submitButton);

    //
    targetElement.querySelector(":scope > footer")?.insertAdjacentElement('afterend', replyForm);
}

function reRouteComment(comment: HTMLElement): void {
    const replyLink = comment.getElementsByClassName("reply")[0];
    if (replyLink != null) {
        replyLink.addEventListener('click', function (e) {
            e.preventDefault();

            // Get commentId
            const targLink = replyLink.getAttribute("href") ?? err("No link attached to Comment Expand Link");
            const params = new URL(targLink, window.location.origin).searchParams;
            const commentId = params.get("commentId") ?? err("Comment link had no commentId")

            // TODO: Generate a reply form
            toggleCommentReplyForm(comment, commentId)
        })
    }
    const commentExpand = comment.getElementsByClassName("commentExpand")[0]
    if (commentExpand != null) {
        // It's probably better to use a boolean to toggle states, but for now, the oneshot loop works really well
        commentExpand.addEventListener('click', function handler(e) {
            // Prevent the standard no-js behavior
            e.preventDefault();

            // Get commentId
            const targLink = commentExpand.getAttribute("href") ?? err("No link attached to Comment Expand Link");
            const params = new URL(targLink, window.location.origin).searchParams;
            const commentId = params.get("commentId") ?? err("Comment link had no commentId")

            // Add reply list to the end of the comment
            appendReplies(comment, commentId);

            // Replace the link functionality with a "Hide replies" function
            commentExpand.addEventListener('click', function (e) {
                // Prevent the standard no-js behavior
                e.preventDefault();

                hideReplies(comment);

                // Once the replies are hidden, re-attach the default "show" behavior
                commentExpand.addEventListener('click', handler, { once: true })
            }, { once: true })
        }, { once: true })
    }
}

// Add interactivity to the root comments (The children are not going to be here at this point)
const comments = document.getElementsByClassName("comment");

for (let i = 0; i < comments.length; i++) {
    reRouteComment(comments[i] as HTMLElement);
}
