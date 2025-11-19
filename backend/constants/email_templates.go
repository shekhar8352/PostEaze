package constants

const (
	NotificationEmailSubject = "New Notification from PostEaze"
	TeamInviteEmailSubject   = "You're invited to join a team on PostEaze"
)

const NotificationEmailBody = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { padding: 20px; }
        .content { background-color: #f9f9f9; padding: 20px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>New Notification</h2>
        <div class="content">
            <p>%s</p>
        </div>
        <p>Best,<br>The PostEaze Team</p>
    </div>
</body>
</html>
`

const TeamInviteEmailBody = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { padding: 20px; }
        .button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Team Invitation</h2>
        <p>You have been invited to join a team on PostEaze.</p>
        <p>Click the button below to accept the invitation:</p>
        <a href="%s" class="button">Join Team</a>
        <p>If you didn't expect this invitation, you can ignore this email.</p>
        <p>Best,<br>The PostEaze Team</p>
    </div>
</body>
</html>
`
