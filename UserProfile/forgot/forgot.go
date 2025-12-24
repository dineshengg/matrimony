package forgot

import (
	"fmt"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/dineshengg/matrimony/common/utils"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func ForgotRoutingFunctions(router *routing.RouteGroup) {
	//no authentication flow where jwt token check is not required
	noauth := router.Group("/noauth")
	noauth.Post("/forgot-password", forgotPasswordHandler)
	//show the reset password page - GET because it's accessed via email link
	noauth.Get("/reset-link", resetLinkHandler)
	noauth.Post("/reset-password", resetPasswordHandler)	
}


func forgotPasswordHandler(ctx *routing.Context) error {
	//1. Email id is present => send email to reset password
	//2. Email id is not present => show error message that email id doesnt exists and provide a link to create account
	//3. Any other error happens => show error message its not your fault please come after some time

	Email := string(ctx.FormValue("email"))

	// Validate if email exists
	matid, err := checkIfUserExistsAlready(Email, "")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte("Failed to check user existence"))
		log.Debugf("Failed to check user existence - %v", err)
		return fmt.Errorf("Failed to check user existence - %v", err)
	}
	if matid != "" {
		// TODO: Send email to this email id
		ctx.Write([]byte(fmt.Sprintf("Email id exists, email was sent to this id %s with reset password link", Email)))
		//ctx.Response.Header.Set("Location", "/services/sent-email")
		//store in the forgot table and the no of times reset password was requested
		resetPassword(Email, matid)
		err = utils.LogEmailEvent(utils.PasswordResetEmail, Email, matid)
		if err != nil {
			log.Errorf("failed to send password reset email to user %s:%v", Email, err)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusSeeOther)
		ctx.Write([]byte(fmt.Sprintf("Email id doesnt exists, redirecting to create account page .....")))
	}
	return nil
}
func resetLinkHandler(ctx *routing.Context) error {
	// Get GUID from query parameter
	fmt.Printf("Reset link handler called\n")
	GUID := string(ctx.QueryArgs().Peek("guid"))
	var result struct {
		Email   string
		ResetAt time.Time
	}
	if GUID == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("Reset password GUID is missing")
		return fmt.Errorf("Reset password GUID is missing")
	}

	// Check if the GUID exists in the forgot table
	query := `SELECT email, reset_at FROM forgot WHERE guid = $1`
	err := utils.GetDB().Raw(query, GUID).Scan(&result).Error
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debugf("Invalid or expired reset password link - %v", err)
		return fmt.Errorf("Invalid or expired reset password link - %v", err)
	}

	// Check if the reset link has expired (reset_at should be in the future, within 24 hours from now)
	if time.Now().After(result.ResetAt)  {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("Reset password link has expired")
		return fmt.Errorf("Reset password link has expired")
	}

	fmt.Printf("GUID is valid, rendering reset password page for email %s\n", result.Email)
	if page, err := utils.RenderTemplatePage(ctx, "resources/login/newpassword.html", map[string]interface{}{"GUID": GUID}); err == nil {
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.Write(page)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return nil
	} else {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Error rendering reset password page %v", err)))
		log.Debugf("Error rendering reset password page - %v", err)
		return fmt.Errorf("Error rendering reset password page - %v", err)
	}

	return nil
}

func resetPasswordHandler(ctx *routing.Context) error {
	//1. Guid is valid and present => allow user to reset password
	//2. Guid is not valid or not present => show error message that reset link is invalid
	//3. Any other error happens => show error message its not your fault please come after some time

	// Get guid from query parameter or form value
	guid := string(ctx.QueryArgs().Peek("guid"))
	if guid == "" {
		guid = string(ctx.FormValue("guid"))
	}
	var userEmail string
	if guid == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Reset password GUID is missing"))
		log.Debug("Reset password GUID is missing")
		return fmt.Errorf("Reset password GUID is missing")
	}

	// Check if the GUID exists in the forgot table
	query := `SELECT email FROM forgot WHERE guid = $1`
	err := utils.GetDB().Raw(query, guid).Scan(&userEmail).Error
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Invalid or expired reset password link"))
		log.Debugf("Invalid or expired reset password link - %v", err)
		return fmt.Errorf("Invalid or expired reset password link - %v", err)
	}

	// Allow user to reset password
	newPassword := string(ctx.FormValue("newpassword"))
	confirmPassword := string(ctx.FormValue("confirmpassword"))

	if newPassword == "" || confirmPassword == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("New password and confirm password are required"))
		log.Debug("New password and confirm password are required")
		return fmt.Errorf("New password and confirm password are required")
	}

	if newPassword != confirmPassword {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("New password and confirm password do not match"))
		log.Debug("New password and confirm password do not match")
		return fmt.Errorf("New password and confirm password do not match")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte("Failed to hash new password"))
		log.Debugf("Failed to hash new password - %v", err)
		return fmt.Errorf("Failed to hash new password - %v", err)
	}

	// Update the user's password in the profiles table
	updateQuery := `UPDATE profiles SET password = $1 WHERE email = $2`
	err = utils.GetDB().Exec(updateQuery, string(hashedPassword), userEmail).Error
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte("Failed to update password"))
		log.Debugf("Failed to update password - %v", err)
		return fmt.Errorf("Failed to update password - %v", err)
	}

	// Delete the used reset token from forgot table
	deleteQuery := `DELETE FROM forgot WHERE guid = $1`
	err = utils.GetDB().Exec(deleteQuery, guid).Error
	if err != nil {
		log.Errorf("Failed to delete reset token for guid %s: %v", guid, err)
		// Don't return error here since password was already reset successfully
	}

	ctx.Write([]byte("Password has been reset successfully"))
	return nil

}