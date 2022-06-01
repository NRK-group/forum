if (window.history.replaceState) {
  window.history.replaceState(null, null, window.location.href);
}
let session;

// Get the login modal
let lmodal = document.getElementById("loginModal");

// Get the button that opens the login modal
let lbtn = document.getElementById("loginBtn");

// Get the button that opens the register modal on the login modal
let rbtnl = document.getElementById("registerBtnl");

// Get the <span> element that closes the login modal
let lspan = document.getElementsByClassName("lclose")[0];

// Get the button that opens that login the user
let fbtn_login = document.getElementById("form-btn-login");

// close and reset the login modal
const Closelogin = () => {
  lmodal.style.display = "none";
  document.getElementById("userName").value = "";
  document.getElementById("password").value = "";
};

// When the user clicks the button, open the login modal
fbtn_login.onclick = function (event) {
  event.preventDefault();
  let data = new FormData();

  data.append("userName", document.getElementById("userName").value);
  data.append("password", document.getElementById("password").value);

  fetch("http://localhost:8800/login", {
    method: "POST",
    body: data,
  })
    .then(function (response) {
      return response.text();
    })
    .then(function (text) {
      //text is the server's response
      if (text[0] === "0") {
        document.getElementById("login-err").innerText = text.substring(1);
      } else if (text[0] === "1") {
        document.getElementById("login-err").innerText = "";

        for (let i = Not_Login_div.length - 1; i >= 0; --i) {
          Not_Login_div[i].style.display = "none";
        }

        for (let i = Login_div.length - 1; i >= 0; --i) {
          Login_div[i].style.display = "flex";
        }
        Closelogin();
        window.location.reload();
      }
    });
};

// When the user clicks the button, open the login modal
lbtn.onclick = function () {
  lmodal.style.display = "block";
};

// When the user clicks on login , close the login modal and open the register modal
rbtnl.onclick = function () {
  rmodal.style.display = "block";
  Closelogin();
};

// When the user clicks on <span> (x), close the login modal
lspan.onclick = function () {
  Closelogin();
};

//--------------------------------------------

// Get the register modal
let rmodal = document.getElementById("registerModal");

// Get the button that opens the login modal
let rbtn = document.getElementById("registerBtn");

// Get the button that opens the login modal on the register modal
let lbtnr = document.getElementById("loginBtnr");

// Get the <span> element that closes the login modal
let rspan = document.getElementsByClassName("rclose")[0];

// Get the button that opens that register a user
let fbtn_register = document.getElementById("form-btn-register");

let Not_Login_div = document.querySelectorAll('[id="Not_Login"]');
let Login_div = document.querySelectorAll('[id="Login"]');

const Closeregister = () => {
  document.getElementById("rUserName").value = "";
  document.getElementById("rPassword").value = "";
  document.getElementById("rEmail").value = "";
  rmodal.style.display = "none";
};

// When the user clicks the button, open the login modal
fbtn_register.onclick = function (event) {
  event.preventDefault();
  let data = new FormData();

  data.append("userName", document.getElementById("rUserName").value);
  data.append("password", document.getElementById("rPassword").value);
  data.append("email", document.getElementById("rEmail").value);

  fetch("http://localhost:8800/register", {
    method: "POST",
    body: data,
  })
    .then(function (response) {
      return response.text();
    })
    .then(function (text) {
      //text is the server's response
      if (text[0] === "0") {
        document.getElementById("register-err").innerText = text.substring(1);
      } else if (text[0] === "1") {
        document.getElementById("register-err").innerText = "";
        Closeregister();
        alert(text.substring(1));
      }
    });
};

// When the user clicks the button, open the login modal
rbtn.onclick = function () {
  rmodal.style.display = "block";
};

// When the user clicks on login , close the register modal and open the login modal
lbtnr.onclick = function () {
  Closeregister();
  lmodal.style.display = "block";
};

// When the user clicks on <span> (x), close the login modal
rspan.onclick = function () {
  Closeregister();
};

//--------------------------------------------

// logout function

// Get the button that logout the user
let logoutBtn = document.getElementById("logoutBtn");

// When the user clicks on logoutBtn that logout the user
logoutBtn.onclick = function () {
  let data = new FormData();
  fetch("http://localhost:8800/logout", {
    method: "POST",
    body: data,
  })
    .then(function (response) {
      return response.text();
    })
    .then(function (text) {
      //text is the server's response
      alert(text);
      window.location.reload();
    });

  for (let i = Not_Login_div.length - 1; i >= 0; --i) {
    Not_Login_div[i].style.display = "flex";
  }

  for (let i = Login_div.length - 1; i >= 0; --i) {
    Login_div[i].style.display = "none";
  }
};

