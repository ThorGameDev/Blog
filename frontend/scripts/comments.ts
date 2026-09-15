import { err } from "./lib/utils.ts"

async function appendReplies(targetElement: HTMLElement) {
    const response = await fetch("/api/blog/getReplies?lang=en&commentId=4");
    // TODO: Error handling
    const data = await response.text();
    // there has to be better ways to do this
    targetElement.innerHTML += data;
    // TODO: reRoute each generated reply
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
        commentExpand.addEventListener('click', function (e) {
            e.preventDefault();
            appendReplies(comment);
        })
    }
}

// Add interactivity to the root comments (The children are not going to be here at this point)
const comments = document.getElementsByClassName("comment");

for (let i = 0; i < comments.length; i++) {
    reRouteComment(comments[i] as HTMLElement);
}
