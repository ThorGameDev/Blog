document.body.appendChild(document.createTextNode("Hello World! From typescript!"));

async function callApi(){
    const data = await fetch("api/helloWorld")
    const text = await data.text()
    document.body.appendChild(document.createTextNode(text));
}
callApi()
