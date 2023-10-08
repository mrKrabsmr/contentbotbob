from django.utils.translation import gettext_lazy as _

text_messages = {
    "password_validate_1": _("The two password fields didn’t match."),

    "password_validate_2": _("""
the password does not meet the requirements:
- The password must contain at least one lowercase letter and at least one uppercase letter.
- The password must contain at least one digit.
- The password must be at least 8 characters long.
- The password must contain at least one special character from the list @$!%*?&.
"""),

    "email_validate": _("incorrect email"),

    "sent_email": _("write the code that you received by email"),
}
