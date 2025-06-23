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
const errorDiv = document.getElementById('error');
const errorField = document.getElementById('error-field');
errorField.style.display = "none";

[firstnameinput, secondnameinput, emailidinput, phonenoinput, dateinput, passwordInput, confirmPasswordInput].forEach(function(input) {
  input.addEventListener("input", function() {
    input.style.borderColor = "";
    input.placeholder = "";
  });
});

nextBtnFirst.addEventListener("click", function(event){
  event.preventDefault();
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


  
  const phonePattern = /^(?:\+91)?\d{10}$/;
  phoneno = phonenoinput.value.trim();

  if (phoneno === "" || !phonePattern.test(phoneno)){
    showError(phonenoinput, "Please enter your correct phone number")
    valid = false;
    return;
  }

  validateEmailAndPhone(emailid, phoneno).then(result => {

  errorField.style.display = "none";
  errorDiv.textContent = " ";
  errorDiv.hidden = true;

  if (result.error) {
    console.error("Error validating email and phone:", result.error);
    showError(null, 'Error validating email and phone: check internet connection');
    return;
  }

  if (result.email) {
    showError(emailidinput, "Email already registered.");
    emailidinput.focus();
    return;
  }
  if (result.phone) {
    showError(phonenoinput, "Phone number already registered.");
    phonenoinput.focus();
    return;
  }
  // Proceed to next step if both are unique
  slidePage.style.marginLeft = "-25%";
  bullet[current - 1].classList.add("active");
  progressCheck[current - 1].classList.add("active");
  progressText[current - 1].classList.add("active");
  current += 1;
});
});

nextBtnSec.addEventListener("click", function(event){
  event.preventDefault();
  let valid = true;
  
  const dateValue = dateinput ? dateinput.value.trim() : "";
  const datePattern = /^\d{4}\/\d{2}\/\d{2}$/;

  if( dateValue === ""){ 
    showError(dateinput, "Please give valid DOB");
    dateinput.focus();
    valid = false;
    return;
  }

  console.log("Date value entered:", dateValue);
  const parts = dateValue.split('-').map(Number);
  const year = parts[0];
  const month = parts[1];
  const day = parts[2];

  console.log("Parsed date components:", { year, month, day });

  // Check if the created Date object's components match the input components
  // This handles invalid dates like "2023/02/30" (February 30th)
  if ((year <= 1950 || year > 2050) || (month < 1 || month >12 ) || (day < 1 || day >= 31)) {
    console.log(dateValue);
    showError(dateinput, "Please enter a valid date as per calendar (dd/mm/yyyy)");
    dateinput.focus();
    valid = false;
    return ;
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
  } else {
    if (errorDiv) {
      console.log("Error Div is present");
      errorField.style.display = "block";
      errorDiv.textContent = message;
      errorDiv.hidden = false;
    }
  }
}

// function to validate if email id and phone is already registered if present show error message
async function validateEmailAndPhone(emailid, phoneno) {
  try {
    const formData = new URLSearchParams();
    formData.append('email', emailid);
    formData.append('phone', phoneno);

    const response = await fetch('/api/new-profile/validate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      },
      body: formData.toString()
    });

    if (!response.ok) {
      showError(null, await response.text());
    }

    const result = await response.json();
    return result;
  } catch (error) {
    return { error: true };
  }
}