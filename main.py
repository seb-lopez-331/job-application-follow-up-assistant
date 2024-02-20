from google.oauth2 import service_account
from dotenv import load_dotenv
import pandas as pd
import gspread
import json
import os

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

if __name__ == "__main__":
    # Load environment
    load_dotenv()

    # Set up a client with scoped credentials
    client = setup_client()

    # Read spreadsheets data into a DataFrame
    spreadsheet_url = os.environ['SPREADSHEET_URL']
    records_df = read_spreadsheets_into_dataframe(spreadsheet_url)
    
    print(records_df)

