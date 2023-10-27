from pydantic import dataclasses

from apps.users.models import UserCodeActivation
from apps.users.utils import get_activation_code


class UserCodeActivationService:
    _queryset = UserCodeActivation.objects.all()

    @classmethod
    def create_or_change_activation_code(cls, user):
        code_activation, created = cls._queryset.get_or_create(user=user)
        if not created:
            code_activation.code = get_activation_code()
            code_activation.save()

    @classmethod
    def get_object(cls, code):
        return cls._queryset.filter(code=code).first()
