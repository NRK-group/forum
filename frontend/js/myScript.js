



// Get the login modal
var lmodal = document.getElementById("loginModal");

// Get the button that opens the login modal
var lbtn = document.getElementById("loginBtn");

// Get the <span> element that closes the login modal
var lspan = document.getElementsByClassName("lclose")[0];

// When the user clicks the button, open the login modal 
lbtn.onclick = function() {
  lmodal.style.display = "block";
}

// When the user clicks on <span> (x), close the login modal
lspan.onclick = function() {
  lmodal.style.display = "none";
}

// Get the login modal
var rmodal = document.getElementById("registerModal");

// Get the button that opens the login modal
var rbtn = document.getElementById("registerBtn");

// Get the <span> element that closes the login modal
var rspan = document.getElementsByClassName("rclose")[0];

// When the user clicks the button, open the login modal 
rbtn.onclick = function() {
  rmodal.style.display = "block";
}

// When the user clicks on <span> (x), close the login modal
rspan.onclick = function() {
  rmodal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == lmodal) {
    lmodal.style.display = "none";
  } else if (event.target == rmodal) {
    rmodal.style.display = "none";
  }
}