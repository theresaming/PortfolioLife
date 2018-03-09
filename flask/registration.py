from wtforms import Form, StringField, PasswordField, validators


class RegistrationForm(Form):
    username = StringField('Username', [validators.Length(min=4, max=25)], description = "name", render_kw={"placeholder": "name"})
    email = StringField('Email Address', [validators.Length(min=6, max=35)], description = "email", render_kw={"placeholder": "email"})
    password = PasswordField('Create Password', [
        validators.DataRequired(),
        validators.EqualTo('confirm', message='Passwords must match')
    ], description = "password", render_kw={"placeholder": "password"})
    confirm = PasswordField('Confirm Password', description = "Confirm password", render_kw={"placeholder": "retype password"})
