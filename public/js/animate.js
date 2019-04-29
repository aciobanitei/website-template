
// Make http request to server.
function makeHttpRequest(parameters, callbackFunction) {
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = callbackFunction;

  var message = JSON.stringify(parameters);
  xhttp.open("POST", "?t=" + Math.random(), true);
  xhttp.send(message);
}
  
  function slideOut( element ) {
    element.style.animation ="slide-out-top 1s";
    setTimeout(function() {
        //remove INTRO
        element.parentNode.removeChild(element);
        //start title animation
        setupHeader(true);
    }, 900);
}


function setupHeader(animate){
    var element = document.getElementById("HEADER");
    var str = element.getAttribute("data-title").split(" ");
    var t = "";
    for (var i = 0; i < str.length; i++) {
          if (animate)
          {
                t = t + "<li style='opacity:0;animation: blur-in .5s linear forwards;animation-delay:"+i*0.1+"s'>"+str[i]+" </li>";
          }
          else
          {
                t = t + "<li>"+str[i]+" </li>";
          }
    }
    element.innerHTML ="<div class='site-header'><ul>"+t+"</ul></div>";
}



// exiting intro 
function onIntroExit() {
  
  document.body.style.overflow = "visible";
  //document.getElementById("INTRO").classList.add("w3-hide");
  
  
  document.getElementById("HEADER").classList.remove("w3-hide");
  document.getElementById("NAV_CONTAINER").classList.remove("w3-hide");
  document.getElementById("CONTENT").classList.remove("w3-hide");
  document.getElementById("FOOTER").classList.remove("w3-hide");


  slideOut(document.getElementById("INTRO"));
}



// Menu actions and updates
function onMenuItemSelect(item) {
  callback = function() {
    if (this.readyState == 4 && this.status == 200) {
    document.body.innerHTML = this.responseText;
     // state has changed -- update history
     window.history.pushState(null, null, '/' + item);
    }

  };
  makeHttpRequest({'path':"/" + item}, callback);
  updateMenuButtons(item);
  onSmallMenuHide();
}

			

// When the user scrolls the page, execute onWindowScroll
window.onscroll = function() {onWindowScroll()};

// Add the sticky class to the navbar when you reach its scroll position. Remove "sticky" when you leave the scroll position
function onWindowScroll() {
  // Get the navbar

  var navbar = document.getElementById("NAV_CONTAINER");
  // Get the offset position of the navbar
  if (window.pageYOffset >= document.getElementById("HEADER").offsetHeight) {
    navbar.classList.add("sticky");
    document.getElementById("NAV_MOBILE").classList.add("sticky");
    //document.getElementById("NAV_MOBILE").style.top = navbar.offsetHeight;
  } else {
    navbar.classList.remove("sticky");
    document.getElementById("NAV_MOBILE").classList.remove("sticky");
  }

}





// Used to toggle the menu on small screens when clicking on the menu button
function onMenuClick() {
    var x = document.getElementById("NAV_MOBILE");
    if (x.className.indexOf(" w3-show") == -1) {
        x.className += " w3-show";
    } else {
        x.className = x.className.replace(" w3-show", "");
    }
}

function onSmallMenuHide() {
    var x = document.getElementById("NAV_MOBILE");
    if (x.className.indexOf("w3-show") > -1) {
        x.className = x.className.replace(" w3-show", "");
    }
}

function updateMenuButtons(id) {
    var menu = document.getElementById("NAV");

    var children = menu.children;
    for (var i = 0; i < children.length; i++) {
      if (children[i].className.indexOf("active") > -1 && children[i].id != id)
          children[i].className = children[i].className.replace(" active", " w3-hide-small");
      else {
        if (children[i].id == id)
          children[i].className = children[i].className.replace(" w3-hide-small", " active");
      }
    }

    var mobile_menu = document.getElementById("NAV_MOBILE");
    var mobileid = "mobile" + id;
    children = mobile_menu.children;
    for (var i = 0; i < children.length; i++) {
      if (children[i].className.indexOf("w3-hide ") > -1 && children[i].id != mobileid)
          children[i].className = children[i].className.replace("w3-hide ", "w3-show ");
      else {
        if (children[i].id == mobileid)
          children[i].className = children[i].className.replace("w3-show ", "w3-hide ");
      }
    }
}




