import uuid

from django.db import models


class ModelCore(models.Model):
    """
    level-1
    Base сlass from which the other classes are inherited

    fields:
    - `id` uuid field
    """
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)

    class Meta:
        abstract = True







