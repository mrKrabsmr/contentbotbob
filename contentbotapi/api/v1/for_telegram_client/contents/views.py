import http

from rest_framework.response import Response
from rest_framework.views import APIView

from apps.contents.serializers import ContentSerializer
from apps.contents.services import ContentService as C_s

from config.settings import CLIENT_KEY


class ContentAPIView(APIView):

    def get(self, request, client_key, outer_id, *args, **kwargs):
        if client_key != CLIENT_KEY:
            return Response({"message": "incorrect key"}, 403)

        content = C_s.get_suitable_one(outer_id)
        serializer = ContentSerializer(content)

        return Response(serializer.data, http.HTTPStatus.OK)