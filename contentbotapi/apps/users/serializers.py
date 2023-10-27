from rest_framework import serializers
from rest_framework.exceptions import ValidationError
from rest_framework.serializers import Serializer, ModelSerializer

from apps.text import text_messages
from apps.users.models import User


class LoginSerializer(Serializer):
    login = serializers.CharField()
    password = serializers.CharField()


class RegisterSerializer(ModelSerializer):
    pattern = "*"

    class Meta:
        model = User
        fields = ("username", "email", "password")

    def validate(self, attrs):
        password = attrs.get("password")

        # if not re.match(self.pattern, password):
        #     raise ValidationError(text_messages["password_validate_2"])

        if "@" not in attrs.get("email"):
            raise ValidationError(text_messages["incorrect email"])

        return attrs
