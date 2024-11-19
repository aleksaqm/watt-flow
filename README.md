# watt-flow

## Authors

- Aleksa Perovic SV24/2021
- Danilo Cvijetic SV25/2021
- Vladimir Cornenki SV53/2021

Init App:

- When initializing app for the first time, in config.env file set RESTART field to True
- After that SuperAdmin account is created
- To be able to activate SuperAdmin account you will have to change account's default password
- That default password is stored in txt file : ./data/admin_password.txt
- To change that password write url: http://localhost:5173/superadmin, there you change default password
- Now you can log in with SuperAdmin account with: username: admin, password: your_new_password
- That is it. Use this steps when running the app for the first time or when you want to restart it for some reason.
- Warning: Set RESTART field in config.env back to false for future.
