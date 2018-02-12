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

@app.route('/login', methods=['POST'])
def do_admin_login():
    POST_USERNAME = str(request.form['username'])
    POST_PASSWORD = str(request.form['password'])

    Session = sessionmaker(bind=engine)
    s = Session()
    query = s.query(User).filter(User.username.in_([POST_USERNAME]), User.password.in_([POST_PASSWORD]) )
    result = query.first()
    if result:
        session['logged_in'] = True
    else:
        flash('wrong password!')
    return home()

@app.route('/registration', methods=['GET','POST'])
def register():
    form = RegistrationForm(request.form)
    if request.method == 'POST':
        username = form.username.data
        password = form.password.data
        email = form.email.data
        data = {
            'email': email,
            'name': username,
            'password': password,
            }
        jsonStr = json.dumps(data)
        r = requests.post("http://67.205.168.129:8080/user/register", jsonStr)
        jsonDict = json.loads(r.text)
        print "jsonDict: ", jsonDict
        print "jsonDict['message']", jsonDict['message']
        print "jsonDict['success']", jsonDict['success']
        if jsonDict['success'] == True:
            flash('Thanks for registering')
            return render_template('login.html')
        else:
            flash(jsonDict['message'])
    return render_template('registration.html', form=form)

@app.route('/api/register', methods=['POST'])
def reg():
    content = request.get_json()
    print (content)

# @app.route('/authorize/facebook')
# def oauth_authorize(provider):
#     if not current_user.is_anonymous():
#         return redirect(url_for('login'))
#     oauth = OAuthSignIn.get_provider(provider)
#     return oauth.authorize()
#
# @app.route('/callback/facebook')
# def oauth_callback(provider):
#     if not current_user.is_anonymous():
#         return redirect(url_for('index'))
#     oauth = OAuthSignIn.get_provider(provider)
#     social_id, username, email = oauth.callback()
#     if social_id is None:
#         flash('Authentication failed.')
#         return redirect(url_for('logout'))
#     user = User.query.filter_by(social_id=social_id).first()
#     if not user:
#         user = User(social_id=social_id, nickname=username, email=email)
#         db.session.add(user)
#         db.session.commit()
#     login_user(user, True)
#     return redirect(url_for('login'))

@app.route("/logout")
def logout():
    session['logged_in'] = False
    return home()

if __name__ == "__main__":
    app.run(debug=True,host='0.0.0.0', port=4000)
