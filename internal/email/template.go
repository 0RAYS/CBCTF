package email

import "fmt"

// buildHTML generates a branded HTML email with a call-to-action button.
//
// Parameters:
//   - title      – displayed as the <h1> heading
//   - greeting   – plain text paragraph shown below the heading
//   - buttonText – label on the CTA button
//   - buttonURL  – href of the CTA button
//   - footerNote – small note below the button (e.g. expiry warning)
func buildHTML(title, greeting, buttonText, buttonURL, footerNote string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>%s</title>
</head>
<body style="margin:0;padding:0;background:#0a0a0a;font-family:'Segoe UI',system-ui,Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#0a0a0a;padding:48px 0;">
    <tr>
      <td align="center">
        <table width="560" cellpadding="0" cellspacing="0"
               style="background:#0f0f0f;border:1px solid #2a2a2a;border-radius:8px;overflow:hidden;max-width:560px;width:100%%;">

          <!-- Header -->
          <tr>
            <td style="padding:32px 40px 28px;border-bottom:1px solid #2a2a2a;">
              <p style="margin:0 0 12px;font-size:11px;letter-spacing:3px;text-transform:uppercase;color:#4d4d4d;font-weight:500;">CTF Platform</p>
              <h1 style="margin:0;font-size:22px;font-weight:600;color:#f9f9f9;line-height:1.3;">%s</h1>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding:32px 40px 8px;">
              <p style="margin:0 0 32px;font-size:14px;line-height:1.75;color:#b3b3b3;">%s</p>

              <!-- CTA Button -->
              <table cellpadding="0" cellspacing="0">
                <tr>
                  <td style="border-radius:6px;border:1px solid #597ef7;">
                    <a href="%s"
                       style="display:inline-block;padding:11px 28px;font-size:14px;font-weight:500;
                              color:#597ef7;text-decoration:none;letter-spacing:0.3px;border-radius:6px;
                              background:transparent;">
                      %s
                    </a>
                  </td>
                </tr>
              </table>

              <!-- Fallback link -->
              <p style="margin:24px 0 0;font-size:12px;color:#4d4d4d;line-height:1.6;word-break:break-all;">
                If the button does not work, copy and paste this link into your browser:<br/>
                <a href="%s" style="color:#597ef7;text-decoration:none;opacity:0.8;">%s</a>
              </p>
            </td>
          </tr>

          <!-- Footer note -->
          <tr>
            <td style="padding:20px 40px 32px;">
              <p style="margin:0;font-size:12px;color:#4d4d4d;line-height:1.6;">%s</p>
            </td>
          </tr>

          <!-- Divider + bottom bar -->
          <tr>
            <td style="border-top:1px solid #1a1a1a;padding:20px 40px;">
              <p style="margin:0;font-size:11px;color:#4d4d4d;text-align:center;letter-spacing:0.2px;">
                This is an automated message. Please do not reply to this email.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
		title,
		title,
		greeting,
		buttonURL,
		buttonText,
		buttonURL,
		buttonURL,
		footerNote,
	)
}

// buildTestHTML generates a "SMTP configuration test" branded HTML email.
func buildTestHTML(from string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>SMTP Configuration Test</title>
</head>
<body style="margin:0;padding:0;background:#0a0a0a;font-family:'Segoe UI',system-ui,Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#0a0a0a;padding:48px 0;">
    <tr>
      <td align="center">
        <table width="560" cellpadding="0" cellspacing="0"
               style="background:#0f0f0f;border:1px solid #2a2a2a;border-radius:8px;overflow:hidden;max-width:560px;width:100%%;">

          <!-- Header -->
          <tr>
            <td style="padding:32px 40px 28px;border-bottom:1px solid #2a2a2a;">
              <p style="margin:0 0 12px;font-size:11px;letter-spacing:3px;text-transform:uppercase;color:#4d4d4d;font-weight:500;">CTF Platform</p>
              <h1 style="margin:0;font-size:22px;font-weight:600;color:#f9f9f9;line-height:1.3;">SMTP Configuration Test</h1>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding:32px 40px 8px;">
              <!-- Status badge -->
              <table cellpadding="0" cellspacing="0" style="margin-bottom:24px;">
                <tr>
                  <td style="border-radius:6px;border:1px solid #2a2a2a;background:#0a0a0a;padding:10px 16px;">
                    <span style="font-size:13px;color:#b3b3b3;font-weight:500;">&#10003;&nbsp;&nbsp;Connection verified</span>
                  </td>
                </tr>
              </table>

              <p style="margin:0 0 24px;font-size:14px;line-height:1.75;color:#b3b3b3;">
                This is a test email to confirm your SMTP configuration is working correctly. No action is required.
              </p>

              <!-- Sender info -->
              <table width="100%%" cellpadding="0" cellspacing="0"
                     style="background:#0a0a0a;border:1px solid #2a2a2a;border-radius:6px;">
                <tr>
                  <td style="padding:14px 20px;">
                    <span style="font-size:11px;color:#4d4d4d;display:block;margin-bottom:4px;letter-spacing:0.5px;text-transform:uppercase;">Sender address</span>
                    <span style="font-size:13px;color:#8a8a8a;font-family:'Courier New',monospace;">%s</span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Spacer -->
          <tr><td style="height:32px;"></td></tr>

          <!-- Divider + bottom bar -->
          <tr>
            <td style="border-top:1px solid #1a1a1a;padding:20px 40px;">
              <p style="margin:0;font-size:11px;color:#4d4d4d;text-align:center;letter-spacing:0.2px;">
                This is an automated test message from the CTF Platform admin panel.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, from)
}