//--------------------------------------------

// onload
const Onload = (cookie) => {
  if (cookie !== "") {
    session = cookie.split("&");

    if (window.location.search.substring(1).length > 2) {
      const query = window.location.search.substring(1);
      const token = query.split("access=")[1];
      token.split("-");
      const date = new Date(Date.now() + 3600 * 1000 * 24);
      document.cookie =
        "session_token=" + token + "; expires=" + date + "; path=/";
      window.location.replace("http://localhost:8800");
    }
    if (session.length > 2) {
      for (let i = Not_Login_div.length - 1; i >= 0; --i) {
        Not_Login_div[i].style.display = "none";
      }

      for (let i = Login_div.length - 1; i >= 0; --i) {
        Login_div[i].style.display = "flex";
      }
      document.getElementById("Form-comment").style.display = "flex";
    } else {
      document.getElementById("Form-comment").style.display = "none";
    }
  }
};

//--------------------------------------------

// Get the post btn to open the post modal
let postModalBtn = document.getElementById("postModalBtn");
// Get the post modal
let pmodal = document.getElementById("postModal");

// Get the <span> element that closes the login modal
let pspan = document.getElementsByClassName("pclose")[0];

// Get the post btn to post
let postBtn = document.getElementById("form-btn-post");

let imgInput = document.getElementById("img-input-id");

// close and reset the post modal
const Closepost = () => {
  pmodal.style.display = "none";
  document.getElementById("categories").value = "GO";
  document.getElementById("title").value = "";
  document.getElementById("post").value = "";
};

// When the user clicks the button it make a new post
postBtn.onclick = function (event) {
  event.preventDefault();
  let data = new FormData();
  let flag = false;
  let inputs = document.querySelectorAll('[name="option[]"]');
  let categories = "";
  for (let i = inputs.length - 1; i >= 0; --i) {
    if (inputs[i].checked) {
      categories = categories + " " + inputs[i].value;
      flag = true;
    }
  }

  if (flag) {
    if (document.getElementById("title").value === "") {
      alert("You need title ");
    } else {
      data.append("categories", categories);
      data.append("title", document.getElementById("title").value);
      data.append("post", document.getElementById("post").value);
      data.append("file", imgInput.files[0]);

      fetch("http://localhost:8800/post", {
        method: "POST",
        body: data,
      })
        .then(function (response) {
          return response.text();
        })
        .then(function (text) {
          //text is the server's response
          console.log(text);
          window.location.reload();
          Closepost();
        });
    }
  } else {
    alert("You much pick one of the categories");
  }
};

// When the user clicks the button, open the login modal
postModalBtn.onclick = function () {
  pmodal.style.display = "block";
};

// When the user clicks on <span> (x), close the login modal
pspan.onclick = () => Closepost();

// set a limit on the file upload
imgInput.onchange = function () {
  if (this.files[0].size > 20000000) {
    alert("File is too big!");
    this.value = "";
  }
};

//--------------------------------------------

let postid = "";

// When the user clicks the button to open the comment modal
const Comment = function (postID) {
  //cmodal.style.display = "block";
  postid = postID;
  document.getElementById(postID).classList.toggle("show");
};

// Close the dropdown if the user clicks outside of it
window.onclick = function (event) {
  if (!event.target.matches(".dropbtn")) {
    var dropdowns = document.getElementsByClassName("dropdown-content");
    var i;
    for (i = 0; i < dropdowns.length; i++) {
      var openDropdown = dropdowns[i];
      if (openDropdown.classList.contains("show")) {
        openDropdown.classList.remove("show");
      }
    }
  }
};

// When the user clicks anywhere outside of the modal, close it
window.onclick = function (event) {
  if (event.target == lmodal) {
    Closelogin();
  } else if (event.target == rmodal) {
    Closeregister();
  } else if (event.target == pmodal) {
    Closepost();
  }
};

//--------------------------------------------

function bindItemsInput() {
  let inputs = document.querySelectorAll('[name="option[]"]');
  let radioForCheckboxes = document.getElementById("radio-for-checkboxes");
  function checkCheckboxes() {
    let isAtLeastOneServiceSelected = false;
    for (let i = inputs.length - 1; i >= 0; --i) {
      if (inputs[i].checked) isAtLeastOneCheckboxSelected = true;
    }
    radioForCheckboxes.checked = isAtLeastOneCheckboxSelected;
  }
  for (let i = inputs.length - 1; i >= 0; --i) {
    inputs[i].addEventListener("change", checkCheckboxes);
  }
}
bindItemsInput(); // call in window
