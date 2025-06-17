const slidePage = document.querySelector(".slide-page");
const nextBtnFirst = document.querySelector(".firstNext");
const prevBtnSec = document.querySelector(".prev-1");
const nextBtnSec = document.querySelector(".next-1");
const prevBtnThird = document.querySelector(".prev-2");
const submitBtn = document.querySelector(".submit");
const resetBtn = document.querySelector(".reset");
const progressText = document.querySelectorAll(".step p");
const progressCheck = document.querySelectorAll(".step .check");
const bullet = document.querySelectorAll(".step .bullet");
let current = 1;

const firstnameinput = document.getElementById("firstname");
const secondnameinput = document.getElementById("secondname");
const emailidinput = document.getElementById("email");
const phonenoinput = document.getElementById("phone");
const dateinput = document.getElementById("dob");
const passwordInput = document.getElementById("password");
const confirmPasswordInput = document.getElementById("confirmpassword");

[firstnameinput, secondnameinput, emailidinput, phonenoinput, dateinput, passwordInput, confirmPasswordInput].forEach(function(input) {
  input.addEventListener("input", function() {
    input.style.borderColor = "";
    input.placeholder = "";
  });
});

nextBtnFirst.addEventListener("click", function(event){
  event.preventDefault();
  //TODO- validate the first name, second name and emailid and phone no.
  //alert("debug1");
  showError("Please enter a valid email address.");
  let valid = true;
  firstname = firstnameinput.value.trim();
  secondname = secondnameinput.value.trim();
  if (firstname === "") {
    showError(firstnameinput, "Please enter your first");
    firstnameinput.focus();
    valid = false;
    return;
  }
  if (secondname === "" ) {
    secondnameinput.focus();
    showError(secondnameinput, "Please enter your second name");
    valid = false;
    return;
  }
  
  const emailid = emailidinput.value.trim();
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (emailid === "" || !emailPattern.test(emailid)) {
    //alert("email id empty");
    showError(emailidinput, "Please enter a valid email address.");
    emailidinput.focus();
    valid = false;
    return;
  }

  
  const phonePattern = /^[0-9]{8,15}$/;
  phoneno = phonenoinput.value.trim();

  if (phoneno === "" || !phonePattern.test(phoneno)){
    showError(phonenoinput, "Please enter your correct phone number")
    valid = false;
    return;
  }

  slidePage.style.marginLeft = "-25%";
  bullet[current - 1].classList.add("active");
  progressCheck[current - 1].classList.add("active");
  progressText[current - 1].classList.add("active");
  current += 1;
  
});
nextBtnSec.addEventListener("click", function(event){
  event.preventDefault();
  let valid = true;
  
  const dateValue = dateinput ? dateinput.value.trim() : "";
  const datePattern = /^(0[1-9]|1[0-2])\/(0[1-9]|[12][0-9]|3[01])\/(19|20)\d\d$/;

  if (dateinput && (dateValue === "" || !datePattern.test(dateValue))) {
    showError(dateinput, "Please enter a valid date (dd/mm/yyyy)");
    dateinput.focus();
    valid = false;
    //return;
  }
  slidePage.style.marginLeft = "-50%";
  bullet[current - 1].classList.add("active");
  progressCheck[current - 1].classList.add("active");
  progressText[current - 1].classList.add("active");
  current += 1;
});
submitBtn.addEventListener("click", function(){
  let valid = true;

  // Password validation regex: at least 8 chars, includes allowed special chars
  const passwordPattern = /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@;#$%^&*])[A-Za-z\d@;#$%^&*]{8,}$/;
  const password = passwordInput.value.trim();
  const confirmPassword = confirmPasswordInput.value.trim();

  //validate the passwords
  if (password === "" || !passwordPattern.test(password)) {
    showError(passwordInput, "Password must be at least 8 characters and include A-Z, a-z, 0-9, and @;#$%^&*");
    passwordInput.focus();
    valid = false;
    event.preventDefault();
    return;
  }

  // Validate confirm password
  if (confirmPassword === "") {
    showError(confirmPasswordInput, "Please confirm your password");
    confirmPasswordInput.focus();
    valid = false;
    event.preventDefault();
    return;
  }

  if (password !== confirmPassword) {
    showError(confirmPasswordInput, "Passwords do not match");
    confirmPasswordInput.focus();
    valid = false;
    event.preventDefault();
    return;
  }

  bullet[current - 1].classList.add("active");
  progressCheck[current - 1].classList.add("active");
  progressText[current - 1].classList.add("active");
  current += 1;
});

resetBtn.addEventListener("click", function(){
	event.preventDefault();
	bullet[current-2].classList.remove("active");
	bullet[current-3].classList.remove("active");
	progressCheck[current-2].classList.remove("active");
	progressCheck[current-3].classList.remove("active");
	progressCheck[current-2].classList.remove("active");
	progressCheck[current-3].classList.remove("active");
	location.reload();
	current=1
})

prevBtnSec.addEventListener("click", function(event){
  event.preventDefault();
  slidePage.style.marginLeft = "0%";
  bullet[current - 2].classList.remove("active");
  progressCheck[current - 2].classList.remove("active");
  progressText[current - 2].classList.remove("active");
  current -= 1;
});
prevBtnThird.addEventListener("click", function(event){
  event.preventDefault();
  slidePage.style.marginLeft = "-25%";
  bullet[current - 2].classList.remove("active");
  progressCheck[current - 2].classList.remove("active");
  progressText[current - 2].classList.remove("active");
  current -= 1;
});

//show error message in this form
function showError(input, message) {
  if (input && input.tagName === "INPUT") {
    input.value = "";
    input.placeholder = message;
    input.style.borderColor = "red";
    input.style.outline = "none";
    input.style.fontSize = "14px";
    input.focus();
  }
}
