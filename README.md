# job-application-follow-up-assistant
A script that regularly checks my Google Sheets, where I track the status of job applications, and sends me email reminders to follow up.

## Prerequisites
- Python >= 3.12.2
- crontab

## Environment Variables

This project utilizes environment variables for configuration. Inside the `.env` file, you will find the necessary files for setting up your environment.

This repository offers a `.env.template` file with an overview of the variables that are placed inside of the `.env` file.

If you require access to the API key or any other sensitive information, please contact me directly, and I will provide you with the `.env` file.

## How to use

First, we would like to install our required dependencies using pip. We can do this by running the following:

`./pre-script.sh`

Secondly run crontab

`crontab -e`

Inside your crontab, add something like:

`0 9 * * * python3 main.py`

at the end of the crontab file.

This will run your script every day at 9:00 AM. If you wish to configure it for another time or frequency, experiment with the crontab guru https://crontab.guru/
for different cron commands.

## Important Note

Configure your own app password on gmail. Put that in your .env
