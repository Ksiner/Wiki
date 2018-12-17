var server = "http://localhost:8080";
/*jsSHA = require("jssha");
var shaObj = new jsSHA("SHA-512", "TEXT");*/
var jsonresponse

var allarticles
var temp = ""
//jsSHA = require("jssha");
function getTree() { //Получить дерево категорий и статей, вынести всё в drop-down меню
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

function getMainArticles() { //Запросить все статьи и занести их в стартовую страницу
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      allarticles = JSON.parse(xmlhttp.responseText);
      setArticles();
    }
  };
  xmlhttp.open("GET", server + "/init", true);
  xmlhttp.send();
}

function throughTree(level, JSONarray) { //Вынести дерево в dropdown меню
  for (let i = 0; i < JSONarray.length; i += 1) {

    //debugger;

    var parentid = level == 0 ? "menu" : JSONarray[i]["category"].parentid;
    var currentid = JSONarray[i]["category"].id;
    addNewList(JSONarray[i]["category"].name, parentid, currentid);
    addArticleOnList(JSONarray[i].articles);
    if (JSONarray[i].childs !== null) {
      throughTree(level + 1, JSONarray[i].childs);
    } else {
      addButtons(currentid);
    }
    if (i === JSONarray.length - 1) {
      addButtons(parentid);
    }
  }
}

function addNewList(name, idParent, idChild) { //Отобразить новую категорию в меню
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


function addArticleOnList(articles) { //Отобразить статьи в категории
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
    elA.onclick = goToArticle;
    elA.href = "#";
    elA.tabIndex = "-1";
    elA.classList.add("test");
    elA.classList.add("article");
    elA.textContent = articles[i].header;
    Li.appendChild(elA);
    dropdown.appendChild(Li);
  }
}

function addButtons(parentid) { //Добавление кнопок для добавления статей и категорий
  var dropdown = document.getElementById(parentid);
  var Li = document.createElement("li");
  //var elA = document.createElement("a");
  var catButton = document.createElement("button");
  var artButton = document.createElement("button");
  var articleForm = document.createElement("form");
  var categoryForm = document.createElement("form");
  var articleNameInput = document.createElement("input");
  var CategoryNameInput = document.createElement("input");
  catButton.innerText = "Категория";
  artButton.innerText = "Статья";
  articleNameInput.name = "art";
  //artButton.type = "submit";
  artButton.onclick = addNewArticle
  articleForm.appendChild(articleNameInput);
  articleForm.appendChild(artButton);

  //articleForm.action = "/" + parentid + "/article/create";

  CategoryNameInput.name = "cat";
  //catButton.type = "submit";
  catButton.onclick = addNewCategory;
  categoryForm.appendChild(CategoryNameInput);
  categoryForm.appendChild(catButton);
  //categoryForm.action = "/" + parentid + "/create";


  Li.appendChild(categoryForm);
  Li.appendChild(articleForm);
  dropdown.appendChild(Li);
}

function addNewArticle(event) { //Отдать название новой статьи по id категории, получить статью НЕПРОВЕРЕНО
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      addArticleOnList([JSON.parse(xmlhttp.responseText)]);
    }
  }
  let catid = null
  if (event.target.parentElement.parentElement.parentElement.id !== "menu")
    catid = event.target.parentElement.parentElement.parentElement.id
  xmlhttp.open("POST", server + "/" + catid + "/article/create", true);
  xmlhttp.send(JSON.stringify({
    art: event.target.parentElement.children[0].value
  }));
}

function addNewCategory(event) { // Отдать название новой категории, по id категории НЕ ПРОВЕРЕНО
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      jsonobj = JSON.parse(xmlhttp.responseText)
      addNewList(jsonobj.name, jsonobj.parentid, jsonobj.id);
    }
  }
  catid = null
  if (event.target.parentElement.parentElement.parentElement.id !== "menu")
    catid = event.target.parentElement.parentElement.parentElement.id;
  xmlhttp.open("POST", server + "/" + catid + "/create", true);
  xmlhttp.send(JSON.stringify({
    cat: event.target.parentElement.children[0].value
  }));
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

