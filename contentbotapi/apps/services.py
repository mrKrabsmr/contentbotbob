from django.db.models import QuerySet
from pydantic import BaseModel


class ServiceCore(BaseModel):
    """
    level-1
    Base сlass from which the other classes are inherited
    """
    _queryset: QuerySet


class ServiceMediaCore(ServiceCore):
    """
    level-2
    Base сlass from which the other classes are inherited

    Note: For objects where using media data
    """
    _media_queryset: QuerySet
