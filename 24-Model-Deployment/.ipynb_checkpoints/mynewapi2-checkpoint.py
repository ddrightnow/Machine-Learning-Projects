#print('Hello')
from flask import Flask,request,jsonify
import joblib
import pandas as pd

#CREATE FLASK APP
app = Flask(__name__)

#CONNECT POST API CALL --> PREDICT() FUNCTION #http://localhost:5000/predict

#predict would define what the last section is on url
#Post as a method means when you get to the url, you can post some information
@app.route('/predict',methods=['POST']) 
def predict():
    
    #GET JSON REQUEST
    feat_data = request.json
    
    #CONVERT JSON TO PANDAS DF & ENSURE MATCHES WITH COLUMN NAMES
    df = pd.DataFrame(feat_data)
    df = df.reindex(columns = col_names)
    
    #MAKE THE PREDICTION
    prediction = list(model.predict(df))
    
    return jsonify ({'prediction':str(prediction)}) #prediction


#LOAD MY MODEL AND COLUMN NAMES
if __name__ == '__main__':
    
    model = joblib.load('final_model.pkl')
    col_names = joblib.load('column_names.pkl')
    
    app.run(debug=True)