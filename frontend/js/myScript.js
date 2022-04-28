



// Get the login modal
let lmodal = document.getElementById("loginModal");

// Get the button that opens the login modal
let lbtn = document.getElementById("loginBtn");

// Get the button that opens the register modal on the login modal
let rbtnl = document.getElementById("registerBtnl");

// Get the <span> element that closes the login modal
let lspan = document.getElementsByClassName("lclose")[0];

// When the user clicks the button, open the login modal 
lbtn.onclick = function() {
  lmodal.style.display = "block";
}

// When the user clicks on login , close the login modal and open the register modal
rbtnl.onclick = function() {
  rmodal.style.display = "block";
  lmodal.style.display = "none";
}


// When the user clicks on <span> (x), close the login modal
lspan.onclick = function() {
  lmodal.style.display = "none";
}

//--------------------------------------------


// Get the register modal
let rmodal = document.getElementById("registerModal");

// Get the button that opens the login modal
let rbtn = document.getElementById("registerBtn");

// Get the button that opens the login modal on the register modal
let lbtnr = document.getElementById("loginBtnr");

// Get the <span> element that closes the login modal
let rspan = document.getElementsByClassName("rclose")[0];

// When the user clicks the button, open the login modal 
rbtn.onclick = function() {
  rmodal.style.display = "block";
}

// When the user clicks on login , close the register modal and open the login modal
lbtnr.onclick = function() {
  rmodal.style.display = "none";
  lmodal.style.display = "block";
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