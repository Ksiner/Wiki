var server = "http://localhost:8080";

function init() {

}

function getMain() {
  let xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
    if ( xmlhttp.readyState == 4 && xmlhttp.status == 200 ) {
      alert( xmlhttp.responseText );
    }
  };
  debugger;
  xmlhttp.open( "GET", server, true );
  xmlhttp.send();
}
getMain();