// PAGINATION
var list;
var pageList;
var currentPage = 0;
var numberPerPage = 3;
var numberOfPages = 0;
var listWindow = 5;
var listStartIndex = 0;

function getNumberOfPages() {
    return Math.ceil(list.length / numberPerPage);
}

function nextPage() {
    if (currentPage < numberOfPages - 1)
     currentPage += 1;

    if ((listStartIndex + listWindow - 1) < currentPage)
      listStartIndex +=1;

    loadList();
}

function previousPage() {
    if (currentPage > 0)
      currentPage -= 1;

    if (listStartIndex > currentPage)
      listStartIndex -=1;
    loadList();
}

function jumpBackwards() {
    if (currentPage >= listWindow) {
      currentPage -= listWindow;
      if (listStartIndex - listWindow > 0) {
          listStartIndex -= listWindow;
      }
      else{
        listStartIndex = 0;
      }
    }
    else {
       currentPage = 0;
       listStartIndex = 0;
    }
    loadList();
}

function jumpForwards() {
     if (currentPage <= (numberOfPages - listWindow)) {
      currentPage += listWindow;
      if ((listStartIndex + listWindow) < numberOfPages - listWindow) {
          listStartIndex += listWindow;
      }
      else{
        listStartIndex = numberOfPages - listWindow;
      }
    }
    else {
       currentPage = numberOfPages;
    }
    loadList();
}

function goToPage(p) {
    currentPage = p;
    if (p < listStartIndex) {
      if (p > 2)
        listStartIndex = p - 2;
      else
        listStartIndex = 0;
    }

    loadList();
}

function loadList() {
    var begin = (currentPage * numberPerPage);
    var end = begin + numberPerPage;

    pageList = list.slice(begin, end);
    drawList();
    buildListElement();
}

function drawList() {
    var s = "";

    for (r = 0; r < pageList.length; r++) {
        s += pageList[r] + "<br/>";
    }
    document.getElementById("PROJECTLIST").innerHTML = s;
}

function buildListElement() {
    var s = "";
    if (listStartIndex > 0)
      s += "<a onclick=\"jumpBackwards()\">&laquo;</a>";
    else
      s += "<a class=\"disabled\">&laquo;</a>";

    if (currentPage > 0)
      s+= "<a onclick=\"previousPage()\">&lt;</a>";
    else
       s+= "<a class=\"disabled\">&lt;</a>";

    for (r = 0; r < listWindow; r++)
    {
        if ((r+listStartIndex) < numberOfPages) {
           if (r + listStartIndex == currentPage)
          s += "<a class=\"active\">" + (r + listStartIndex + 1) + "</a>";
          else
          s += "<a onclick=\"goToPage(" + (r +listStartIndex) +")\">" + (r + listStartIndex + 1) + "</a>";
        }
    }
    if (currentPage < numberOfPages - 1)
      s+= "<a onclick=\"nextPage()\">&gt;</a>";
    else
      s+= "<a class=\"disabled\">&gt;</a>";


   if ((listStartIndex  + listWindow) < numberOfPages - 1)
      s += "<a onclick=\"jumpForwards()\" >&raquo;</a>";
   else
      s += "<a class=\"disabled\">&raquo;</a>";

   document.getElementById("LIST").innerHTML = s;
}

function load_pagination(pageid) {
    updateMenuButtons(pageid);
    onSmallMenuHide();

    if (pageid != "projects")
      return;

    var blogposts = document.getElementById("PROJECTLIST");

    if (!blogposts)
      return;

    var children = blogposts.children;

    list = new Array();
    pageList = new Array();

    for (var i = 0; i < children.length; i++) {
        list.push(children[i].outerHTML);
    }

    numberOfPages  = getNumberOfPages();
    listStartIndex = 0;
    currentPage    = 0;
    loadList();
}
