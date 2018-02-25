from flask import Flask, flash, render_template, request, session, json, redirect, make_response
import requests
import os
from werkzeug import secure_filename
from werkzeug.datastructures import FileStorage
from wtforms import Form, StringField, PasswordField, validators

app = Flask(__name__)
app.secret_key = os.urandom(12)

app.config['SQLALCHEMY_DATABASE_URI'] = 'mysql://jd:jd2018@67.205.168.129/junior_design'

API_URL = "http://67.205.168.129:8080/"

class RegistrationForm(Form):
    username = StringField('Username', [validators.Length(min=4, max=25)], description = "name")
    email = StringField('Email Address', [validators.Length(min=6, max=35)], description = "email")
    password = PasswordField('Create Password', [
        validators.DataRequired(),
        validators.EqualTo('confirm', message='Passwords must match')
    ], description = "password")
    confirm = PasswordField('Confirm Password', description = "Confirm password")

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
def login():
    if not session.get('logged_in'):
        return redirect("/login", code=302)
    else:
        return load_home()


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
        if jsonDict['success']:
            session['logged_in'] = True
            resp = make_response(redirect('/'))
            resp.set_cookie('token', jsonDict['token'])
            return resp
        else:
            flash(jsonDict['message'])
    return render_template('login.html')


@app.route('/registration', methods=['GET','POST'])
def register():
    form = RegistrationForm(request.form)
    if request.method == 'POST':
        name = form.username.data
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
        print "jsonDict['message']", jsonDict['message']
        # print "jsonDict['success']", jsonDict['success']
        if jsonDict['success']:
            flash('Thanks for registering')
            token = jsonDict['token']
            return render_template('login.html')
        else:
            flash(jsonDict['message'])
    return render_template('registration.html', form=form)


@app.route("/logout")
def logout():
    session['logged_in'] = False
    requests.post(API_URL + "user/logout", headers={'token': request.cookies.get('token')})
    resp = make_response(redirect('/'))
    resp.set_cookie('token', '', expires=0)
    return resp


@app.route("/home")
def load_home():
    if session.get('logged_in'):
        token = request.cookies.get('token')
        # Get photos from API
        getPhotos = requests.get(API_URL + "user/pictures", headers={'token': token})
        jsonDict = json.loads(getPhotos.text)

        pictureUrlArr = []
        # Add photos to array
        if (jsonDict['success'] == True):
            pictureUrlArr = []
            for pictureElements in jsonDict['pictures']:
                pictureUrlArr.append(pictureElements['url'])
        else:
            flash(jsonDict['message'])
        return render_template('home.html', imageArr = pictureUrlArr)
    else:
        return login()

@app.route("/upload")
def load_upload():
    return render_template('uploadPhotos.html')

ALLOWED_EXTENSIONS = set(['png', 'jpg', 'jpeg', 'gif'])

def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

@app.route('/uploader', methods = ['GET', 'POST'])
def upload_file():
    if request.method == 'POST':
        f = request.files['file']
        # saves files locally
        # f.save(secure_filename(f.filename))
        if f.filename == '':
            flash('No selected file')
            return redirect(request.url)
        if f and allowed_file(f.filename):
            filename = secure_filename(f.filename)
            # print request.cookies.get('token')
            req = requests.post(API_URL + "picture/upload", headers={'token': request.cookies.get('token')},
                files = {'file': (filename, f, None, None)})
            jsonDict = json.loads(req.text)
            if jsonDict['success']: # if upload successful
                flash('Upload successful')
                return render_template('home.html')
                # return str(req.status_code) + '<br/><br>' + jsonDict['url'] + '<br/><br>' + jsonDict['pictureID']
            else:
                return str(req.status_code) + ': ' + jsonDict['message']

@app.route("/delete", methods = ['GET'])
def load_delete():
    req = requests.get(API_URL + "user/pictures", headers={'token': request.cookies.get('token')})
    jsonDict = json.loads(req.text)
    # print jsonDict['success']
    imageUrlArr = []
    for pictures in jsonDict['pictures']:
        imageUrlArr.append(pictures['url'])
    return render_template('deletePhotos.html', imageArr = imageUrlArr)



if __name__ == "__main__":
    app.run(debug=True,host='0.0.0.0', port=4000)
    # max upload thingy
    app.config['MAX_CONTENT_LENGTH'] = 16 * 1024 * 1024 # 16 MB
    # TODO: throw a 404 page if filesize is too large
