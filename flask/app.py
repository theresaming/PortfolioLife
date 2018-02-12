from flask import Flask
from flask import flash, redirect, render_template, request, session, abort, jsonify, json
import httplib, urllib
import requests
import os
from sqlalchemy.orm import sessionmaker
from wtforms import Form, BooleanField, StringField, PasswordField, validators

from sqldb  import *

app = Flask(__name__)
app.secret_key = os.urandom(12)

app.config['SQLALCHEMY_DATABASE_URI'] = 'mysql://jd:jd2018@67.205.168.129/junior_design'
# db = SQLAlchemy(app)

class RegistrationForm(Form):
    username = StringField('Username', [validators.Length(min=4, max=25)])
    email = StringField('Email Address', [validators.Length(min=6, max=35)])
    password = PasswordField('New Password', [
        validators.DataRequired(),
        validators.EqualTo('confirm', message='Passwords must match')
    ])
    confirm = PasswordField('Repeat Password')


app.config['OAUTH_CREDENTIALS'] = {
    'facebook': {
        'id': '1987111774844195',
        'secret': '45f2751acb4510955a067be66c206e19'
    },
    'google': {
        'id': '500260639279-e8e8eniirbnm82m0sngtt0pogn86pbvd.apps.googleusercontent.com',
        'secret': 'QpD_M7SzNVpmMVNPXTrJmy1m'
    }
}

@app.route('/')
def home():
    if not session.get('logged_in'):
        return render_template('login.html')
    else:
        return "Hello!  <a href='/logout'>Logout</a>"

@app.route('/login', methods=['GET', 'POST'])
def do_admin_login():
    if request.method == 'POST':
        email = request.form['email']
        print email
        password = request.form['password']
        print password
        data = {
            "email": email,
            "password": password
        }
        jsonStr = json.dumps(data)
        r = requests.post('http://67.205.168.129:8080/user/login', jsonStr)
        jsonDict = json.loads(r.text)
        # print jsonDict['success']
        # print jsonDict['message']
        if jsonDict['success'] == True:
            session['logged_in'] = True
            return home()
        else:
            flash(jsonDict['message'])
    return home()

@app.route('/registration', methods=['GET','POST'])
def register():
    form = RegistrationForm(request.form)
    if request.method == 'POST':
        name = form.name.data
        password = form.password.data
        email = form.email.data
        data = {
            'email': email,
            'name': name,
            'password': password,
            }
        jsonStr = json.dumps(data)
        r = requests.post("http://67.205.168.129:8080/user/register", jsonStr)
        jsonDict = json.loads(r.text)
        # print "jsonDict: ", jsonDict
        # print "jsonDict['message']", jsonDict['message']
        # print "jsonDict['success']", jsonDict['success']
        if jsonDict['success'] == True:
            flash('Thanks for registering')
            return render_template('login.html')
        else:
            flash(jsonDict['message'])
    return render_template('registration.html', form=form)

@app.route("/logout")
def logout():
    session['logged_in'] = False
    return home()

if __name__ == "__main__":
    app.run(debug=True,host='0.0.0.0', port=4000)