function goToArticle(event) { //ПЕРЕХОД НА СТРАНИЦУ СО СТАТЬЁЙ
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
function setArticles() { //Занести статьи на главную страницу
  /*counter = 0
  elements = document.querySelectorAll(".col-sm-4 .panel.panel-primary");
  for (i = 0; i < elements.length; i += 1) {
    elements[i].children[0].innerText = allarticles[counter].header;
    elements[i].children[1].children[0].src = "data:image/png;base64," + allarticles[counter].pic;
    elements[i].children[2].innerText = allarticles[counter].content.substring(0, 3) + "...";
    counter++;
  }*/

  main = document.getElementById("allArticles");
  main.innerText="";
  main.innerHTML="";
  for(i=0;i<allarticles.length;i++){
    addArticleInPage(allarticles[i]);
  }
}

function addArticleInPage(article){
  var column = document.createElement("div");
      column.classList.add("col-sm-4");
    var panel = document.createElement("div");
    panel.classList.add("panel");
    panel.classList.add("panel-primary");
    var header = document.createElement("div");
    header.classList.add("panel-heading");
    var imgDiv = document.createElement("div");
    imgDiv.classList.add("panel-body");
    var footer = document.createElement("div");
    footer.classList.add("panel-footer");
    var myImg = document.createElement("img");
    myImg.classList.add("img-responsive");
    myImg.style = "width:100%";
    myImg.alt = "Image";
    header.innerText = article.header;
    footer.innerText = article.content.substring(0, 25) + "...";
    myImg.src = "data:image/png;base64," + article.pic;
    imgDiv.appendChild(myImg);
    panel.appendChild(header);
    panel.appendChild(imgDiv);
    panel.appendChild(footer);
    column.appendChild(panel);
    main.appendChild(column);
}

function searchOnInput(event){
  main = document.getElementById("allArticles");
  main.innerText="";
  main.innerHTML="";
  for(i=0;i<allarticles.length;i++){
    if(allarticles[i].content.includes(event.target.value))
    {
      addArticleInPage(allarticles[i]);
    }
  }
}

function login(event){//Вставить Хеширование
  log = document.getElementById("log");
  pswd = document.getElementById("pswd");
  let hasher = new jsSHA("SHA-512", "TEXT");
  hasher.update(log.value);
  let hashLog = hasher.getHash("B64");
  hasher = new jsSHA("SHA-512", "TEXT");
  hasher.update(pswd.value);
  let hashpass = hasher.getHash("B64");
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      jsonobj = JSON.parse(xmlhttp.responseText);
      localStorage.setItem("tokenArticles",jsonobj[0]);
    }
    else
    {
      alert ("Такой комбинации логина и пароля не обнаружено");
    }
  }
  xmlhttp.open("POST", server + "/login", true);
  xmlhttp.send(JSON.stringify({
    login: hashLog,
    pass: hashpass
  }));
}

function registration(event){
  log = document.getElementById("log");
  pswd = document.getElementById("pswd");
  let hasher = new jsSHA("SHA-512", "TEXT");
  hasher.update(log.value);
  let hashLog = hasher.getHash("B64");
  hasher = new jsSHA("SHA-512", "TEXT");
  hasher.update(pswd.value);
  let hashpass = hasher.getHash("B64");
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      jsonobj = JSON.parse(xmlhttp.responseText);
      localStorage.setItem("tokenArticles",jsonobj[0]);
    }
    else
    {
      alert ("Невозможно регистрация такого пользователя");
    }
  }
  xmlhttp.open("POST", server + "/reg", true);
  xmlhttp.send(JSON.stringify({
    login: hashLog,
    pass: hashpass
  }));
}
element = document.getElementById("example-search-input")
element.oninput = searchOnInput;
element = document.querySelector("#entrebtn");
element.onclick = login;
element = document.querySelector("#regisbtn");
element.onclick=registration;
getTree();
getMainArticles();