export default ({ app }, inject) => {
  // Set the function directly on the context.app object
  app.htmlDecode = (input) => {
    var e = document.createElement('div');
    e.innerHTML = input;
    return e.childNodes.length === 0 ? "" : e.childNodes[0].nodeValue;
  }
}