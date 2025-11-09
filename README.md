### EMAIL SMTP App
This is a simple program that allows you to programmatically email multiple receipients in one go. 

## Prerequisite
You need to get "App Password" from Google. Follow following steps:
- Set up 2FA for you google account. 
- Navigate to *Settings*, then to *Security* and search for *App passwords*
- Create App password and save it in your .env file

Please create your own .env file locally and populate following values:
SENDER_EMAIL= `<your email>` \
APP_PASSWORD=`<your app password that you generated>`

SMTP_HOST= `"smtp.gmail.com"`\
SMTP_PORT= `587`

You can modify recipients and populate email body by modifying `main.go` file. 