from wtforms import Form, StringField, PasswordField, validators


class RegistrationForm(Form):
    username = StringField('Username', [validators.Length(min=4, max=25)], description = "name")
    email = StringField('Email Address', [validators.Length(min=6, max=35)], description = "email")
    password = PasswordField('Create Password', [
        validators.DataRequired(),
        validators.EqualTo('confirm', message='Passwords must match')
    ], description = "password")
    confirm = PasswordField('Confirm Password', description = "Confirm password")
