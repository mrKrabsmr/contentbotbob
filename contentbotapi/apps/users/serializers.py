import re

from rest_framework import serializers
from rest_framework.exceptions import ValidationError
from rest_framework.serializers import Serializer, ModelSerializer

from apps.users import text_messages
from apps.users.models import User


class LoginSerializer(Serializer):
    login = serializers.CharField()
    password = serializers.CharField()


class RegisterSerializer(ModelSerializer):
    confirm_password = serializers.CharField()

    pattern = "/.*([a-z]+[A-Z]+[0-9]+|[a-z]+[0-9]+[A-Z]+|[A-Z]+[a-z]+[0-9]+ \
    |[A-Z]+[0-9]+[a-z]+|[0-9]+[a-z]+[A-Z]+|[0-9]+[A-Z]+[a-z]+).*/;"

    class Meta:
        model = User
        fields = ("username", "email", "password")

    def validate(self, attrs):
        password = attrs.get("password")

        if password != attrs.get("confirm_password"):
            raise ValidationError(text_messages["password_validate_1"])

        if not re.match(self.pattern, password):
            raise ValidationError(text_messages["password_validate_2"])

        if "@" not in attrs.get("email"):
            raise ValidationError(text_messages["incorrect email"])

        return attrs

