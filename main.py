%%capture
!pip install gspread
from google.oauth2 import service_account
import pandas as pd
import gspread
import json
import os

google_json = os.environ["GOOGLE_JSON"]