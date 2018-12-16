var server = "http://localhost:8080";
/*jsSHA = require("jssha");
var shaObj = new jsSHA("SHA-512", "TEXT");*/
var jsonresponse
var temp = ""

function getTree() {
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      //alert(xmlhttp.responseText);
      jsonresponse = JSON.parse(xmlhttp.responseText);
      throughTree(0, jsonresponse);
      $(document).ready(function () {
        $('.dropdown a.test').on("click", function (e) {
          $(this).next('ul').toggle();
          e.stopPropagation();
          e.preventDefault();
        });
      });
      document.querySelector('.btn.btn-default.dropdown-toggle').disabled = false;
      //addNewList(1,"addedByJS");

      //addNewList("test1","menu1","menu2")
    }
  };
  xmlhttp.open("GET", server + "/catTree", true);
  xmlhttp.send();
}

function getMainArticles() {
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      jsonresponse = JSON.parse(xmlhttp.responseText);
      setArticles(jsonresponse);
    }
  };
  xmlhttp.open("GET", server + "/init", true);
  xmlhttp.send();
}

function throughTree(level, JSONarray) {
  for (let i = 0; i < JSONarray.length; i += 1) {
    
    //debugger;
    
    var parentid = level == 0 ? "menu" : JSONarray[i]["category"].parentid;
    var currentid = JSONarray[i]["category"].id;
    addNewList(JSONarray[i]["category"].name, parentid, currentid);
    addArticleOnList(JSONarray[i].articles);
    if (JSONarray[i].childs !== null) {
      throughTree(level + 1, JSONarray[i].childs);
    }
    else{
      addButtons(currentid);
    }
    if (i===JSONarray.length-1)
    {
      addButtons(parentid);
    }
  }
}

function addNewList(name, idParent, idChild) {
  var dropdown = document.getElementById(idParent);
  var Li = document.createElement("li");
  var elA = document.createElement("a");
  var elSpan = document.createElement("span");
  var elementUL = document.createElement("ul");
  elSpan.classList.add("caret");
  elA.href = "#";
  elA.tabIndex = "-1";
  elA.classList.add("test");
  elA.classList.add("category");
  elA.textContent = name;
  elementUL.classList.add("dropdown-menu");
  elementUL.id = idChild;
  elA.appendChild(elSpan);
  Li.appendChild(elA);
  Li.appendChild(elementUL);
  dropdown.appendChild(Li);
  
}

//TODO Добавить идентификатор для получения текста
function addArticleOnList(articles) {
  if (articles === null)
    return;
  for (i = 0; i < articles.length; i += 1) {
    var dropdown = document.getElementById(articles[i].catid);
    var Li = document.createElement("li");
    var elA = document.createElement("a");
    Li.id = articles[i].id;
    /*elA.onclick=(e)=>{
      console.log(e.target.parentElement.id+" was clicked");
    };*/
    elA.onclick = requestArticle;
    elA.href = "#";
    elA.tabIndex = "-1";
    elA.classList.add("test");
    elA.classList.add("article");
    elA.textContent = articles[i].header;
    Li.appendChild(elA);
    dropdown.appendChild(Li);
  }
}

function addButtons(parentid){
  var dropdown = document.getElementById(parentid);
  var Li = document.createElement("li");
  //var elA = document.createElement("a");
  var catButton = document.createElement("button");
  var artButton = document.createElement("button");
  var articleForm = document.createElement("form");
  var categoryForm = document.createElement("form");
  var articleNameInput = document.createElement("input");
  var CategoryNameInput = document.createElement("input");

  articleNameInput.name = "art";
  artButton.type="submit";
  articleForm.appendChild(articleNameInput);
  articleForm.appendChild(artButton);
  articleForm.action="/"+parentid+"/article/create";

  CategoryNameInput.name = "cat";
  catButton.type="submit";
  categoryForm.appendChild(CategoryNameInput);
  categoryForm.appendChild(catButton);
  categoryForm.action="/"+parentid+"/create";


  Li.appendChild(categoryForm);
  Li.appendChild(articleForm);
  dropdown.appendChild(Li);
}

/*function addNewArticle(event){
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      jsonresponse = JSON.parse(xmlhttp.responseText);
      setArticles(jsonresponse);
    }
  };
  xmlhttp.open("POST", server + "/catTree", true);
  xmlhttp.send();
}*/

function requestArticle(event) {
  articleId = event.target.parentElement.id;
  categoryId = event.target.parentElement.parentElement.id;
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      alert(xmlhttp.responseText);
      //Переход на страницу с артиклом
    }
  }
  xmlhttp.open("GET", server + "/" + categoryId + "/article/" + articleId, true);
  xmlhttp.send();

}

function setArticles(JSONarticles) {
  counter = 0
  elements = document.querySelectorAll(".col-sm-4 .panel.panel-primary");
  for (i = 0; i < elements.length; i += 1) {
    elements[i].children[0].innerText = JSONarticles[counter].header;
    elements[i].children[1].children[0].src = "data:image/png;base64," + JSONarticles[counter].pic;
    elements[i].children[2].innerText = JSONarticles[counter].content.substring(0,3)+"...";
    counter++;
  }
}
getTree();
getMainArticles();