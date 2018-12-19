var server = "http://localhost:8080";
var jsonresponse
var userName;
var allarticles
var currentArticle;
var userName;

//Получить дерево категорий и статей, вынести всё в drop-down меню
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

//Запросить все статьи и занести их в стартовую страницу
function getMainArticles() {
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

//Вынести дерево в dropdown меню
function throughTree(level, JSONarray) {
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

//Отобразить новую категорию в меню
function addNewList(name, idParent, idChild) {
  if(idParent==="null"){
    idParent="menu";
  }
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
  dropdown.insertBefore(Li, dropdown.childNodes[0]);

}

//Отобразить статьи в категории
function addArticleOnList(articles) {
  if (articles === null)
    return;
  for (i = 0; i < articles.length; i += 1) {
    let categoryid= articles[i].catid==="null"?"menu":articles[i].catid;
    var dropdown = document.getElementById(categoryid);
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
    if (dropdown.childElementCount >= 4 && dropdown.lastElementChild.childElementCount > 1)
      dropdown.insertBefore(Li, dropdown.childNodes[dropdown.childElementCount - 1]);
    else
      dropdown.appendChild(Li);
  }
}
//Добавление кнопок для добавления статей и категорий
function addButtons(parentid) {
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
  artButton.type="button";
  catButton.type="button";
  articleNameInput.name = "art";
  //artButton.type = "submit";
  artButton.onclick = addNewArticle
  articleForm.appendChild(articleNameInput);
  articleForm.appendChild(artButton);
  articleForm.classList.add("addition-form");

  CategoryNameInput.name = "cat";
  catButton.onclick = addNewCategory;
  categoryForm.appendChild(CategoryNameInput);
  categoryForm.appendChild(catButton);
  categoryForm.classList.add("addition-form");



  Li.appendChild(categoryForm);
  Li.appendChild(articleForm);
  dropdown.appendChild(Li);
}

//Отдать название новой статьи по id категории, получить статью НЕПРОВЕРЕНО
function addNewArticle(event) {
  let catid = null
  if (event.target.parentElement.parentElement.parentElement.id !== "menu")
    catid = event.target.parentElement.parentElement.parentElement.id;
  var string = server + "/" + catid + "/article/create";
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.open("POST",string , true);
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      let article = JSON.parse(xmlhttp.responseText)
      addArticleOnList([article]);
      articlePage(article);
    }
  }
  //xmlhttp.setRequestHeader("token", localStorage.getItem("tokenArticles"));
  xmlhttp.send(JSON.stringify({
    art: event.target.parentElement.children[0].value
  }));

  /*var url = server + "/" + "e46c203e-f701-11e8-a8a5-bcaec5906742" + "/article/create";
  let xmlhttp = new XMLHttpRequest();
    xmlhttp.open("POST",url,true);
    xmlhttp.onreadystatechange = function(){
        if (xmlhttp.readyState === 4 && xmlhttp.status === 200){
            alert("Message sended!");
        }
    }
    xmlhttp.send(JSON.stringify({art: event.target.parentElement.children[0].value}));*/

}

// Отдать название новой категории, по id категории НЕ ПРОВЕРЕНО
function addNewCategory(event) {
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
  //xmlhttp.setRequestHeader("token", localStorage.getItem("tokenArticles"));
  xmlhttp.send(JSON.stringify({
    cat: event.target.parentElement.children[0].value
  }));
}

//ПЕРЕХОД НА СТРАНИЦУ СО СТАТЬЁЙ
function goToArticle(event) { //
  articleId = event.target.parentElement.id;
  let article;
  for (let i = 0; i < allarticles.length; i++) {
    if (allarticles[i].id === articleId) {
      article = allarticles[i];
      break;
    }
  }
  articlePage(article);
}

//Добавляет все статьи на главную страницу
function setArticles() {
  let main = document.getElementById("allArticles");
  main.innerText = "";
  main.innerHTML = "";
  for (i = 0; i < allarticles.length; i++) {
    addArticleInPage(allarticles[i]);
  }
}

function goToArticleFromPage(event) {
  let target = event.target;
  while (target.id === "") {
    target = target.parentElement;
  }

  let article;
  for (let i = 0; i < allarticles.length; i++) {
    if (allarticles[i].id === target.id) {
      article = allarticles[i];
      break;
    }
  }
  articlePage(article);
}

//Отображает статью
function articlePage(article) {
  currentArticle = article;
  if(!document.getElementById("allArticles").classList.contains("hidden")){
    changeMainPage();
  }
  if (document.getElementById("constHeader").classList.contains("hidden")){
    changeArticle(null);
  }
  let articlePage = document.getElementById("articlePage");
  document.getElementById("constHeader").textContent = article.header;
  let content = document.getElementById("constContent").querySelector(".art-container_cont-text")
  content.innerText = article.content;
  let image = articlePage.querySelector(".art-container_img-container")[0];
  let hiddenContent = document.getElementById("Content").querySelector(".art-container_cont-text");
  hiddenContent.innerText = article.content;
  let hiddenHeader = document.getElementById("Header");
  hiddenHeader.textContent = article.header;
  //image.src = article.image;
}

function changeMainPage() {
  let main = document.getElementById("allArticles");
  let articlePage = document.getElementById("articlePage");
  changeVision(main, articlePage);
  // if (main.classList.contains("hidden")){
  //   main.classList.remove("hidden");
  //   articlePage.classList.add("hidden");
  // }
  // else
  // {
  //   articlePage.classList.remove("hidden");
  //   main.classList.add("hidden");
  // }
}

