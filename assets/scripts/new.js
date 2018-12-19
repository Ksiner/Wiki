var colratio = Math.floor(1077/73);
$("#Content .art-container_cont-text").attr("cols",Math.floor($("body").width()/colratio));
$(window).resize(function(){
    $("#Content .art-container_cont-text").attr("cols",Math.floor($("body").width()/colratio));
})
