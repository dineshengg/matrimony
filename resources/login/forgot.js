const errorDiv = document.getElementById('error');
const errorField = document.getElementById('error-field');

const submitBtn = document.querySelector(".submit");
submitBtn.addEventListener("click", function(event){
  event.preventDefault();
  
  const emailidinput = document.getElementById("email");
  const emailid = emailidinput.value.trim();
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  
  if (emailid === "" || !emailPattern.test(emailid)) {
    showError(emailidinput, "Please enter a valid email address.");
    emailidinput.focus();
    return;
  }

  // Clear any previous errors
  if (errorDiv) {
    errorField.style.display = "none";
    errorDiv.hidden = true;
  }

  // Submit the form to backend
  submitForgotPassword(emailid);
});

async function submitForgotPassword(email) {
  try {
    const formData = new FormData();
    formData.append('email', email);

    const response = await fetch('/api/noauth/forgot-password', {
      method: 'POST',
      body: formData
    });

    if (response.ok) {
      const message = await response.text();
      showSuccess("Password reset email sent successfully! Please check your inbox.");
      console.log('Success:', message);
      // Optionally redirect after 3 seconds
      setTimeout(() => {
        window.location.href = '/static/login/login.html';
      }, 3000);
    } else {
      const errorMsg = await response.text();
      showError(null, errorMsg || "Failed to process request. Please try again.");
      console.error('Error:', errorMsg);
      //optionally redirect to blank html page
      setTimeout(() => {
        window.location.href = '/static/login/blankuser.html';
      }, 3000);
    }
  } catch (error) {
    console.error('Fetch error:', error);
    showError(null, "Network error. Please check your connection and try again.");
  }
}

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
      errorDiv.style.color = "red";
    }
  }
}

function showSuccess(message) {
  if (errorDiv) {
    errorField.style.display = "block";
    errorDiv.textContent = message;
    errorDiv.hidden = false;
    errorDiv.style.color = "green";
  }
}