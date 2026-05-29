package email

import "fmt"

// buildTestHTML generates a simple "SMTP test" branded HTML email.
func buildTestHTML(from string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>SMTP Test Email</title>
</head>
<body style="margin:0;padding:0;background:#0f0f0f;font-family:'Segoe UI',Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#0f0f0f;padding:40px 0;">
    <tr>
      <td align="center">
        <table width="560" cellpadding="0" cellspacing="0"
               style="background:#161616;border:1px solid #2a2a2a;border-radius:12px;overflow:hidden;max-width:560px;">
          <!-- Header -->
          <tr>
            <td style="background:#1a2e1a;padding:32px 40px;border-bottom:2px solid #52c41a;">
              <p style="margin:0;font-size:13px;letter-spacing:3px;text-transform:uppercase;color:#52c41a;font-weight:600;">CTF Platform</p>
              <h1 style="margin:8px 0 0;font-size:24px;font-weight:700;color:#ffffff;">SMTP Configuration Test</h1>
            </td>
          </tr>
          <!-- Body -->
          <tr>
            <td style="padding:36px 40px;">
              <div style="display:inline-block;padding:10px 20px;background:#1a2e1a;border:1px solid #52c41a;
                          border-radius:8px;margin-bottom:24px;">
                <span style="font-size:22px;margin-right:8px;">✓</span>
                <span style="font-size:15px;font-weight:600;color:#52c41a;vertical-align:middle;">Connection Successful</span>
              </div>
              <p style="margin:0 0 16px;font-size:15px;line-height:1.7;color:#c8c8c8;">
                This is a test email sent from your CTF Platform to verify that your SMTP configuration is working correctly.
              </p>
              <table width="100%%" cellpadding="0" cellspacing="0"
                     style="background:#0f0f0f;border:1px solid #2a2a2a;border-radius:8px;padding:16px;">
                <tr>
                  <td style="padding:8px 16px;">
                    <span style="font-size:12px;color:#666666;display:block;margin-bottom:2px;">Sender Address</span>
                    <span style="font-size:14px;color:#c8c8c8;font-family:monospace;">%s</span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <!-- Bottom bar -->
          <tr>
            <td style="background:#0f0f0f;padding:20px 40px;border-top:1px solid #2a2a2a;">
              <p style="margin:0;font-size:11px;color:#444444;text-align:center;">
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
// title       – displayed as the <h1> heading
// greeting    – plain text paragraph shown below the heading
// buttonText  – label on the CTA button
// buttonURL   – href of the CTA button
// footerNote  – small note below the button (e.g. expiry warning)
func buildHTML(title, greeting, buttonText, buttonURL, footerNote string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>%s</title>
</head>
<body style="margin:0;padding:0;background:#0f0f0f;font-family:'Segoe UI',Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background:#0f0f0f;padding:40px 0;">
    <tr>
      <td align="center">
        <table width="560" cellpadding="0" cellspacing="0"
               style="background:#161616;border:1px solid #2a2a2a;border-radius:12px;overflow:hidden;max-width:560px;">
          <!-- Header -->
          <tr>
            <td style="background:#1a1a2e;padding:32px 40px;border-bottom:2px solid #597ef7;">
              <p style="margin:0;font-size:13px;letter-spacing:3px;text-transform:uppercase;color:#597ef7;font-weight:600;">CTF Platform</p>
              <h1 style="margin:8px 0 0;font-size:24px;font-weight:700;color:#ffffff;">%s</h1>
            </td>
          </tr>
          <!-- Body -->
          <tr>
            <td style="padding:36px 40px;">
              <p style="margin:0 0 28px;font-size:15px;line-height:1.7;color:#c8c8c8;">%s</p>
              <!-- CTA Button -->
              <table cellpadding="0" cellspacing="0">
                <tr>
                  <td style="border-radius:8px;background:#597ef7;">
                    <a href="%s"
                       style="display:inline-block;padding:14px 32px;font-size:15px;font-weight:600;
                              color:#ffffff;text-decoration:none;letter-spacing:0.5px;border-radius:8px;">
                      %s
                    </a>
                  </td>
                </tr>
              </table>
              <!-- Fallback link -->
              <p style="margin:24px 0 0;font-size:12px;color:#666666;word-break:break-all;">
                If the button does not work, copy and paste this link into your browser:<br/>
                <a href="%s" style="color:#597ef7;text-decoration:none;">%s</a>
              </p>
            </td>
          </tr>
          <!-- Footer note -->
          <tr>
            <td style="padding:0 40px 28px;">
              <p style="margin:0;font-size:12px;color:#555555;line-height:1.6;">%s</p>
            </td>
          </tr>
          <!-- Bottom bar -->
          <tr>
            <td style="background:#0f0f0f;padding:20px 40px;border-top:1px solid #2a2a2a;">
              <p style="margin:0;font-size:11px;color:#444444;text-align:center;">
                This is an automated message. Please do not reply directly to this email.
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
