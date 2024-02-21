from google.oauth2 import service_account
from dotenv import load_dotenv
from datetime import datetime, timedelta
from email.message import EmailMessage
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
import pandas as pd
import gspread
import json
import os
import smtplib, ssl

def setup_client():
    """
    Set up and authenticate a Google Sheets client using service account credentials.

    This function loads the GOOGLE_JSON environment variable, which should contain the
    service account credentials in JSON format. It then creates a credentials object
    from the service account credentials JSON, assigns the appropriate API scope to
    the credentials, and authorizes the client using the scoped credentials.

    Returns:
        gspread.Client: A client object authenticated with the service account credentials
                         and scoped for accessing Google Sheets.
                         
    Raises:
        KeyError: If the GOOGLE_JSON environment variable is not set or if it does not contain
                  valid service account credentials JSON.
        json.JSONDecodeError: If the content of the GOOGLE_JSON environment variable is not
                              valid JSON.
    """
    # Load the GOOGLE_JSON environment variable and store it
    google_json = os.environ["GOOGLE_JSON"]

    # Create a credentials object from the service account credentials json
    service_account_info = json.loads(google_json)
    credentials = service_account.Credentials.from_service_account_info(service_account_info)

    # Assign the right API scope to the credentials so that it has the right permissions to
    # interface with the Google Sheet spreadsheet
    scope = ['https://spreadsheets.google.com/feeds','https://www.googleapis.com/auth/drive']
    creds_with_scope = credentials.with_scopes(scope)

    # Create a new client with these scoped credentials
    client = gspread.authorize(creds_with_scope)
    return client

def read_spreadsheets_into_dataframe(url):
    """
    Read data from a Google Sheets spreadsheet into a pandas DataFrame.

    Args:
        url (str): The URL of the Google Sheets spreadsheet to read from.

    Returns:
        pandas.DataFrame: A DataFrame containing the data from the specified Google Sheets spreadsheet.

    Raises:
        gspread.exceptions.APIError: If there is an error accessing the Google Sheets API.
        gspread.exceptions.SpreadsheetNotFound: If the specified spreadsheet URL does not exist.
        gspread.exceptions.WorksheetNotFound: If the specified worksheet (tab) does not exist in the spreadsheet.
        gspread.exceptions.NoValidUrlKeyFound: If the provided URL does not contain a valid Google Sheets key.
    """
    # Create a new sheet instance from a spreadsheet URL
    spreadsheet = client.open_by_url(url)

    # Zoom in on the worksheet (the tab in the Google spreadsheet) we want to get the data from
    worksheet = spreadsheet.get_worksheet(0)

    # Get all the records as a JSON
    records_data = worksheet.get_all_records()

    # Convert the JSON records into a pandas DataFrame
    records_df = pd.DataFrame.from_dict(records_data)
    return records_df

def find_people_to_follow_up(df):
    """
    Find people to follow up based on specific criteria.

    Args:
        df (pandas.DataFrame): Input DataFrame containing job application data.

    Returns:
        pandas.DataFrame: DataFrame containing rows that meet the criteria for follow-up.

    Notes:
        - The function assumes that the input DataFrame has columns representing various attributes of job applications,
          such as 'Last_Spoken_On', 'Status', 'Recruiter', and 'Hiring_Manager'.
        - It filters the DataFrame based on the following criteria:
          1. 'Last_Spoken_On' matches the date one week ago.
          2. 'Status' is empty (indicating no recent update).
          3. Either 'Recruiter' or 'Hiring_Manager' is not empty (indicating a contact person is available).
    """
    # Replace blank spaces with0 '_'
    df.columns = [column.replace(" ", "_") for column in df.columns] 

    # After converting the first row to column names, the index of the DataFrame will be
    # off by one. To fix this, we use the .reset_index method
    df = df.reset_index(drop=True)

    # Grab today's date to determine the date it was a week ago
    now = datetime.now()
    week_ago = now - timedelta(days=7)

    # Format 'week_ago' to match the date format (mm/dd/yyyy) in the dataset
    fmt_week_ago = week_ago.strftime("%m/%d/%y")
    print(df.columns)

    # Find all jobs where your last communications with the hiring team was one week ago
    # This way, we provide a notification that one week has been up, and it's therefore
    # a good time to reach out.
    df.query(f'Last_Spoken_On == {fmt_week_ago} and Status == "" and (Recruiter != "" or Hiring_Manager != "")', inplace=True)

    return df

def send_reminder_emails(df):
    """
    Send reminder emails to follow up on job applications.

    Args:
        df (pandas.DataFrame): DataFrame containing job application data.

    Notes:
        - The function assumes that the input DataFrame has columns representing various attributes of job applications,
          such as 'Recruiter', 'Hiring_Manager', 'Company', and 'Role'.
        - It pulls out email credentials from environment variables (EMAIL_ADDRESS and APP_PASSWORD).
        - It sends reminder emails for each row in the DataFrame, using the recruiter or hiring manager as the contact.
        - The subject and content of the email are personalized based on the job application details.
    """
    # Configure SSL port
    port = 465

    # Pull out your email credentials from .env file
    email = os.environ['EMAIL_ADDRESS']
    password = os.environ['APP_PASSWORD']

    # Create a secure SSL context
    context = ssl.create_default_context()
    with smtplib.SMTP_SSL("smtp.gmail.com", port, context=context) as server:
        server.login(email, password)

        # Send email reminders to follow up
        # Note: the reminders will be sent by your own email address
        for index, row in df.iterrows():
            # Parse the information
            recruiter = row.iloc[1]
            hiring_manager = row.iloc[4]
            company = row.iloc[5]
            role = row.iloc[6]

            # Contact will either be the recruiter or hiring manager
            contact = recruiter or hiring_manager

            # Construct the message metadata
            message = MIMEMultipart()
            message["Subject"] = f'Follow-Up Reminder: Job Application Status of {role}, {company}'
            message["From"] = "Job Application Follow-Up Assistant"
            message["To"] = email

            # Create the plain text message
            text = """
            Hey %s,

            Hoping all is well! According to our databases, it's been a week since you've last
            contacted %s. Please follow up with them on the status of your application to
            %s at %s.
            
            Thanks,
            Your Faithful Follow-Up Assistant
            """ % (email, contact, role, company)
            
            # Turn the text into plain a MIMEText object
            part = MIMEText(text, "plain")

            # Add the part to the message
            message.attach(part)

            # Send the mail
            server.login(email, password)
            server.sendmail(email, email, message.as_string())

if __name__ == "__main__":
    """
    Main block of the script responsible for executing the workflow.

    Notes:
        - This block is executed when the script is run as the main program.
        - It loads environment variables using dotenv.
        - It sets up a client with scoped credentials for accessing Google Sheets.
        - It reads data from a specified spreadsheet URL into a DataFrame.
        - It filters the DataFrame to identify people to follow up with based on specific criteria.
        - It sends reminder emails to the identified contacts.
    """
    # Load environment
    load_dotenv()

    # Set up a client with scoped credentials
    client = setup_client()

    # Read spreadsheets data into a DataFrame
    spreadsheet_url = os.environ['SPREADSHEET_URL']
    records_df = read_spreadsheets_into_dataframe(spreadsheet_url)
    
    # Filter out data to only return the people we might wish to follow up with
    filtered_df = find_people_to_follow_up(records_df)

    # Send reminder emails
    send_reminder_emails(filtered_df)

