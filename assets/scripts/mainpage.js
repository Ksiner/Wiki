var server = "http://localhost:8080";
jsSHA = require("jssha");
var shaObj = new jsSHA("SHA-512", "TEXT");
var jsonresponse
var temp = ""

function getMain() {
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      alert(xmlhttp.responseText);
      jsonresponse = JSON.parse(xmlhttp.responseText);
      //addNewList(1,"addedByJS");
      throughTree(0, jsonresponse);
    }
  };
  xmlhttp.open("GET", server + "/catTree", true);
  xmlhttp.send();
}

function throughTree(level, JSONarray) {
  for (i = 0; i < JSONarray.length; i += 1) {
    addNewList(level, JSONarray[i]["category"].name, level=0?"menu":JSONarray[i]["category"].id);

    if (!(JSONarray[i].childs === null)) {
      for (g = 0; g < JSONarray[i].childs.length; g += 1)
        throughTree(level + 1, JSONarray[i].childs);
    }
    if (!(JSONarray[i].articles === null)) {
      addArticleOnList(level, catId, name);
    }
  }
}

function addNewList(level, name, id) {
  var dropdown = document.getElementById(id);
  var Li = document.createElement("li");
  var elA = document.createElement("a");
  var elSpan = document.createElement("span");
  var elementUL = document.createElement("ul");
  elSpan.classList.add("caret");
  elA.href = "#";
  elA.tabIndex = "-1";
  elA.classList.add("test");
  elA.textContent = name;
  elementUL.classList.add("dropdown-menu");
  elementUL.id = id;
  elA.appendChild(elSpan);
  Li.appendChild(elA);
  Li.appendChild(elementUL);
  dropdown.appendChild(Li);
}

//TODO Добавить идентификатор для получения текста
function addArticleOnList(level, catName, name) {

}
getMain();
//addNewList(2,"addedByJS");