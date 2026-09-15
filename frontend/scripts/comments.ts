import { err } from "./lib/utils.ts"

const langCode = document.getElementsByTagName("html")[0].getAttribute("lang")

async function appendReplies(targetElement: HTMLElement, commentId: string): Promise<void> {
    // unhide hidden comments if they exist
    const preloaded = targetElement.querySelectorAll(":scope > .comment")
    if (preloaded.length > 0) {
        for (let i = 0; i < preloaded.length; i++) {
            preloaded[i].setAttribute('class', 'comment')
        }
        return
    }

    // Get replies from backend
    const response = await fetch(`/api/blog/getReplies?lang=${langCode}&commentId=${commentId}`);
    if (!response.ok) {
        err("Could not get replies from backend")
    }
    const data = await response.text();
    targetElement.insertAdjacentHTML('beforeend', data);

    const comments = targetElement.querySelectorAll(":scope > .comment")
    for (let i = 0; i < comments.length; i++) {
        reRouteComment(comments[i] as HTMLElement);
    }
}

async function hideReplies(targetElement: HTMLElement): Promise<void> {
    const comments = targetElement.querySelectorAll(":scope > .comment")
    for (let i = 0; i < comments.length; i++) {
        comments[i].setAttribute('class', 'comment hidden')
    }
}

function reRouteComment(comment: HTMLElement): void {
    const replyLink = comment.getElementsByClassName("reply")[0];
    if (replyLink != null) {
        replyLink.addEventListener('click', function (e) {
            e.preventDefault();
            // TODO: Generate a reply form
            replyLink.replaceWith();
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