//Добавляет новую статью на главную страницу
function addArticleInPage(article) {
  let main = document.getElementById("allArticles");
  var column = document.createElement("div");
  column.classList.add("col-sm-4");
  column.onclick = goToArticleFromPage;
  column.id = article.id;
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

//Поиск при вводе
function searchOnInput(event) {
  let main = document.getElementById("allArticles");
  main.innerText = "";
  main.innerHTML = "";
  for (i = 0; i < allarticles.length; i++) {
    if (allarticles[i].content.includes(event.target.value)) {
      addArticleInPage(allarticles[i]);
    }
  }
}

//Отправляет хешированные логин и пароль, получет токен пользователя
function login(event) {
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.open("POST", server + "/login", true);
  log = document.getElementById("log");
  pswd = document.getElementById("pswd");
  userName = log.value;
  let hasher = new jsSHA("SHA-512", "TEXT");
  hasher.update(log.value);
  let hashLog = hasher.getHash("B64");
  hasher = new jsSHA("SHA-512", "TEXT");
  hasher.update(pswd.value);
  let hashpass = hasher.getHash("B64");
  
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
      jsonobj = JSON.parse(xmlhttp.responseText);
      localStorage.setItem("tokenArticles", jsonobj["token"]);
      changeLoginForms();
      setName()
    }
  }
  xmlhttp.send(JSON.stringify({
    login: hashLog,
    pass: hashpass
  }));
}

function VKRegistration(){
  VK.Observer.subscribe('auth.login', function(response){
    changeLoginForms();
    setName(response.user.first_name)
  });
  VK.Auth.login(null);
}

//Создаёт нового пользователя, отправляя хешированные логин и пароль, получет токен пользователя
function registration(event) {
  log = document.getElementById("log");
  pswd = document.getElementById("pswd");
  userName = log.value;
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
      localStorage.setItem("tokenArticles", jsonobj["token"]);
      changeLoginForms();
      setName();
    }
  }
  xmlhttp.open("POST", server + "/reg", true);
  xmlhttp.send(JSON.stringify({
    login: hashLog,
    pass: hashpass
  }));
}

function setName(name){
  let profile = document.getElementById("profile")[0];
  if(name===null){
    profile.innerText = userName;
  }
  else
  {
    profile.innerText = name;
  }
}

function changeVision(first, second) {
  if (first.classList.contains("hidden")) {
    first.classList.remove("hidden");
    second.classList.add("hidden");
  } else {
    first.classList.add("hidden");
    second.classList.remove("hidden");
  }
}

//Меняет отображение верхнего меню (Форма входа, на форму профиля и наоборот)
function changeLoginForms() {
  let logForm = document.getElementById("login")
  let profForm = document.getElementById("profile");
  changeVision(logForm, profForm);
  /*if(logForm.classList.contains("hidden")){
    logForm.classList.remove("hidden");
    profForm.classList.add("hidden");
  }
  else
  {
    logForm.classList.add("hidden");
    profForm.classList.remove("hidden");
  }*/
}

//Выходит из профиля, отправляя токен пользователя на сервер
function logout() {
  userName=null;
  if(localStorage.getItem("tokenArticles")===null)
    return;
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
    }
  }
  xmlhttp.open("POST", server + "/logout", true);
  xmlhttp.send(JSON.stringify({
    token: localStorage.getItem("tokenArticles"),
  }));
  changeLoginForms();
  localStorage.removeItem("tokenArticles");
  document.getElementById("pswd").value="";
}

function changeArticle(event) {
  let first = document.getElementById("constHeader");
  let second = document.getElementById("Header");
  changeVision(first, second)
  first = document.getElementById("constContent");
  second = document.getElementById("Content");
  changeVision(first, second)
  first = document.getElementById("imagebtn");
  if (first.classList.contains("hidden"))
    first.classList.remove("hidden");
  else
    first.classList.add("hidden");
  first = document.getElementById("savebtn");
  second = document.getElementById("changebtn");
  changeVision(first, second)
}

function saveArticle(event){
  categoryId = currentArticle.catid;
  let header = document.getElementById("Header").value;
  let content = document.getElementById("Content").querySelector(".art-container_cont-text").value;
  currentArticle.header = header;
  currentArticle.content = content;
  document.getElementById("constContent").querySelector(".art-container_cont-text").innerText = content;
  document.getElementById("constHeader").textContent = header;
  changeArticle(null);
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function () {
    if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
    }
  }
  xmlhttp.open("POST", server + "/" + categoryId + "/article", true);
  xmlhttp.send(JSON.stringify(currentArticle));
  UpdateTree(currentArticle);
}

function UpdateTree(article){
  let articles = document.getElementById("menu").querySelectorAll(".article");
  for (let i=0;i<articles.length;i++){
    if(article.id === articles[i].parentElement.id){
      articles[i].textContent=article.header;
    }
  }
}

function exit(){
  logout();
}


element = document.getElementById("example-search-input")
element.oninput = searchOnInput;
element = document.querySelector("#entrebtn");
element.onclick = login;
element = document.querySelector("#regisbtn");
element.onclick = registration;
element = document.querySelector("#logoutbtn");
element.onclick = logout;
element = document.querySelector("#changebtn");
element.onclick = changeArticle;
element = document.querySelector("#savebtn");
element.onclick = saveArticle;
element = document.querySelector("#vk");
element.onclick = VKRegistration;
window.onbeforeunload=exit;
getTree();
getMainArticles();
response = VK.Auth.getLoginStatus(function(response){
  let e = response;
  });

// VK.Observer.subscribe('auth.login', function(response){
//   changeLoginForms();
//   setName(response.user.first_name)
// });
// VK.Auth.login(null);