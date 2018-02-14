from flask import Flask, flash, render_template, request, session, json, make_response
import requests
import os
from werkzeug import secure_filename
from werkzeug.datastructures import FileStorage
from wtforms import Form, StringField, PasswordField, validators

app = Flask(__name__)
app.secret_key = os.urandom(12)

app.config['SQLALCHEMY_DATABASE_URI'] = 'mysql://jd:jd2018@67.205.168.129/junior_design'


# allowed image type initialization WHOOOWHEEEEEE
# ALLOWED_EXTENSIONS = ['jpg', 'jpeg', 'png']
# FILE_CONTENT_TYPES = { # these will be used to set the content type of S3 object. It is binary by default.
#     'jpg': 'image/jpeg',
#     'jpeg': 'image/jpeg',
#     'png': 'image/png'
# }

class RegistrationForm(Form):
    username = StringField('Username', [validators.Length(min=4, max=25)], description = "name")
    email = StringField('Email Address', [validators.Length(min=6, max=35)], description = "email")
    password = PasswordField('Create Password', [
        validators.DataRequired(),
        validators.EqualTo('confirm', message='Passwords must match')
    ], description = "password")
    confirm = PasswordField('Confirm Password', description = "confirm password")

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
        return render_template('login.html')
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
            token = jsonDict['token']
            return login()
        else:
            flash(jsonDict['message'])
    return login()


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
    return login()


@app.route("/home")
def load_home():
    if session.get('logged_in'):
        return render_template('home.html')
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
            #TODO: send to server
            #SQL injection protector shit
            filename = secure_filename(f.filename)
            # data = {
            #     'upload': f,read()
            #     }
            # jsonStr = json.dumps(data)
            # r = requests.post('http://67.205.168.129:8080/picture/upload', jsonStr)
            return 'file uploaded successfully'

@app.route("/delete")
def load_delete():
    return render_template('deletePhotos.html')


if __name__ == "__main__":
    app.run(debug=True,host='0.0.0.0', port=4000)
    # max upload thingy
    app.config['MAX_CONTENT_LENGTH'] = 16 * 1024 * 1024 # 16 MB
    # TODO: throw a 404 page if filesize is too large
